package ui_test

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

func TestMainView_IncludesOutpostSchematic(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	profiles, err := game.LoadEmbeddedRunProfiles()
	if err != nil {
		t.Fatalf("LoadEmbeddedRunProfiles: %v", err)
	}
	m := ui.NewModel(content, profiles)
	m.Started = true
	m.State = game.NewState(content)
	m.TermWidth = 100
	m.TermHeight = 30
	m.Screen = 2 // screenMain
	view := m.View()
	if !strings.Contains(view, "Outpost") {
		t.Fatalf("main view missing Outpost panel:\n%s", view)
	}
	if !strings.Contains(view, "[HB]") {
		t.Fatalf("main view missing habitat schematic token:\n%s", view)
	}
}
