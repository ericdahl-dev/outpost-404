package game

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestReferenceSeeds_Count(t *testing.T) {
	if len(ReferenceSeeds) < 5 {
		t.Fatalf("ReferenceSeeds: got %d, want at least 5", len(ReferenceSeeds))
	}
}

func TestReferenceStrategies_Count(t *testing.T) {
	if len(ReferenceStrategies()) < 4 {
		t.Fatalf("ReferenceStrategies: got %d, want at least 4", len(ReferenceStrategies()))
	}
}

func TestBalanceBaseline_ReferenceStrategies(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}

	for _, strategy := range ReferenceStrategies() {
		t.Run(strategy.ID, func(t *testing.T) {
			path := filepath.Join(ScriptsDir(), strategy.ScriptFile)
			script, err := LoadSimScript(path)
			if err != nil {
				t.Fatalf("LoadSimScript(%s): %v", path, err)
			}

			for _, seed := range ReferenceSeeds {
				t.Run(seedLabel(seed), func(t *testing.T) {
					final, err := Simulate(content, seed, script.Actions)
					if err != nil {
						t.Fatalf("Simulate: %v", err)
					}
					if err := CheckBaselineOutcome(strategy, seed, final); err != nil {
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
	strategy := BaselineStrategy{
		ID:             "test_strategy",
		MinEndDay:      10,
		MinBeaconParts: 2,
		RequireAlive:   true,
		Expected: map[int64]BaselineOutcome{
			seed: {Day: 14, GameOver: false, Won: false, BeaconParts: 3},
		},
	}

	tests := []struct {
		name    string
		state   State
		wantErr string
	}{
		{
			name:    "below MinEndDay",
			state:   State{Day: 9},
			wantErr: `strategy "test_strategy" seed 42: day 9 below min 10`,
		},
		{
			name:    "wrong end day",
			state:   State{Day: 13, GameOver: false, Won: false, BeaconParts: 3},
			wantErr: `strategy "test_strategy" seed 42: day 13 want 14`,
		},
		{
			name:    "won mismatch",
			state:   State{Day: 14, GameOver: false, Won: true, BeaconParts: 3},
			wantErr: `strategy "test_strategy" seed 42: won true want false`,
		},
		{
			name:    "beacon below MinBeaconParts",
			state:   State{Day: 14, GameOver: false, Won: false, BeaconParts: 1},
			wantErr: `strategy "test_strategy" seed 42: beacon 1 below min 2`,
		},
		{
			name:    "beacon below expected",
			state:   State{Day: 14, GameOver: false, Won: false, BeaconParts: 2},
			wantErr: `strategy "test_strategy" seed 42: beacon 2 want 3`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckBaselineOutcome(strategy, seed, tc.state)
			if err == nil {
				t.Fatal("CheckBaselineOutcome: want error, got nil")
			}
			if got := err.Error(); got != tc.wantErr {
				t.Fatalf("error = %q, want %q", got, tc.wantErr)
			}
		})
	}
}
