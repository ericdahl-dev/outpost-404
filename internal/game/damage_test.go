package game

import "testing"

func TestDamageBuilding_MarksDamaged(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2}

	s.damageBuilding("hydroponics")

	b := s.Buildings["hydroponics"]
	if !b.Damaged {
		t.Fatal("expected hydroponics damaged")
	}
}

func TestApplyBuildingProduction_DamagedRunsHalfDailyEffects(t *testing.T) {
	s := newTestState()
	s.Buildings["hydroponics"] = Building{DefID: "hydroponics", Level: 2, Damaged: true}
	s.Food = 50
	s.Power = 50

	s.applyBuildingProduction()

	if s.Food != 56 {
		t.Fatalf("Food = %d, want 56 (half of +12)", s.Food)
	}
}

func TestApplyEvent_DamageBuildingMarksFacility(t *testing.T) {
	content := testContentWithEvents()
	content.Events = []EventDef{{
		ID:             "storm",
		Title:          "Storm",
		Description:    "Hit arrays.",
		Effects:        map[string]int{"power": -5},
		MinDay:         1,
		DamageBuilding: "solar_array",
	}}
	s := NewState(content)
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}

	s.applyEvent(content.Events[0])

	if !s.Buildings["solar_array"].Damaged {
		t.Fatal("expected solar_array damaged from event")
	}
}

func TestRepairBuilding_ClearsDamageAndChargesByLevel(t *testing.T) {
	s := newTestState()
	s.Buildings["workshop"] = Building{DefID: "workshop", Level: 2, Damaged: true}
	s.Credits = 100
	wantCost := RepairCost(2)

	s.RepairBuilding("workshop")

	b := s.Buildings["workshop"]
	if b.Damaged {
		t.Fatal("expected damage cleared")
	}
	if s.Credits != 100-wantCost {
		t.Fatalf("Credits = %d, want %d", s.Credits, 100-wantCost)
	}
}
