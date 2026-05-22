package ui_test

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

func TestMainLayout_wideAt120x35(t *testing.T) {
	l := ui.MainLayoutFor(120, 35)
	if l.Mode != ui.LayoutModeWide {
		t.Fatalf("mode = %v, want wide", l.Mode)
	}
	if l.Stacked || l.OutpostBelow {
		t.Fatalf("wide should be 3-column row: stacked=%v outpostBelow=%v", l.Stacked, l.OutpostBelow)
	}
}

func TestMainLayout_mediumAt100x30(t *testing.T) {
	l := ui.MainLayoutFor(100, 30)
	if l.Mode != ui.LayoutModeMedium {
		t.Fatalf("mode = %v, want medium", l.Mode)
	}
	if !l.OutpostBelow || l.Stacked {
		t.Fatalf("medium want outpost below: stacked=%v outpostBelow=%v", l.Stacked, l.OutpostBelow)
	}
}

func TestMainLayout_narrowAt80x24(t *testing.T) {
	l := ui.MainLayoutFor(80, 24)
	if l.Mode != ui.LayoutModeNarrow {
		t.Fatalf("mode = %v, want narrow", l.Mode)
	}
	if !l.Stacked {
		t.Fatal("narrow should stack panels")
	}
}

func TestRenderStatusHeader_IncludesDayAndScenario(t *testing.T) {
	s := game.NewState(game.Content{})
	s.Day = 18
	s.ScenarioID = "standard"
	s.DifficultyID = "normal"
	p := game.RunProfiles{
		Scenarios:    []game.ScenarioDef{{ID: "standard", Name: "Standard Landing"}},
		Difficulties: []game.DifficultyDef{{ID: "normal", Name: "Normal"}},
	}
	h := ui.RenderStatusHeader(s, p, 100)
	for _, want := range []string{"Outpost 404", "Day 18", "Standard Landing", "Normal"} {
		if !strings.Contains(h, want) {
			t.Fatalf("missing %q in %q", want, h)
		}
	}
}
