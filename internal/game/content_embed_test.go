package game

import "testing"

func TestLoadEmbeddedContent_LoadsBuildingsAndEvents(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	if len(content.Buildings) == 0 || len(content.Events) == 0 {
		t.Fatal("expected embedded buildings and events")
	}
}
