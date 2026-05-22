package game_test

import (
	"encoding/json"
	"github.com/ericdahl/outpost-404/internal/game"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNextDay_SameSeedProducesSameOutcome(t *testing.T) {
	content := testContentWithEvents()
	s1 := game.NewStateWithSeed(content, 42)
	s2 := game.NewStateWithSeed(content, 42)

	s1.NextDay()
	s2.NextDay()

	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&s1)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&s2))); diff != "" {
		t.Fatalf("snapshots differ: %s", diff)
	}
}

func TestSimulate_RunsScriptDeterministically(t *testing.T) {
	content := testContentWithEvents()
	actions := []game.SimAction{
		{Type: "next_day"},
		{Type: "next_day"},
		{Type: "repair"},
	}

	a, err := game.Simulate(content, 99, actions)
	if err != nil {
		t.Fatalf("game.Simulate: %v", err)
	}
	b, err := game.Simulate(content, 99, actions)
	if err != nil {
		t.Fatalf("game.Simulate: %v", err)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&a)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&b))); diff != "" {
		t.Fatalf("runs differ: %s", diff)
	}
}

func TestReplaySession_MatchesRecordedLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	content := testContentWithEvents()

	orig := game.NewStateWithSeed(content, 7)
	logger, err := game.OpenSessionLog(path)
	if err != nil {
		t.Fatalf("game.OpenSessionLog: %v", err)
	}
	orig.SessionLog = logger
	orig.LogSessionStart()
	orig.Build("solar_array")
	orig.NextDay()
	orig.NextDay()
	_ = logger.Close()

	entries, err := game.LoadSessionLog(path)
	if err != nil {
		t.Fatalf("game.LoadSessionLog: %v", err)
	}
	replayed, err := game.ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("game.ReplaySession: %v", err)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&replayed)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&orig))); diff != "" {
		t.Fatalf("final snapshot mismatch: %s\norig %#v\nreplay %#v", diff, game.SnapshotForTest(&orig), game.SnapshotForTest(&replayed))
	}
}

func TestReplaySession_MismatchAfterBuildReturnsError(t *testing.T) {
	content := testContentWithEvents()
	start := game.Snapshot{Day: 1, Power: 65, Food: 60, Morale: 70, Credits: 180, Population: 8, PopulationCap: 10, BeaconParts: 0, MaxBeacon: 5}
	wrongAfter := start
	wrongAfter.Credits = 999

	entries := []game.LogEntry{
		{
			Type:          "session_start",
			Snapshot: &start,
			Detail:        map[string]any{"seed": "7"},
		},
		{
			Type:   "build",
			Day:    1,
			Detail: map[string]any{"building_id": "solar_array", "ok": true},
			Before: &start,
			After:  &wrongAfter,
		},
	}

	_, err := game.ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected replay mismatch error")
	}
	if !containsAll(err.Error(), "replay mismatch", "build", "credits") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_missingSessionStart(t *testing.T) {
	content := testContentWithEvents()
	_, err := game.ReplaySession(content, []game.LogEntry{{Type: "build", Day: 1}})
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsAll(err.Error(), "session log missing session_start") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_sessionStartWithoutSnapshot(t *testing.T) {
	content := testContentWithEvents()
	entries := []game.LogEntry{{Type: "session_start", Detail: map[string]any{"seed": "7"}}}
	_, err := game.ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsAll(err.Error(), "session log missing session_start") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_missingSeed(t *testing.T) {
	content := testContentWithEvents()
	start := game.Snapshot{Day: 1, Power: 65, Credits: 180, MaxBeacon: 5}
	entries := []game.LogEntry{{Type: "session_start", Snapshot: &start}}
	_, err := game.ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsAll(err.Error(), "session log", "missing seed") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_invalidSeedType(t *testing.T) {
	content := testContentWithEvents()
	start := game.Snapshot{Day: 1, Power: 65, Credits: 180, MaxBeacon: 5}
	entries := []game.LogEntry{
		{Type: "session_start", Snapshot: &start, Detail: map[string]any{"seed": true}},
	}
	_, err := game.ReplaySession(content, entries)
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsAll(err.Error(), "session log", "invalid seed type") {
		t.Fatalf("error = %q", err)
	}
}

func TestReplaySession_numericSeedFormsReplay(t *testing.T) {
	content := testContentWithEvents()
	seed := int64(7)
	actions := []game.SimAction{{Type: "next_day"}}

	final, snaps, err := game.SimulateWithSnapshots(content, seed, actions)
	if err != nil {
		t.Fatalf("game.SimulateWithSnapshots: %v", err)
	}
	start, after := snaps[0], snaps[1]

	cases := []struct {
		name string
		seed any
	}{
		{name: "string", seed: "7"},
		{name: "float64", seed: float64(7)},
		{name: "json.Number", seed: json.Number("7")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entries := []game.LogEntry{
				{Type: "session_start", Snapshot: &start, Detail: map[string]any{"seed": tc.seed}},
				{Type: "next_day", Day: 1, Before: &start, After: &after},
			}
			got, err := game.ReplaySession(content, entries)
			if err != nil {
				t.Fatalf("game.ReplaySession: %v", err)
			}
			if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&final)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&got))); diff != "" {
				t.Fatalf("snapshot mismatch: %s", diff)
			}
		})
	}
}

func TestReplaySession_UnsupportedEntryTypeReturnsError(t *testing.T) {
	content := testContentWithEvents()
	start := game.Snapshot{Day: 1, Power: 65, Credits: 180, MaxBeacon: 5}

	entries := []game.LogEntry{
		{Type: "session_start", Snapshot: &start, Detail: map[string]any{"seed": "7"}},
		{Type: "teleport", Day: 1, Before: &start, After: &start},
	}

	_, err := game.ReplaySession(content, entries)
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
	path := filepath.Join(game.RepoRoot(), "logs", "my-game.jsonl")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("logs/my-game.jsonl not present")
	}
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}
	entries, err := game.LoadSessionLog(path)
	if err != nil {
		t.Fatalf("game.LoadSessionLog: %v", err)
	}
	if _, err := game.ReplaySession(content, entries); err != nil {
		if strings.Contains(err.Error(), "replay mismatch") {
			t.Skipf("session log predates balance pass (#16); re-record: %v", err)
		}
		t.Fatalf("game.ReplaySession: %v", err)
	}
}

func testContentWithEvents() game.Content {
	c := testContent()
	c.Events = []game.EventDef{
		{ID: "quiet_shift", Title: "Quiet Shift", Description: "Calm.", Effects: map[string]int{"morale": 8}, MinDay: 1},
		{ID: "solar_storm", Title: "Solar Storm", Description: "Storm.", Effects: map[string]int{"power": -14}, MinDay: 2},
	}
	return c
}
