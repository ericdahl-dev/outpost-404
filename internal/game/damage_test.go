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

func TestDamageEvent_SolarArrayDamagedViaApplyEventByID(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	s := NewState(content)
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 2}
	s.Day = 10

	// Find an event that damages solar_array
	var damageEventID string
	for _, e := range s.Content.Events {
		if e.DamageBuilding == "solar_array" {
			damageEventID = e.ID
			break
		}
	}
	if damageEventID == "" {
		t.Skip("no solar_array damage event in events.json — add one first")
	}

	s.applyEventByID(damageEventID)
	if !s.Buildings["solar_array"].Damaged {
		t.Fatal("solar_array should be Damaged after damage event")
	}

	// Daily output should be halved
	def, _ := s.FindBuilding("solar_array")
	effects := dailyEffectsScaled(def, s.Buildings["solar_array"])
	b := s.Buildings["solar_array"]
	b.Damaged = false
	fullEffects := dailyEffectsScaled(def, b)
	if effects["power"] != fullEffects["power"]/2 {
		t.Fatalf("damaged power %d want %d (half of %d)", effects["power"], fullEffects["power"]/2, fullEffects["power"])
	}
}

func TestDamageEvent_RandomBuilt_DamagesBuiltFacility(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	s := NewState(content)
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Day = 20

	var randDamageEventID string
	for _, e := range s.Content.Events {
		if e.DamageRandomBuilt {
			randDamageEventID = e.ID
			break
		}
	}
	if randDamageEventID == "" {
		t.Skip("no damageRandomBuilt event in events.json — add one first")
	}

	s.applyEventByID(randDamageEventID)
	// at least one built building should be damaged
	anyDamaged := false
	for _, b := range s.Buildings {
		if b.Damaged {
			anyDamaged = true
			break
		}
	}
	if !anyDamaged {
		t.Fatal("expected at least one building damaged after damageRandomBuilt event")
	}
}
