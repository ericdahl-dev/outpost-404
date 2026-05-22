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
