package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"path/filepath"
	"testing"
)

// conservativeMidScript is the mid-game reference path for issue #16 (hydro + solar through day 11+).
func conservativeMidScript(t *testing.T) game.SimScript {
	t.Helper()
	script, err := game.LoadSimScript(filepath.Join(game.ScriptsDir(), "conservative_mid.json"))
	if err != nil {
		t.Fatalf("game.LoadSimScript: %v", err)
	}
	return script
}

func TestEarlyBalance_ConservativeMidSurvivesDay11OnAllReferenceSeeds(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}
	script := conservativeMidScript(t)

	for _, seed := range game.ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := game.Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("game.Simulate: %v", err)
			}
			if final.GameOver {
				t.Fatalf("collapsed on day %d (food=%d power=%d)", final.Day, final.Food, final.Power)
			}
			if final.Day < 11 {
				t.Fatalf("day=%d, want >= 11", final.Day)
			}
		})
	}
}

func TestEarlyBalance_ConservativeShortAliveThroughDay5OnAllReferenceSeeds(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}
	script, err := game.LoadSimScript(filepath.Join(game.ScriptsDir(), "conservative.json"))
	if err != nil {
		t.Fatalf("game.LoadSimScript: %v", err)
	}

	for _, seed := range game.ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := game.Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("game.Simulate: %v", err)
			}
			if final.GameOver || final.Day < 5 {
				t.Fatalf("day=%d game_over=%v", final.Day, final.GameOver)
			}
		})
	}
}

// TestEarlyBalance_ReferenceWinOrMidGameMilestone requires a win on some reference run,
// or at least half of reference seeds reaching day 12 alive on conservative_mid.
func TestEarlyBalance_ReferenceWinOrMidGameMilestone(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}

	wins := 0
	for _, strategy := range game.ReferenceStrategies() {
		script, err := game.LoadSimScript(filepath.Join(game.ScriptsDir(), strategy.ScriptFile))
		if err != nil {
			t.Fatalf("game.LoadSimScript %s: %v", strategy.ScriptFile, err)
		}
		for _, seed := range game.ReferenceSeeds {
			final, err := game.Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("game.Simulate: %v", err)
			}
			if final.Won {
				wins++
			}
		}
	}

	if wins > 0 {
		return
	}

	mid := conservativeMidScript(t)
	day12Alive := 0
	for _, seed := range game.ReferenceSeeds {
		final, err := game.Simulate(content, seed, mid.Actions)
		if err != nil {
			t.Fatalf("game.Simulate: %v", err)
		}
		if !final.GameOver && final.Day >= 12 {
			day12Alive++
		}
	}
	minAlive := (len(game.ReferenceSeeds) + 1) / 2
	if day12Alive < minAlive {
		t.Fatalf("no wins and only %d/%d seeds reached day 12 alive on conservative_mid (need %d)",
			day12Alive, len(game.ReferenceSeeds), minAlive)
	}
}
