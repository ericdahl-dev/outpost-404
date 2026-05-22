package game

import (
	"path/filepath"
	"testing"
)

func TestSurvival30_WinsOnAllReferenceSeeds(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_30.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}
	strategy := survival30Baseline()

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("Simulate: %v", err)
			}
			if err := CheckBaselineOutcome(strategy, seed, final); err != nil {
				t.Fatal(err)
			}
			if !final.Won || final.Day != 31 {
				t.Fatalf("seed %d: day=%d won=%v", seed, final.Day, final.Won)
			}
		})
	}
}

func TestSurvival30_NoVitalHitsZeroDuringScript(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_30.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, snaps, err := SimulateWithSnapshots(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("SimulateWithSnapshots: %v", err)
			}
			for i, snap := range snaps {
				if snap.GameOver {
					continue
				}
				if snapshotVitalsDepleted(snap) {
					t.Fatalf("seed %d action %d day %d: vitals depleted (power=%d food=%d morale=%d pop=%d)",
						seed, i, snap.Day, snap.Power, snap.Food, snap.Morale, snap.Population)
				}
			}
			if !final.Won {
				t.Fatalf("seed %d: expected win, got day=%d game_over=%v", seed, final.Day, final.GameOver)
			}
		})
	}
}

func TestSurvival30_EndMarginsOnAllReferenceSeeds(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_30.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, snaps, err := SimulateWithSnapshots(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("SimulateWithSnapshots: %v", err)
			}
			end := snaps[len(snaps)-1]
			if !snapshotMeetsSurvivalEndMargins(end) {
				t.Fatalf("seed %d end: power=%d food=%d (want power>=%d food>=%d)",
					seed, end.Power, end.Food, SurvivalMinEndPower, SurvivalMinEndFood)
			}
			if final.Power != end.Power || final.Food != end.Food {
				t.Fatalf("seed %d: final state mismatch with last snapshot", seed)
			}
		})
	}
}

func TestCheckEnd_SurvivalLossWhenDay31ButFoodDepleted(t *testing.T) {
	s := newTestState()
	s.Day = 31
	s.Food = 0
	s.Power = 50
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if !s.GameOver || s.Won {
		t.Fatal("expected collapse on day 31 when food is 0, not survival win")
	}
}

func TestCheckEnd_SurvivalWinWhenDay31AndVitalsPositive(t *testing.T) {
	s := newTestState()
	s.Day = 31
	s.Food = 10
	s.Power = 15
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if !s.GameOver || !s.Won {
		t.Fatal("expected survival win with positive vitals on day 31")
	}
}

func TestCheckEnd_BeaconWinStillBeatsCollapse(t *testing.T) {
	s := newTestState()
	s.BeaconParts = 5
	s.Power = 0
	s.Food = 0

	s.CheckEnd()

	if !s.Won {
		t.Fatal("beacon completion should win even when vitals are depleted")
	}
}

func TestSurvival30_ScriptHasThirtyDayAdvances(t *testing.T) {
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_30.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}
	n := 0
	for _, a := range script.Actions {
		if a.Type == "next_day" {
			n++
		}
	}
	if n != 30 {
		t.Fatalf("survival_30.json has %d next_day actions, want 30", n)
	}
}

