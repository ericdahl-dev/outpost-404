package game

import "testing"

func TestAdvanceDay_HydroAndSolarProvideDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Power = 50
	s.Food = 50

	s.advanceDay()

	if s.Food != 46 {
		t.Fatalf("Food = %d, want 46 (50 + 4 hydro - 8 upkeep)", s.Food)
	}
	if s.Power != 44 {
		t.Fatalf("Power = %d, want 44 (50 + 5 solar - 11 upkeep)", s.Power)
	}
}
