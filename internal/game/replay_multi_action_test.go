package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// multiActionSeed yields a logged next_day with event_id under testContentWithEvents.
const multiActionSeed = int64(4242)

func TestReplaySession_MultiActionMatchesLiveSession(t *testing.T) {
	content := testContentWithEvents()
	live, entries := recordMultiActionSession(t, content, multiActionSeed)

	assertSessionHasActionTypes(t, entries, "build", "damage", "repair", "trade", "beacon", "next_day")
	assertNextDayLogsEventID(t, entries)

	replayed, err := ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("ReplaySession: %v", err)
	}
	if diff := snapshotDiff(normalizeSnapshot(live.snapshot()), normalizeSnapshot(replayed.snapshot())); diff != "" {
		t.Fatalf("final snapshot mismatch: %s", diff)
	}
}

func TestReplaySession_FixtureJSONLMatchesLiveSession(t *testing.T) {
	content := testContentWithEvents()
	live, _ := recordMultiActionSession(t, content, multiActionSeed)

	path := filepath.Join("testdata", "sessions", "multi_action.jsonl")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile fixture: %v", err)
	}
	entries, err := parseSessionLogLines(raw)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	assertSessionHasActionTypes(t, entries, "build", "damage", "repair", "trade", "beacon", "next_day")

	replayed, err := ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("ReplaySession fixture: %v", err)
	}
	if diff := snapshotDiff(normalizeSnapshot(live.snapshot()), normalizeSnapshot(replayed.snapshot())); diff != "" {
		t.Fatalf("fixture replay mismatch vs live: %s", diff)
	}
}

func TestReplaySession_GameEndRestoresTerminalState(t *testing.T) {
	content := testContentWithEvents()
	live, entries := recordBeaconWinSession(t, content, 9001)

	gameEnd, ok := findLogEntry(entries, "game_end")
	if !ok {
		t.Fatal("session missing game_end entry")
	}
	if gameEnd.After == nil || !gameEnd.After.Won || !gameEnd.After.GameOver {
		t.Fatalf("game_end after snapshot: %+v", gameEnd.After)
	}

	replayed, err := ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("ReplaySession: %v", err)
	}
	if diff := snapshotDiff(normalizeSnapshot(*gameEnd.After), normalizeSnapshot(replayed.snapshot())); diff != "" {
		t.Fatalf("game_end.After mismatch: %s", diff)
	}
	if !replayed.GameOver || !replayed.Won {
		t.Fatalf("terminal flags: game_over=%v won=%v", replayed.GameOver, replayed.Won)
	}
	if replayed.Message != live.Message {
		t.Fatalf("message = %q, want %q", replayed.Message, live.Message)
	}
	if replayed.BeaconParts != live.MaxBeaconParts {
		t.Fatalf("beacon parts = %d, want %d", replayed.BeaconParts, live.MaxBeaconParts)
	}
}

func recordMultiActionSession(t *testing.T, content Content, seed int64) (State, []LogEntry) {
	t.Helper()
	return recordSession(t, content, seed, func(s *State) {
		s.Build("solar_array")
		s.DamageBuilding("solar_array")
		s.RepairBuilding("solar_array")
		s.Trade()
		s.WorkOnBeacon()
		s.NextDay()
	})
}

func recordBeaconWinSession(t *testing.T, content Content, seed int64) (State, []LogEntry) {
	t.Helper()
	return recordSession(t, content, seed, func(s *State) {
		s.BeaconParts = 4
		s.Power = 50
		s.Credits = 100
		s.WorkOnBeacon()
	})
}

func recordSession(t *testing.T, content Content, seed int64, play func(*State)) (State, []LogEntry) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	logger, err := OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}

	s := NewStateWithSeed(content, seed)
	s.SessionLog = logger
	s.LogSessionStart()
	play(&s)
	_ = logger.Close()

	entries, err := LoadSessionLog(path)
	if err != nil {
		t.Fatalf("LoadSessionLog: %v", err)
	}
	return s, entries
}

func assertSessionHasActionTypes(t *testing.T, entries []LogEntry, types ...string) {
	t.Helper()
	seen := make(map[string]bool, len(types))
	for _, e := range entries {
		seen[e.Type] = true
	}
	for _, typ := range types {
		if !seen[typ] {
			t.Fatalf("session missing %q entry; types present: %v", typ, logEntryTypes(entries))
		}
	}
}

func assertNextDayLogsEventID(t *testing.T, entries []LogEntry) {
	t.Helper()
	for _, e := range entries {
		if e.Type != "next_day" {
			continue
		}
		id, _ := e.Detail["event_id"].(string)
		if id != "" {
			return
		}
	}
	t.Fatalf("expected next_day with event_id in detail; entries: %v", logEntryTypes(entries))
}

func findLogEntry(entries []LogEntry, typ string) (LogEntry, bool) {
	for _, e := range entries {
		if e.Type == typ {
			return e, true
		}
	}
	return LogEntry{}, false
}

func logEntryTypes(entries []LogEntry) []string {
	var types []string
	for _, e := range entries {
		types = append(types, e.Type)
	}
	return types
}

func TestRegenerateMultiActionFixture(t *testing.T) {
	if os.Getenv("UPDATE_SESSION_FIXTURE") != "1" {
		t.Skip("set UPDATE_SESSION_FIXTURE=1 to rewrite testdata/sessions/multi_action.jsonl")
	}
	content := testContentWithEvents()
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	logger, err := OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}
	s := NewStateWithSeed(content, multiActionSeed)
	s.SessionLog = logger
	s.LogSessionStart()
	s.Build("solar_array")
	s.DamageBuilding("solar_array")
	s.RepairBuilding("solar_array")
	s.Trade()
	s.WorkOnBeacon()
	s.NextDay()
	_ = logger.Close()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	out := filepath.Join("testdata", "sessions", "multi_action.jsonl")
	if err := os.WriteFile(out, raw, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func parseSessionLogLines(raw []byte) ([]LogEntry, error) {
	var entries []LogEntry
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		if line == "" {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
