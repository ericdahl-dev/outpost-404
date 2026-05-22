// One-off: go run ./scripts/gen-multi-action-fixture.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ericdahl/outpost-404/internal/game"
)

const seed = int64(4242)

func main() {
	content := testContentWithEvents()
	dir := filepath.Join("internal", "game", "testdata", "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}
	path := filepath.Join(dir, "multi_action.jsonl")
	logger, err := game.OpenSessionLog(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open log: %v\n", err)
		os.Exit(1)
	}
	s := game.NewStateWithSeed(content, seed)
	s.SessionLog = logger
	s.LogSessionStart()
	s.Build("solar_array")
	s.Repair()
	s.Trade()
	s.WorkOnBeacon()
	s.NextDay()
	_ = logger.Close()
	fmt.Println("wrote", path)
}

func testContentWithEvents() game.Content {
	c := game.Content{
		Buildings: []game.BuildingDef{
			{ID: "solar_array", Name: "Solar Array", Cost: 70, MaxLevel: 3, Effects: map[string]int{"power": 20}, DailyEffects: map[string]int{"power": 6}},
			{ID: "hydroponics", Name: "Hydroponics", Cost: 65, MaxLevel: 3, Effects: map[string]int{"food": 22}, DailyEffects: map[string]int{"food": 6}},
		},
		Events: []game.EventDef{
			{ID: "quiet_shift", Title: "Quiet Shift", Description: "Calm.", Effects: map[string]int{"morale": 8}, MinDay: 1},
			{ID: "solar_storm", Title: "Solar Storm", Description: "Storm.", Effects: map[string]int{"power": -14}, MinDay: 2},
		},
	}
	return c
}
