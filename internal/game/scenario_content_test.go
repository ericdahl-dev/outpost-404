package game

import "testing"

func TestLoadEmbeddedRunProfiles_HasFiveScenariosAndThreeDifficulties(t *testing.T) {
	p, err := LoadEmbeddedRunProfiles()
	if err != nil {
		t.Fatalf("LoadEmbeddedRunProfiles: %v", err)
	}
	if len(p.Scenarios) < 5 {
		t.Fatalf("scenarios = %d, want >= 5", len(p.Scenarios))
	}
	if len(p.Difficulties) != 3 {
		t.Fatalf("difficulties = %d, want 3", len(p.Difficulties))
	}
	if _, ok := p.FindScenario("beacon_rush"); !ok {
		t.Fatal("missing beacon_rush scenario")
	}
}
