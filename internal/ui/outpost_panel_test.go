package ui

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func TestFormatFacilityToken_Unbuilt(t *testing.T) {
	got := formatFacilityToken("SA", facilitySpec{maxLevel: 3})
	if got != "[SA]" {
		t.Fatalf("got %q, want [SA]", got)
	}
}

func TestFormatFacilityToken_Level2(t *testing.T) {
	got := formatFacilityToken("SA", facilitySpec{level: 2, maxLevel: 3})
	if got != "[SA²]" {
		t.Fatalf("got %q, want [SA²]", got)
	}
}

func TestFormatFacilityToken_MaxLevel(t *testing.T) {
	got := formatFacilityToken("SA", facilitySpec{level: 3, maxLevel: 3})
	if got != "[SA★]" {
		t.Fatalf("got %q, want [SA★]", got)
	}
}

func TestFormatFacilityToken_Damaged(t *testing.T) {
	got := formatFacilityToken("WS", facilitySpec{level: 2, maxLevel: 2, damaged: true})
	if got != "[WS!]" {
		t.Fatalf("got %q, want [WS!]", got)
	}
}

func TestFormatFacilityToken_Offline(t *testing.T) {
	got := formatFacilityToken("HY", facilitySpec{level: 1, maxLevel: 3, offline: true})
	if got != "[HY×]" {
		t.Fatalf("got %q, want [HY×]", got)
	}
}

func TestFormatFacilityToken_BeaconProgress(t *testing.T) {
	got := formatFacilityToken("BE", facilitySpec{
		beaconParts: 3, maxBeacon: 5, isBeacon: true,
	})
	if got != "[BE 3/5]" {
		t.Fatalf("got %q, want [BE 3/5]", got)
	}
}

func TestFormatFacilityToken_BeaconComplete(t *testing.T) {
	got := formatFacilityToken("BE", facilitySpec{
		beaconParts: 5, maxBeacon: 5, isBeacon: true,
	})
	if got != "[BE★]" {
		t.Fatalf("got %q, want [BE★]", got)
	}
}

func TestRenderOutpostPanel_ShowsBuiltAndUnbuilt(t *testing.T) {
	s := game.NewState(game.Content{Buildings: []game.BuildingDef{
		{ID: "solar_array", Name: "Solar Array", MaxLevel: 3, DailyEffects: map[string]int{"power": 1}},
		{ID: "habitat", Name: "Habitat", MaxLevel: 3},
	}})
	s.Buildings["solar_array"] = game.Building{DefID: "solar_array", Level: 2}
	s.MaxBeaconParts = 5

	body := RenderOutpostPanel(s, 40)
	if !strings.Contains(body, "[SA²]") {
		t.Fatalf("expected built solar token in:\n%s", body)
	}
	if !strings.Contains(body, "[HB]") {
		t.Fatalf("expected unbuilt habitat in:\n%s", body)
	}
	if !strings.Contains(body, "[BE") {
		t.Fatalf("expected beacon token in:\n%s", body)
	}
}

func TestRenderOutpostPanel_DamagedFacility(t *testing.T) {
	s := game.NewState(game.Content{Buildings: []game.BuildingDef{
		{ID: "workshop", Name: "Workshop", MaxLevel: 2},
	}})
	s.Buildings["workshop"] = game.Building{DefID: "workshop", Level: 1, Damaged: true}

	body := RenderOutpostPanel(s, 40)
	if !strings.Contains(body, "[WS!]") {
		t.Fatalf("expected damaged workshop in:\n%s", body)
	}
}
