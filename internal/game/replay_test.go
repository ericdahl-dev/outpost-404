package game

import (
	"path/filepath"
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

func testContentWithEvents() Content {
	c := testContent()
	c.Events = []EventDef{
		{ID: "quiet_shift", Title: "Quiet Shift", Description: "Calm.", Effects: map[string]int{"morale": 8}, MinDay: 1},
		{ID: "solar_storm", Title: "Solar Storm", Description: "Storm.", Effects: map[string]int{"power": -14}, MinDay: 2},
	}
	return c
}
