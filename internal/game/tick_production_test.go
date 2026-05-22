package game

import "testing"

func TestAdvanceDay_HydroAndSolarProvideDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Power = 50
	s.Food = 50

	s.advanceDay()

	if s.Food != 48 {
		t.Fatalf("Food = %d, want 48 (50 + 6 hydro - 8 upkeep)", s.Food)
	}
	if s.Power != 46 {
		t.Fatalf("Power = %d, want 46 (50 + 6 solar - 10 upkeep)", s.Power)
	}
}
