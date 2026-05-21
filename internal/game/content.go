package game

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	outpostdata "github.com/ericdahl/outpost-404/data"
)

// LoadEmbeddedContent loads buildings and events bundled with the binary.
func LoadEmbeddedContent() (Content, error) {
	return LoadContentFS(outpostdata.Files)
}

// LoadContentFS loads JSON content from an fs.FS rooted at buildings.json / events.json.
func LoadContentFS(fsys fs.FS) (Content, error) {
	buildings, err := loadJSONFS[[]BuildingDef](fsys, "buildings.json")
	if err != nil {
		return Content{}, err
	}
	events, err := loadJSONFS[[]EventDef](fsys, "events.json")
	if err != nil {
		return Content{}, err
	}
	return Content{Buildings: buildings, Events: events}, nil
}

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
	return loadJSONFS[T](os.DirFS(filepath.Dir(path)), filepath.Base(path))
}

func loadJSONFS[T any](fsys fs.FS, name string) (T, error) {
	var value T
	bytes, err := fs.ReadFile(fsys, name)
	if err != nil {
		return value, fmt.Errorf("read %s: %w", name, err)
	}
	if err := json.Unmarshal(bytes, &value); err != nil {
		return value, fmt.Errorf("parse %s: %w", name, err)
	}
	return value, nil
}
