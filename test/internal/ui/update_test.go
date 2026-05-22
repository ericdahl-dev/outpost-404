package ui_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

func TestUpdate_WindowSizeBeforeStart_DoesNotPanic(t *testing.T) {
	content, err := game.LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	profiles, err := game.LoadEmbeddedRunProfiles()
	if err != nil {
		t.Fatalf("LoadEmbeddedRunProfiles: %v", err)
	}
	m := ui.NewModel(content, profiles)
	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	if _, ok := next.(ui.Model); !ok {
		t.Fatal("expected Model return type")
	}
}
