package game_test

import (
	"fmt"
	"github.com/ericdahl/outpost-404/internal/game"
	"path/filepath"
	"testing"
)

func TestReferenceSeeds_Count(t *testing.T) {
	if len(game.ReferenceSeeds) < 5 {
		t.Fatalf("game.ReferenceSeeds: got %d, want at least 5", len(game.ReferenceSeeds))
	}
}

func TestReferenceStrategies_Count(t *testing.T) {
	if len(game.ReferenceStrategies()) < 4 {
		t.Fatalf("game.ReferenceStrategies: got %d, want at least 4", len(game.ReferenceStrategies()))
	}
}

func TestBalanceBaseline_ReferenceStrategies(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}

	for _, strategy := range game.ReferenceStrategies() {
		t.Run(strategy.ID, func(t *testing.T) {
			path := filepath.Join(game.ScriptsDir(), strategy.ScriptFile)
			script, err := game.LoadSimScript(path)
			if err != nil {
				t.Fatalf("game.LoadSimScript(%s): %v", path, err)
			}

			for _, seed := range game.ReferenceSeeds {
				t.Run(seedLabel(seed), func(t *testing.T) {
					final, err := game.Simulate(content, seed, script.Actions)
					if err != nil {
						t.Fatalf("game.Simulate: %v", err)
					}
					if err := game.CheckBaselineOutcome(strategy, seed, final); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}

func seedLabel(seed int64) string {
	return fmt.Sprintf("seed_%d", seed)
}

func TestCheckBaselineOutcome_errors(t *testing.T) {
	const seed = int64(42)
	strategy := game.BaselineStrategy{
		ID:             "test_strategy",
		MinEndDay:      10,
		MinBeaconParts: 2,
		RequireAlive:   true,
		Expected: map[int64]game.BaselineOutcome{
			seed: {Day: 14, GameOver: false, Won: false, BeaconParts: 3},
		},
	}

	tests := []struct {
		name    string
		state   game.State
		wantErr string
	}{
		{
			name:    "below MinEndDay",
			state:   game.State{Day: 9},
			wantErr: `strategy "test_strategy" seed 42: day 9 below min 10`,
		},
		{
			name:    "wrong end day",
			state:   game.State{Day: 13, GameOver: false, Won: false, BeaconParts: 3},
			wantErr: `strategy "test_strategy" seed 42: day 13 want 14`,
		},
		{
			name:    "won mismatch",
			state:   game.State{Day: 14, GameOver: false, Won: true, BeaconParts: 3},
			wantErr: `strategy "test_strategy" seed 42: won true want false`,
		},
		{
			name:    "beacon below MinBeaconParts",
			state:   game.State{Day: 14, GameOver: false, Won: false, BeaconParts: 1},
			wantErr: `strategy "test_strategy" seed 42: beacon 1 below min 2`,
		},
		{
			name:    "beacon below expected",
			state:   game.State{Day: 14, GameOver: false, Won: false, BeaconParts: 2},
			wantErr: `strategy "test_strategy" seed 42: beacon 2 want 3`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := game.CheckBaselineOutcome(strategy, seed, tc.state)
			if err == nil {
				t.Fatal("game.CheckBaselineOutcome: want error, got nil")
			}
			if got := err.Error(); got != tc.wantErr {
				t.Fatalf("error = %q, want %q", got, tc.wantErr)
			}
		})
	}
}
