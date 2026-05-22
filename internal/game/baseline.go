package game

import (
	"fmt"
	"path/filepath"
)

// ReferenceSeeds are fixed RNG seeds for balance regression checks.
// 1779403310247544000 matches a real session_start seed (see seed_json_test).
var ReferenceSeeds = []int64{
	1,
	7,
	42,
	99,
	100,
	101,
	1779403310247544000,
}

// BaselineOutcome is the end state recorded for a strategy/seed pair.
type BaselineOutcome struct {
	Day         int
	GameOver    bool
	Won         bool
	BeaconParts int
}

// BaselineStrategy is a reference script plus viability rules and optional exact outcomes.
type BaselineStrategy struct {
	ID             string
	ScriptFile     string
	MinEndDay      int
	RequireAlive   bool
	MinBeaconParts int
	Expected       map[int64]BaselineOutcome
}

// ReferenceStrategies returns the balance baseline script set.
func ReferenceStrategies() []BaselineStrategy {
	return []BaselineStrategy{
		conservativeBaseline(),
		noTradeSurvivalBaseline(),
		beaconRushBaseline(),
		survival30Baseline(),
	}
}

// ScriptsDir returns the repo scripts/ directory from internal/game tests or tools.
func ScriptsDir() string {
	return filepath.Join("..", "..", "scripts")
}

// CheckBaselineOutcome verifies a simulation result against strategy rules.
func CheckBaselineOutcome(strategy BaselineStrategy, seed int64, s State) error {
	if s.Day < strategy.MinEndDay {
		return fmt.Errorf("seed %d: day %d < min %d", seed, s.Day, strategy.MinEndDay)
	}
	if strategy.RequireAlive && s.GameOver {
		return fmt.Errorf("seed %d: game over on day %d (expected alive at script end)", seed, s.Day)
	}
	if strategy.MinBeaconParts > 0 && s.BeaconParts < strategy.MinBeaconParts {
		return fmt.Errorf("seed %d: beacon %d < min %d", seed, s.BeaconParts, strategy.MinBeaconParts)
	}

	want, ok := strategy.Expected[seed]
	if !ok {
		return fmt.Errorf("seed %d: missing expected outcome for strategy %q", seed, strategy.ID)
	}
	got := BaselineOutcome{
		Day:         s.Day,
		GameOver:    s.GameOver,
		Won:         s.Won,
		BeaconParts: s.BeaconParts,
	}
	if got != want {
		return fmt.Errorf("seed %d: got %+v want %+v", seed, got, want)
	}
	return nil
}

func conservativeBaseline() BaselineStrategy {
	return BaselineStrategy{
		ID:           "conservative",
		ScriptFile:   "conservative.json",
		MinEndDay:    5,
		RequireAlive: true,
		Expected: map[int64]BaselineOutcome{
			1:                   {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			7:                   {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			42:                  {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			99:                  {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			100:                 {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			101:                 {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
			1779403310247544000: {Day: 5, GameOver: false, Won: false, BeaconParts: 0},
		},
	}
}

func noTradeSurvivalBaseline() BaselineStrategy {
	return BaselineStrategy{
		ID:         "no_trade_survival",
		ScriptFile: "no_trade_survival.json",
		MinEndDay:  14,
		Expected: map[int64]BaselineOutcome{
			1:                   {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			7:                   {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			42:                  {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			99:                  {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			100:                 {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			101:                 {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
			1779403310247544000: {Day: 14, GameOver: false, Won: false, BeaconParts: 0},
		},
	}
}

func beaconRushBaseline() BaselineStrategy {
	return BaselineStrategy{
		ID:             "beacon_rush",
		ScriptFile:     "beacon_rush.json",
		MinEndDay:      6,
		MinBeaconParts: 2,
		Expected: map[int64]BaselineOutcome{
			1:                   {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			7:                   {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			42:                  {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			99:                  {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			100:                 {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			101:                 {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
			1779403310247544000: {Day: 6, GameOver: false, Won: false, BeaconParts: 2},
		},
	}
}

func survival30Baseline() BaselineStrategy {
	win := BaselineOutcome{Day: 31, GameOver: true, Won: true, BeaconParts: 0}
	return BaselineStrategy{
		ID:         "survival_30",
		ScriptFile: "survival_30.json",
		MinEndDay:  31,
		Expected: map[int64]BaselineOutcome{
			1:                   win,
			7:                   win,
			42:                  win,
			99:                  win,
			100:                 win,
			101:                 win,
			1779403310247544000: win,
		},
	}
}
