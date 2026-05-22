package game_test

import (
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func TestLoadContent_LoadsBuildingsAndEventsFromDataDir(t *testing.T) {
	dir := game.DataDir()

	content, err := game.LoadContent(dir)
	if err != nil {
		t.Fatalf("game.LoadContent: %v", err)
	}
	if len(content.Buildings) == 0 {
		t.Fatal("expected buildings from data/buildings.json")
	}
	if len(content.Events) == 0 {
		t.Fatal("expected events from data/events.json")
	}

	for _, b := range content.Buildings {
		if b.ID == "solar_array" {
			if b.Cost != 70 {
				t.Fatalf("solar_array cost = %d, want 70", b.Cost)
			}
			if b.DailyEffects["power"] != 6 {
				t.Fatalf("solar_array daily power = %d, want 6", b.DailyEffects["power"])
			}
			return
		}
	}
	t.Fatal("solar_array not found in loaded buildings")
}
