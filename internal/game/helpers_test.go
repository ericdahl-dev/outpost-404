package game

func testContent() Content {
	return Content{
		Buildings: []BuildingDef{
			{
				ID:           "solar_array",
				Name:         "Solar Array",
				Cost:         70,
				MaxLevel:     3,
				Effects:      map[string]int{"power": 20},
				DailyEffects: map[string]int{"power": 10},
			},
			{
				ID:           "hydroponics",
				Name:         "Hydroponics",
				Cost:         65,
				MaxLevel:     3,
				Effects:      map[string]int{"food": 22},
				DailyEffects: map[string]int{"food": 6},
			},
			{
				ID:           "habitat",
				Name:         "Habitat",
				Cost:         90,
				MaxLevel:     3,
				Effects:      map[string]int{"populationCap": 3, "morale": 4},
				DailyEffects: map[string]int{"morale": 1},
			},
			{
				ID:           "workshop",
				Name:         "Workshop",
				Cost:         85,
				MaxLevel:     2,
				Effects:      map[string]int{"morale": 10, "credits": 10},
				DailyEffects: map[string]int{"morale": 1},
			},
			{
				ID:           "radio_tower",
				Name:         "Radio Tower",
				Cost:         110,
				MaxLevel: 2,
				Effects:      map[string]int{"credits": 35, "morale": 7},
				DailyEffects: map[string]int{"credits": 2, "morale": 1},
			},
		},
		Events: nil,
	}
}

func newTestState() State {
	return NewState(testContent())
}
