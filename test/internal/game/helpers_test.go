package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
)

func testContent() game.Content {
	return game.Content{
		Buildings: []game.BuildingDef{
			{
				ID:           "solar_array",
				Name:         "Solar Array",
				Cost:         70,
				MaxLevel:     3,
				Effects:      map[string]int{"power": 20},
				DailyEffects: map[string]int{"power": 6},
			},
			{
				ID:           "hydroponics",
				Name:         "Hydroponics",
				Cost:         65,
				MaxLevel:     3,
				Effects:      map[string]int{"food": 22},
				DailyEffects: map[string]int{"food": 6},
			},
		},
		Events: nil,
	}
}

func newTestState() game.State {
	return game.NewState(testContent())
}
