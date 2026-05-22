package ui

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func TestRenderEndScreen_WinIncludesArtAndSeed(t *testing.T) {
	s := game.NewState(game.Content{})
	s.Won = true
	s.GameOver = true
	s.Seed = 99
	s.ScenarioID = "standard"
	s.DifficultyID = "normal"
	s.BeaconParts = 5
	s.MaxBeaconParts = 5
	p := game.RunProfiles{
		Scenarios:    []game.ScenarioDef{{ID: "standard", Name: "Standard"}},
		Difficulties: []game.DifficultyDef{{ID: "normal", Name: "Normal"}},
	}
	body := RenderEndScreen(game.BuildRunReport(s, p), 80)
	for _, want := range []string{"MISSION COMPLETE", "SIGNAL SENT", "seed 99", "Standard", "Normal"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in:\n%s", want, body)
		}
	}
}

func TestRenderEndScreen_LossIncludesCause(t *testing.T) {
	s := game.NewState(game.Content{})
	s.Food = 0
	s.GameOver = true
	s.Day = 22
	s.MinFoodSeen = 8
	s.MinPowerSeen = 65
	s.MinMoraleSeen = 70
	p := game.RunProfiles{}
	r := game.BuildRunReport(s, p)
	r.Cause = "Food reserves depleted"
	body := RenderEndScreen(r, 70)
	if !strings.Contains(body, "OUTPOST COLLAPSED") {
		t.Fatalf("missing loss title:\n%s", body)
	}
	if !strings.Contains(body, "Food reserves depleted") {
		t.Fatalf("missing cause:\n%s", body)
	}
	if !strings.Contains(body, "lowest food") {
		t.Fatalf("missing vital low:\n%s", body)
	}
}

func TestEndScreenActions_IncludesSameSeed(t *testing.T) {
	if !strings.Contains(endScreenActions(false), "[S]") {
		t.Fatal("expected same-seed action")
	}
}
