package game

import (
	"path/filepath"
	"testing"
)

func survivalWinDay() int {
	return SurvivalWinAfterDay + 1
}

// Survival win threshold (vertical slice 1): healthy colony on the last survival day is not a win yet.
func TestCheckEnd_NoSurvivalWinOnLastSurvivalDay(t *testing.T) {
	s := newTestState()
	s.Day = SurvivalWinAfterDay
	s.Power = 50
	s.Food = 50
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if s.GameOver || s.Won {
		t.Fatalf("day %d with positive vitals should keep playing (win on day %d)", s.Day, survivalWinDay())
	}
}

func TestSurvival45_WinsOnAllReferenceSeeds(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_45.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}
	strategy := survival45Baseline()

	for _, seed := range ReferenceSeeds {
		t.Run(seedLabel(seed), func(t *testing.T) {
			final, err := Simulate(content, seed, script.Actions)
			if err != nil {
				t.Fatalf("Simulate: %v", err)
			}
			if err := CheckBaselineOutcome(strategy, seed, final); err != nil {
				t.Fatal(err)
			}
			wantDay := survivalWinDay()
			if !final.Won || final.Day != wantDay {
				t.Fatalf("seed %d: day=%d won=%v, want day=%d won", seed, final.Day, final.Won, wantDay)
			}
		})
	}
}

func TestSurvival45_NoVitalHitsZeroDuringScript(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_45.json"))
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

func TestSurvival45_EndMarginsOnAllReferenceSeeds(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_45.json"))
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

func TestCheckEnd_SurvivalLossWhenWinDayButFoodDepleted(t *testing.T) {
	s := newTestState()
	s.Day = survivalWinDay()
	s.Food = 0
	s.Power = 50
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if !s.GameOver || s.Won {
		t.Fatalf("expected collapse on day %d when food is 0, not survival win", s.Day)
	}
}

func TestCheckEnd_SurvivalWinWhenWinDayAndVitalsPositive(t *testing.T) {
	s := newTestState()
	s.Day = survivalWinDay()
	s.Food = 10
	s.Power = 15
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if !s.GameOver || !s.Won {
		t.Fatalf("expected survival win with positive vitals on day %d", s.Day)
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

func TestSurvival45_ScriptHasFortyFiveDayAdvances(t *testing.T) {
	script, err := LoadSimScript(filepath.Join(ScriptsDir(), "survival_45.json"))
	if err != nil {
		t.Fatalf("LoadSimScript: %v", err)
	}
	n := 0
	for _, a := range script.Actions {
		if a.Type == "next_day" {
			n++
		}
	}
	if n != SurvivalWinAfterDay {
		t.Fatalf("survival_45.json has %d next_day actions, want %d", n, SurvivalWinAfterDay)
	}
}
