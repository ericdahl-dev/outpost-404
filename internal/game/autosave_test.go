package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoad_RoundTripPreservesSimulationState(t *testing.T) {
	profiles := testRunProfiles()
	s := NewRun(testContent(), profiles, 42, "first_landing", "hard")
	s.Build("solar_array")
	s.DamageBuilding("solar_array")
	s.NextDay()
	s.RepairBuilding("solar_array")

	path := filepath.Join(t.TempDir(), "autosave.json")
	if err := SaveAutosave(&s, path); err != nil {
		t.Fatalf("SaveAutosave: %v", err)
	}

	loaded, err := LoadAutosave(path, s.Content, profiles)
	if err != nil {
		t.Fatalf("LoadAutosave: %v", err)
	}

	if loaded.Day != s.Day || loaded.Credits != s.Credits || loaded.Power != s.Power {
		t.Fatalf("vitals mismatch: loaded day=%d credits=%d power=%d", loaded.Day, loaded.Credits, loaded.Power)
	}
	if loaded.ScenarioID != "first_landing" || loaded.DifficultyID != "hard" {
		t.Fatalf("run profile = %s/%s", loaded.ScenarioID, loaded.DifficultyID)
	}
	b := loaded.Buildings["solar_array"]
	if b.Level != 1 || b.Damaged {
		t.Fatalf("solar_array = %+v, want L1 repaired", b)
	}
	if len(loaded.Log) != len(s.Log) {
		t.Fatalf("log len %d want %d", len(loaded.Log), len(s.Log))
	}
}

func TestSaveLoad_RoundTripPreservesRNGDraws(t *testing.T) {
	profiles := testRunProfiles()
	a := NewRun(testContent(), profiles, 99, "standard", "normal")
	b := NewRun(testContent(), profiles, 99, "standard", "normal")
	for range 5 {
		a.NextDay()
		b.NextDay()
	}
	path := filepath.Join(t.TempDir(), "save.json")
	if err := SaveAutosave(&a, path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadAutosave(path, a.Content, profiles)
	if err != nil {
		t.Fatal(err)
	}
	loaded.NextDay()
	b.NextDay()
	if diff := snapshotDiff(normalizeSnapshot(loaded.snapshot()), normalizeSnapshot(b.snapshot())); diff != "" {
		t.Fatalf("RNG diverged after load: %s", diff)
	}
}

func TestLoadAutosave_missingFile(t *testing.T) {
	_, err := LoadAutosave(filepath.Join(t.TempDir(), "nope.json"), testContent(), testRunProfiles())
	if err == nil {
		t.Fatal("expected error for missing save")
	}
}

func TestLoadAutosave_invalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadAutosave(path, testContent(), testRunProfiles())
	if err == nil {
		t.Fatal("expected parse error")
	}
}
