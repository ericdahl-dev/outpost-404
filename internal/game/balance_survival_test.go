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
