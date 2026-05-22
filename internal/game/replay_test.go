package game

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNextDay_SameSeedProducesSameOutcome(t *testing.T) {
	content := testContentWithEvents()
	s1 := NewStateWithSeed(content, 42)
	s2 := NewStateWithSeed(content, 42)

	s1.NextDay()
	s2.NextDay()

	if diff := snapshotDiff(normalizeSnapshot(s1.snapshot()), normalizeSnapshot(s2.snapshot())); diff != "" {
		t.Fatalf("snapshots differ: %s", diff)
	}
}

func TestSimulate_RunsScriptDeterministically(t *testing.T) {
	content := testContentWithEvents()
	actions := []SimAction{
		{Type: "next_day"},
		{Type: "next_day"},
		{Type: "repair"},
	}

	a, err := Simulate(content, 99, actions)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}
	b, err := Simulate(content, 99, actions)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}
	if diff := snapshotDiff(normalizeSnapshot(a.snapshot()), normalizeSnapshot(b.snapshot())); diff != "" {
		t.Fatalf("runs differ: %s", diff)
	}
}

func TestReplaySession_MatchesRecordedLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	content := testContentWithEvents()

	orig := NewStateWithSeed(content, 7)
	logger, err := OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}
	orig.SessionLog = logger
	orig.LogSessionStart()
	orig.Build("solar_array")
	orig.NextDay()
	orig.NextDay()
	_ = logger.Close()

	entries, err := LoadSessionLog(path)
	if err != nil {
		t.Fatalf("LoadSessionLog: %v", err)
	}
	replayed, err := ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("ReplaySession: %v", err)
	}
	if diff := snapshotDiff(normalizeSnapshot(replayed.snapshot()), normalizeSnapshot(orig.snapshot())); diff != "" {
		t.Fatalf("final snapshot mismatch: %s\norig %#v\nreplay %#v", diff, orig.snapshot(), replayed.snapshot())
	}
}

func TestReplaySession_MismatchAfterBuildReturnsError(t *testing.T) {
	content := testContentWithEvents()
	start := Snapshot{Day: 1, Power: 65, Food: 60, Morale: 70, Credits: 180, Population: 8, PopulationCap: 10, BeaconParts: 0, MaxBeacon: 5}
	wrongAfter := start
	wrongAfter.Credits = 999

	entries := []LogEntry{
		{
			Type:     "session_start",
			Snapshot: &start,
			Detail:   map[string]any{"seed": "7"},
		},
		{
			Type:   "build",
			Day:    1,
			Detail: map[string]any{"building_id": "solar_array", "ok": true},
			Before: &start,
			After:  &wrongAfter,
		},
	}

	_, err := ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected replay mismatch error")
	}
	if !containsAll(err.Error(), "replay mismatch", "build", "credits") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_UnsupportedEntryTypeReturnsError(t *testing.T) {
	content := testContentWithEvents()
	start := Snapshot{Day: 1, Power: 65, Credits: 180, MaxBeacon: 5}

	entries := []LogEntry{
		{Type: "session_start", Snapshot: &start, Detail: map[string]any{"seed": "7"}},
		{Type: "teleport", Day: 1, Before: &start, After: &start},
	}

	_, err := ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected unsupported type error")
	}
	if !containsAll(err.Error(), "unsupported", "teleport") {
		t.Fatalf("error = %q", err)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}

func TestReplaySession_UserLogFile(t *testing.T) {
	path := filepath.Join("..", "..", "logs", "my-game.jsonl")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("logs/my-game.jsonl not present")
	}
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	entries, err := LoadSessionLog(path)
	if err != nil {
		t.Fatalf("LoadSessionLog: %v", err)
	}
	if _, err := ReplaySession(content, entries); err != nil {
		if strings.Contains(err.Error(), "replay mismatch") {
			t.Skipf("session log predates balance pass (#16); re-record: %v", err)
		}
		t.Fatalf("ReplaySession: %v", err)
	}
}

func testContentWithEvents() Content {
	c := testContent()
	c.Events = []EventDef{
		{ID: "quiet_shift", Title: "Quiet Shift", Description: "Calm.", Effects: map[string]int{"morale": 8}, MinDay: 1},
		{ID: "solar_storm", Title: "Solar Storm", Description: "Storm.", Effects: map[string]int{"power": -14}, MinDay: 2},
	}
	return c
}
