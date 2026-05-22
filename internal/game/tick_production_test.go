package game

import "testing"

func TestAdvanceDay_HydroAndSolarProvideDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Power = 50
	s.Food = 50

	s.advanceDay()

	wantFood := 50 + 6 - DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (production before upkeep)", s.Food, wantFood)
	}
	wantPower := 50 + 6 - DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (production before upkeep)", s.Power, wantPower)
	}
}
