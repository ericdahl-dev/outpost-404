package game

import (
	"path/filepath"
	"testing"
)

func TestLoadContent_LoadsBuildingsAndEventsFromDataDir(t *testing.T) {
	dir := filepath.Join("..", "..", "data")

	content, err := LoadContent(dir)
	if err != nil {
		t.Fatalf("LoadContent: %v", err)
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
			return
		}
	}
	t.Fatal("solar_array not found in loaded buildings")
}
