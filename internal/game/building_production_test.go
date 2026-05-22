package game

import (
	"path/filepath"
	"testing"
)

func TestLoadEmbeddedContent_BuildingsHaveDailyEffectsWhereExpected(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	wantDaily := map[string]map[string]int{
		"hydroponics":  {"food": 6},
		"solar_array":  {"power": 10},
		"habitat":      {"morale": 1},
		"workshop":     {"morale": 1},
		"radio_tower":  {"credits": 2, "morale": 1},
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
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 1}
	s.Food = 10
	s.Power = 50
	s.Morale = 50

	s.advanceDay()

	want := 10 + 6 - DailyFoodUpkeep(s.Population)
	if s.Food != want {
		t.Fatalf("Food = %d, want %d (production before upkeep)", s.Food, want)
	}
}

func TestHydroDailyProduction_ConsumesPowerPerLevel(t *testing.T) {
	content := testContent()
	for i := range content.Buildings {
		if content.Buildings[i].ID == "hydroponics" {
			content.Buildings[i].DailyEffects = map[string]int{"food": 6, "power": -1}
			break
		}
	}
	s := NewState(content)
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	s.advanceDay()

	foodNet := 12 - DailyFoodUpkeep(s.Population)
	powerNet := -2 - DailyPowerUpkeep(s.Population)
	if s.Food != 50+foodNet || s.Power != 50+powerNet {
		t.Fatalf("Food=%d Power=%d, want %d/%d", s.Food, s.Power, 50+foodNet, 50+powerNet)
	}
}

func TestHydroDailyProduction_AccumulatesOverMultipleDays(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	for range 3 {
		s.advanceDay()
	}

	net := 12 - DailyFoodUpkeep(s.Population)
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
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 2}
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	for range 3 {
		s.advanceDay()
	}

	net := 20 - DailyPowerUpkeep(s.Population)
	want := 50 + net*3
	if s.Power != want {
		t.Fatalf("Power = %d, want %d after 3 days of solar L2 production", s.Power, want)
	}
}

func TestAdvanceDay_Level0Facility_NoDailyProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 0}
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 0}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	s.advanceDay()

	wantFood := 50 - DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (level 0 produces nothing)", s.Food, wantFood)
	}
	wantPower := 50 - DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (level 0 produces nothing)", s.Power, wantPower)
	}
}

func TestAdvanceDay_MissingFacility_NoDailyProduction(t *testing.T) {
	s := newTestState()
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	s.advanceDay()

	wantFood := 50 - DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d (missing building produces nothing)", s.Food, wantFood)
	}
	wantPower := 50 - DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d (missing building produces nothing)", s.Power, wantPower)
	}
}

func TestNextDay_HydroLevel2_ScalesDailyFoodProduction(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2}
	s.Food = 50
	s.Power = 50
	s.Morale = 50

	for range 3 {
		s.NextDay()
	}

	net := 12 - DailyFoodUpkeep(s.Population)
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
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 3}
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	for range 3 {
		s.advanceDay()
	}

	net := 30 - DailyPowerUpkeep(s.Population)
	want := 50 + net*3
	if s.Power != want {
		t.Fatalf("Power = %d, want %d after 3 advanceDay with solar L3", s.Power, want)
	}
}

func TestFormatDailyProductionNote(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	var hydro BuildingDef
	for _, def := range content.Buildings {
		if def.ID == "hydroponics" {
			hydro = def
			break
		}
	}
	got := FormatDailyProductionNote(hydro)
	if got != "Daily: +6 food/lv" {
		t.Fatalf("FormatDailyProductionNote = %q, want Daily: +6 food/lv", got)
	}
}

func TestWorkshopDailyProduction_AddsMoralePerLevel(t *testing.T) {
	s := newTestState()
	s.Buildings["workshop"] = Building{DefID: "workshop", Level: 2}
	s.Morale = 40
	s.Power = 80
	s.Food = 80

	s.advanceDay()

	want := 40 + 2 + ComfortMoraleGain
	if s.Morale != want {
		t.Fatalf("Morale = %d, want %d", s.Morale, want)
	}
}

func TestHabitatDailyProduction_AddsMoralePerLevel(t *testing.T) {
	s := newTestState()
	s.Buildings["habitat"] = Building{DefID: "habitat", Level: 2}
	s.Morale = 50
	s.Power = 80
	s.Food = 80

	s.advanceDay()

	want := 50 + 2 + ComfortMoraleGain
	if s.Morale != want {
		t.Fatalf("Morale = %d, want %d", s.Morale, want)
	}
}

func TestApplyBuildingProduction_UsesContentOrder(t *testing.T) {
	content, err := LoadContent(filepath.Join("..", "..", "data"))
	if err != nil {
		t.Fatalf("LoadContent: %v", err)
	}
	s := NewState(content)
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 1}
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Food = 40
	s.Power = 40

	s.applyBuildingProduction()

	if s.Food != 46 || s.Power != 50 {
		t.Fatalf("Food=%d Power=%d, want 46/50", s.Food, s.Power)
	}
}
