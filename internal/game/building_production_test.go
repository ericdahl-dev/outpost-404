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
		"solar_array":  {"power": 6},
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

	// +6 daily food then -8 upkeep (pop 8) => net -2 from food 10 => 8
	if s.Food != 8 {
		t.Fatalf("Food = %d, want 8 (production before upkeep)", s.Food)
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

	// +12 food/day (6*2), -8 upkeep => +4/day for 3 days => 62
	if s.Food != 62 {
		t.Fatalf("Food = %d, want 62 after 3 days of hydro L2 production", s.Food)
	}
	if s.Day != 4 {
		t.Fatalf("Day = %d, want 4", s.Day)
	}
}

func TestSolarDailyProduction_AccumulatesPowerOverMultipleDays(t *testing.T) {
	s := newTestState()
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	for range 3 {
		s.advanceDay()
	}

	// +6 power/day, -10 upkeep (pop 8) => -4/day => 50-12 = 38
	if s.Power != 38 {
		t.Fatalf("Power = %d, want 38 after 3 days of solar production", s.Power)
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

	if s.Food != 46 || s.Power != 46 {
		t.Fatalf("Food=%d Power=%d, want 46/46", s.Food, s.Power)
	}
}
