package ui

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func TestFormatDailyDelta(t *testing.T) {
	if got := formatDailyDelta(6); got != "+6/day" {
		t.Fatalf("got %q", got)
	}
	if got := formatDailyDelta(-5); got != "-5/day" {
		t.Fatalf("got %q", got)
	}
}

func TestResourceBarLine_IncludesDeltaAndCritical(t *testing.T) {
	line := resourceBarLine("Food", 14, -5, 32, game.WarningFoodCriticalAt, game.WarningFoodUrgentAt)
	if !strings.Contains(line, "14%") {
		t.Fatalf("missing value: %q", line)
	}
	if !strings.Contains(line, "-5/day") {
		t.Fatalf("missing delta: %q", line)
	}
	if !strings.Contains(line, "CRITICAL") {
		t.Fatalf("missing CRITICAL tag: %q", line)
	}
}

func TestFormatStatusStrip_JoinsBadges(t *testing.T) {
	s := game.NewState(game.Content{Buildings: []game.BuildingDef{{ID: "habitat", Name: "Habitat", MaxLevel: 3}}})
	s.Food = game.WarningFoodCriticalAt
	strip := formatStatusStrip(game.StatusBadges(s))
	if !strings.HasPrefix(strip, "Status: ") {
		t.Fatalf("unexpected strip: %q", strip)
	}
	if !strings.Contains(strip, "FOOD CRITICAL") {
		t.Fatalf("missing badge: %q", strip)
	}
}

func TestRenderResourcePanel_IncludesStatusStrip(t *testing.T) {
	s := game.NewState(game.Content{Buildings: []game.BuildingDef{{ID: "habitat", Name: "Habitat", MaxLevel: 3}}})
	s.Food = game.WarningFoodCriticalAt
	body := RenderResourcePanel(s, 36)
	if !strings.Contains(body, "Status:") {
		t.Fatalf("missing status strip:\n%s", body)
	}
	if !strings.Contains(body, "Power") {
		t.Fatalf("missing power bar:\n%s", body)
	}
}
