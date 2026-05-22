package game

import "testing"

func TestProjectedDailyDeltas_UpkeepWithoutBuildings(t *testing.T) {
	s := newTestState()
	s.Power = 30
	s.Food = 30
	s.Population = 8
	d := s.ProjectedDailyDeltas()
	wantPower := -DailyPowerUpkeep(8)
	wantFood := -DailyFoodUpkeep(8)
	if d.Power != wantPower {
		t.Fatalf("Power delta = %d, want %d", d.Power, wantPower)
	}
	if d.Food != wantFood {
		t.Fatalf("Food delta = %d, want %d", d.Food, wantFood)
	}
	if d.Morale != -StressMoraleLoss {
		t.Fatalf("Morale delta = %d, want stress loss %d", d.Morale, -StressMoraleLoss)
	}
}

func TestProjectedDailyDeltas_IncludesBuildingProduction(t *testing.T) {
	s := newTestState()
	s.Population = 8
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	d := s.ProjectedDailyDeltas()
	wantPower := 10 - DailyPowerUpkeep(8)
	if d.Power != wantPower {
		t.Fatalf("Power delta = %d, want %d", d.Power, wantPower)
	}
}

func TestProjectedDailyDeltas_DamagedHalvesProduction(t *testing.T) {
	s := newTestState()
	s.Population = 8
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2, Damaged: true}
	d := s.ProjectedDailyDeltas()
	wantFood := 6 - DailyFoodUpkeep(8)
	if d.Food != wantFood {
		t.Fatalf("Food delta = %d, want %d (half of 12 - upkeep)", d.Food, wantFood)
	}
}

func TestProjectedDailyDeltas_ComfortableMoraleGain(t *testing.T) {
	s := newTestState()
	s.Power = ComfortPowerMin + 1
	s.Food = ComfortFoodMin + 1
	d := s.ProjectedDailyDeltas()
	if d.Morale != ComfortMoraleGain {
		t.Fatalf("Morale delta = %d, want %+d", d.Morale, ComfortMoraleGain)
	}
}
