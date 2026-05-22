package ui_test

import (
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

func TestTitleView_OmitsContinueWithoutSave(t *testing.T) {
	m := ui.NewModel(game.Content{}, game.RunProfiles{})
	m.CanContinue = false
	v := m.View()
	if strings.Contains(v, "Continue") {
		t.Fatalf("view should not offer continue: %q", v)
	}
	if !strings.Contains(v, "Start new run") {
		t.Fatal("expected start option")
	}
}

func TestTitleView_ShowsContinueWhenSavePresent(t *testing.T) {
	m := ui.NewModel(game.Content{}, game.RunProfiles{})
	m.CanContinue = true
	v := m.View()
	if !strings.Contains(v, "Continue") {
		t.Fatalf("view missing continue: %q", v)
	}
}
