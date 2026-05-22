package game_test

import (
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func TestLoadEmbeddedContent_BuildingsHaveDailyEffectsWhereExpected(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}
	wantDaily := map[string]map[string]int{
		"hydroponics": {"food": 6},
		"solar_array": {"power": 6},
	}
	for id, want := range wantDaily {
		def, ok := content.FindBuilding(id)
		if !ok {
			t.Fatalf("missing building %q", id)
		}
		if len(def.DailyEffects) != len(want) {
			t.Fatalf("%s dailyEffects = %v, want %v", id, def.DailyEffects, want)
		}
		for k, v := range want {
			if def.DailyEffects[k] != v {
				t.Fatalf("%s dailyEffects[%q] = %d, want %d", id, k, def.DailyEffects[k], v)
			}
		}
	}
}

func TestAdvanceDay_AppliesBuildingProductionBeforeUpkeep(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 1}
	s.Food = 10
	s.Power = 50
	s.Morale = 50

	game.AdvanceDayForTest(&s)

	want := 10 + 6 - game.DailyFoodUpkeep(s.Population)
	if s.Food != want {
		t.Fatalf("Food = %d, want %d (production before upkeep)", s.Food, want)
	}
}

func TestHydroDailyProduction_AccumulatesOverMultipleDays(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 2}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	for range 3 {
		game.AdvanceDayForTest(&s)
	}

	net := 12 - game.DailyFoodUpkeep(s.Population)
	want := 50 + net*3
	if s.Food != want {
		t.Fatalf("Food = %d, want %d after 3 days of hydro L2 production", s.Food, want)
	}
	if s.Day != 4 {
		t.Fatalf("Day = %d, want 4", s.Day)
	}
}

func TestSolarDailyProduction_AccumulatesPowerOverMultipleDays(t *testing.T) {
	s := newTestState()
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 2}
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	for range 3 {
		game.AdvanceDayForTest(&s)
	}

	net := 12 - game.DailyPowerUpkeep(s.Population)
	want := 50 + net*3
	if s.Power != want {
		t.Fatalf("Power = %d, want %d after 3 days of solar L2 production", s.Power, want)
	}
}

func TestAdvanceDay_Level0Facility_NoDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 0}
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 0}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	game.AdvanceDayForTest(&s)

	wantFood := 50 - game.DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (level 0 produces nothing)", s.Food, wantFood)
	}
	wantPower := 50 - game.DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (level 0 produces nothing)", s.Power, wantPower)
	}
}

func TestAdvanceDay_MissingFacility_NoDailyProduction(t *testing.T) {
	s := newTestState()
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	game.AdvanceDayForTest(&s)

	wantFood := 50 - game.DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (missing building produces nothing)", s.Food, wantFood)
	}
	wantPower := 50 - game.DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (missing building produces nothing)", s.Power, wantPower)
	}
}

func TestNextDay_HydroLevel2_ScalesDailyFoodProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 2}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	for range 3 {
		s.NextDay()
	}

	net := 12 - game.DailyFoodUpkeep(s.Population)
	want := 50 + net*3
	if s.Food != want {
		t.Fatalf("Food = %d, want %d after 3 NextDay with hydro L2", s.Food, want)
	}
	if s.Day != 4 {
		t.Fatalf("Day = %d, want 4", s.Day)
	}
}

func TestNextDay_SolarLevel3_ScalesDailyPowerProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 3}
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	for range 3 {
		s.NextDay()
	}

	net := 18 - game.DailyPowerUpkeep(s.Population)
	want := 50 + net*3
	if s.Power != want {
		t.Fatalf("Power = %d, want %d after 3 NextDay with solar L3", s.Power, want)
	}
}

func TestFormatDailyProductionNote(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("game.LoadEmbeddedContent: %v", err)
	}
	var hydro game.BuildingDef
	for _, def := range content.Buildings {
		if def.ID == "hydroponics" {
			hydro = def
			break
		}
	}
	got := game.FormatDailyProductionNote(hydro)
	if got != "Daily: +6 food/lv" {
		t.Fatalf("game.FormatDailyProductionNote = %q, want Daily: +6 food/lv", got)
	}
}

func TestApplyBuildingProduction_UsesContentOrder(t *testing.T) {
	content, err := game.LoadContent(game.DataDir())
	if err != nil {
		t.Fatalf("game.LoadContent: %v", err)
	}
	s := game.NewState(content)
	s.Buildings["hydroponics"] = game.Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 1}
	s.Food = 40
	s.Power = 40

	game.ApplyBuildingProductionForTest(&s)

	if s.Food != 46 || s.Power != 46 {
		t.Fatalf("Food=%d Power=%d, want 46/46", s.Food, s.Power)
	}
}
