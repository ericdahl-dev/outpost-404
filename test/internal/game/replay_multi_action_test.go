package game_test

import (
	"encoding/json"
	"github.com/ericdahl/outpost-404/internal/game"
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

	assertSessionHasActionTypes(t, entries, "build", "repair", "trade", "beacon", "next_day")
	assertNextDayLogsEventID(t, entries)

	replayed, err := game.ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("game.ReplaySession: %v", err)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&live)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&replayed))); diff != "" {
		t.Fatalf("final snapshot mismatch: %s", diff)
	}
}

func TestReplaySession_FixtureJSONLMatchesLiveSession(t *testing.T) {
	content := testContentWithEvents()
	live, _ := recordMultiActionSession(t, content, multiActionSeed)

	path := game.TestFixturePath("sessions", "multi_action.jsonl")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile fixture: %v", err)
	}
	entries, err := parseSessionLogLines(raw)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	assertSessionHasActionTypes(t, entries, "build", "repair", "trade", "beacon", "next_day")

	replayed, err := game.ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("game.ReplaySession fixture: %v", err)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&live)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&replayed))); diff != "" {
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

	replayed, err := game.ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("game.ReplaySession: %v", err)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(*gameEnd.After), game.NormalizeSnapshotForTest(game.SnapshotForTest(&replayed))); diff != "" {
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

func recordMultiActionSession(t *testing.T, content game.Content, seed int64) (game.State, []game.LogEntry) {
	t.Helper()
	return recordSession(t, content, seed, func(s *game.State) {
		s.Build("solar_array")
		s.Repair()
		s.Trade()
		s.WorkOnBeacon()
		s.NextDay()
	})
}

func recordBeaconWinSession(t *testing.T, content game.Content, seed int64) (game.State, []game.LogEntry) {
	t.Helper()
	return recordSession(t, content, seed, func(s *game.State) {
		s.BeaconParts = 4
		s.Power = 50
		s.Credits = 100
		s.WorkOnBeacon()
	})
}

func recordSession(t *testing.T, content game.Content, seed int64, play func(*game.State)) (game.State, []game.LogEntry) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	logger, err := game.OpenSessionLog(path)
	if err != nil {
		t.Fatalf("game.OpenSessionLog: %v", err)
	}

	s := game.NewStateWithSeed(content, seed)
	s.SessionLog = logger
	s.LogSessionStart()
	play(&s)
	_ = logger.Close()

	entries, err := game.LoadSessionLog(path)
	if err != nil {
		t.Fatalf("game.LoadSessionLog: %v", err)
	}
	return s, entries
}

func assertSessionHasActionTypes(t *testing.T, entries []game.LogEntry, types ...string) {
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

func assertNextDayLogsEventID(t *testing.T, entries []game.LogEntry) {
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

func findLogEntry(entries []game.LogEntry, typ string) (game.LogEntry, bool) {
	for _, e := range entries {
		if e.Type == typ {
			return e, true
		}
	}
	return game.LogEntry{}, false
}

func logEntryTypes(entries []game.LogEntry) []string {
	var types []string
	for _, e := range entries {
		types = append(types, e.Type)
	}
	return types
}

func parseSessionLogLines(raw []byte) ([]game.LogEntry, error) {
	var entries []game.LogEntry
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		if line == "" {
			continue
		}
		var entry game.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
