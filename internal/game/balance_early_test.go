package game

import (
	"path/filepath"
	"testing"
)

// conservativeMidScript is the mid-game reference path for issue #16 (hydro + solar through day 11+).
func conservativeMidScript(t *testing.T) SimScript {
	t.Helper()
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "conservative_mid.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}
	return script
}

func TestEarlyBalance_ConservativeMidSurvivesDay11OnAllReferenceSeeds(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script := conservativeMidScript(t)

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("Simulate: %v", err)
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
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "conservative.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("Simulate: %v", err)
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
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}

	wins := 0
	for _, strategy := range ReferenceStrategies() {
		script, err := LoadSimScript(filepath.Join(ScriptsDir(), strategy.ScriptFile))
		if err != nil {
			t.Fatalf("LoadSimScript %s: %v", strategy.ScriptFile, err)
		}
		for _, seed := range ReferenceSeeds {
			final, err := Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("Simulate: %v", err)
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
	for _, seed := range ReferenceSeeds {
		final, err := Simulate(content, seed, mid.Actions)
		if err != nil {
			t.Fatalf("Simulate: %v", err)
		}
		if !final.GameOver && final.Day >= 12 {
			day12Alive++
		}
	}
	minAlive := (len(ReferenceSeeds) + 1) / 2
	if day12Alive < minAlive {
		t.Fatalf("no wins and only %d/%d seeds reached day 12 alive on conservative_mid (need %d)",
			day12Alive, len(ReferenceSeeds), minAlive)
	}
}
