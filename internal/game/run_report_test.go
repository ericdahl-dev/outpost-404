package game

import "testing"

func TestRecordVitalLows_TracksMinimums(t *testing.T) {
	s := newTestState()
	s.recordVitalLows()
	s.Power = 12
	s.Food = 40
	s.recordVitalLows()
	if s.MinPowerSeen != 12 {
		t.Fatalf("MinPowerSeen = %d, want 12", s.MinPowerSeen)
	}
	if s.MinFoodSeen != 40 {
		t.Fatalf("MinFoodSeen = %d, want 40", s.MinFoodSeen)
	}
}

func TestEndCause_BeaconWin(t *testing.T) {
	s := newTestState()
	s.BeaconParts = 5
	s.MaxBeaconParts = 5
	s.CheckEnd()
	if s.EndCause() != "Signal Beacon complete — rescue inbound" {
		t.Fatalf("got %q", s.EndCause())
	}
}

func TestEndCause_CollapseFood(t *testing.T) {
	s := newTestState()
	s.Food = 0
	s.CheckEnd()
	if s.EndCause() != "Food reserves depleted" {
		t.Fatalf("got %q", s.EndCause())
	}
}

func TestBuildRunReport_IncludesProfileAndLows(t *testing.T) {
	s := newTestState()
	s.ScenarioID = "standard"
	s.DifficultyID = "normal"
	s.Seed = 4242
	s.MinPowerSeen = 10
	s.Power = 0
	s.Food = 0
	s.CheckEnd()
	p := RunProfiles{
		Scenarios:    []ScenarioDef{{ID: "standard", Name: "Standard Landing"}},
		Difficulties: []DifficultyDef{{ID: "normal", Name: "Normal"}},
	}
	r := BuildRunReport(s, p)
	if r.ScenarioName != "Standard Landing" || r.DifficultyName != "Normal" {
		t.Fatalf("profile names: %+v", r)
	}
	if r.Seed != 4242 {
		t.Fatalf("seed = %d", r.Seed)
	}
	if r.LowestPower != 10 {
		t.Fatalf("LowestPower = %d", r.LowestPower)
	}
	if r.Won {
		t.Fatal("expected loss report")
	}
}
