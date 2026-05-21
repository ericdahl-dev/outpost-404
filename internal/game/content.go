package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadContent(dir string) (Content, error) {
	buildings, err := loadJSON[[]BuildingDef](filepath.Join(dir, "buildings.json"))
	if err != nil {
		return Content{}, err
	}

	events, err := loadJSON[[]EventDef](filepath.Join(dir, "events.json"))
	if err != nil {
		return Content{}, err
	}

	return Content{Buildings: buildings, Events: events}, nil
}

func loadJSON[T any](path string) (T, error) {
	var value T
	bytes, err := os.ReadFile(path)
	if err != nil {
		return value, fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(bytes, &value); err != nil {
		return value, fmt.Errorf("parse %s: %w", path, err)
	}
	return value, nil
}
