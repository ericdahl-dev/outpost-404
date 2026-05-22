package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"testing"
)

func TestAdvanceDay_HydroAndSolarProvideDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 1}
	s.Power = 50
	s.Food = 50

	game.AdvanceDayForTest(&s)

	wantFood := 50 + 6 - game.DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (production before upkeep)", s.Food, wantFood)
	}
	wantPower := 50 + 6 - game.DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (production before upkeep)", s.Power, wantPower)
	}
}
