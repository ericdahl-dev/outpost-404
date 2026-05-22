package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSimScript_Array(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.json")
	if err := os.WriteFile(path, []byte(`[
  {"type": "build", "building_id": "solar_array"},
  {"type": "next_day"}
]`), 0o644); err != nil {
		t.Fatal(err)
	}

	script, err := game.LoadSimScript(path)
	if err != nil {
		t.Fatalf("game.LoadSimScript: %v", err)
	}
	if script.SeedSet {
		t.Fatal("expected no seed in array-only script")
	}
	if len(script.Actions) != 2 {
		t.Fatalf("actions: got %d", len(script.Actions))
	}
	if script.Actions[0].Type != "build" || script.Actions[0].BuildingID != "solar_array" {
		t.Fatalf("first action: %#v", script.Actions[0])
	}
}

func TestLoadSimScript_WrappedWithSeed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.json")
	if err := os.WriteFile(path, []byte(`{"seed": 42, "actions": [{"type": "repair"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	script, err := game.LoadSimScript(path)
	if err != nil {
		t.Fatalf("game.LoadSimScript: %v", err)
	}
	if !script.SeedSet || script.Seed != 42 {
		t.Fatalf("seed: set=%v val=%d", script.SeedSet, script.Seed)
	}
	if len(script.Actions) != 1 || script.Actions[0].Type != "repair" {
		t.Fatalf("actions: %#v", script.Actions)
	}
}

func TestResolveSimSeed_FlagWins(t *testing.T) {
	got, err := game.ResolveSimSeed(true, 1, 99)
	if err != nil || got != 99 {
		t.Fatalf("got %d err %v", got, err)
	}
}

func TestParseSeedList(t *testing.T) {
	got, err := game.ParseSeedList("1, 42,99")
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{1, 42, 99}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestResolveSimSeed_RequiresSeed(t *testing.T) {
	if _, err := game.ResolveSimSeed(false, 0, 0); err == nil {
		t.Fatal("expected error without seed")
	}
}

func TestSimulateSeeds_MatchesIndividualSimulate(t *testing.T) {
	content := testContentWithEvents()
	actions := []game.SimAction{{Type: "next_day"}, {Type: "repair"}}
	seeds := []int64{7, 99}

	batch, err := game.SimulateSeeds(content, seeds, actions)
	if err != nil {
		t.Fatalf("game.SimulateSeeds: %v", err)
	}
	if len(batch) != len(seeds) {
		t.Fatalf("got %d states, want %d", len(batch), len(seeds))
	}

	for i, seed := range seeds {
		single, err := game.Simulate(content, seed, actions)
		if err != nil {
			t.Fatalf("game.Simulate seed %d: %v", seed, err)
		}
		if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&batch[i])), game.NormalizeSnapshotForTest(game.SnapshotForTest(&single))); diff != "" {
			t.Fatalf("seed %d batch vs single: %s", seed, diff)
		}
	}
}

func TestFormatSimOutcome(t *testing.T) {
	s := game.State{Day: 3, Won: false, GameOver: true, BeaconParts: 1, MaxBeaconParts: 5,
		Power: 10, Food: 0, Morale: 50, Credits: 20}
	out := game.FormatSimOutcome(7, s)
	if out == "" {
		t.Fatal("empty outcome")
	}
	for _, want := range []string{"seed=7", "day=3", "game_over=true", "food=0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in %q", want, out)
		}
	}
}
