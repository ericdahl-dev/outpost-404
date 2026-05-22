package game

func testContent() Content {
	return Content{
		Buildings: []BuildingDef{
			{
				ID:       "solar_array",
				Name:     "Solar Array",
				Cost:     70,
				MaxLevel: 3,
				Effects:  map[string]int{"power": 20},
			},
			{
				ID:       "hydroponics",
				Name:     "Hydroponics",
				Cost:     65,
				MaxLevel: 3,
				Effects:  map[string]int{"food": 22},
			},
		},
		Events: nil,
	}
}

func newTestState() State {
	return NewState(testContent())
}
