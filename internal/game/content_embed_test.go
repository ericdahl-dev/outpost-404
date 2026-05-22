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
	for _, def := range content.Buildings {
		if def.ID == "hydroponics" && def.DailyEffects["food"] != 6 {
			t.Fatalf("embedded hydroponics daily food = %d, want 6", def.DailyEffects["food"])
		}
	}
}
