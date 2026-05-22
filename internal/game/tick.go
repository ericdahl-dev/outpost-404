package game

func (s *State) NextDay() {
	if s.GameOver {
		return
	}

	popBefore := s.Population

	s.doAction("next_day", nil, func(detail map[string]any) {
		s.advanceDay()
		eventID := s.TriggerRandomEvent()
		s.Clamp()

		if eventID != "" {
			detail["event_id"] = eventID
		}
		if s.Population > popBefore {
			detail["population_growth"] = true
		}
	})
}

// replayNextDay applies a logged day advance using the recorded event_id when present.
func (s *State) replayNextDay(detail map[string]any) {
	if s.GameOver {
		return
	}
	s.advanceDay()
	if eventID, _ := detail["event_id"].(string); eventID != "" {
		s.applyEventByID(eventID)
	}
	s.Clamp()
}

func (s *State) advanceDay() {
	s.Day++
	s.applyBuildingProduction()
	s.Power -= 7 + s.Population/2
	s.Food -= 4 + s.Population/2
	s.Credits += 18

	if s.Power > 55 && s.Food > 45 {
		s.Morale += 2
	} else {
		s.Morale -= 5
	}

	if s.Day%5 == 0 && s.Population < s.PopulationCap && s.Food > 35 && s.Morale > 40 {
		s.Population++
		s.AddLog("A new colonist joined after hearing your beacon tests.")
	}
}

// applyBuildingProduction grants daily output from completed facilities (see data/buildings.json).
func (s *State) applyBuildingProduction() {
	if lvl := s.BuildingLevel("hydroponics"); lvl > 0 {
		s.Food += 4 * lvl
	}
	if lvl := s.BuildingLevel("solar_array"); lvl > 0 {
		s.Power += 5 * lvl
	}
}

func (s *State) applyEventByID(id string) {
	for _, event := range s.Content.Events {
		if event.ID == id {
			s.applyEffects(event.Effects, 1)
			s.AddLog(event.Title + ": " + event.Description)
			return
		}
	}
}

func (s *State) TriggerRandomEvent() string {
	s.ensureRNG()
	if s.rng.Intn(100) > 45 {
		return ""
	}

	candidates := make([]EventDef, 0)
	for _, event := range s.Content.Events {
		if event.MinDay <= s.Day {
			candidates = append(candidates, event)
		}
	}
	if len(candidates) == 0 {
		return ""
	}

	event := candidates[s.rng.Intn(len(candidates))]
	s.applyEffects(event.Effects, 1)
	s.AddLog(event.Title + ": " + event.Description)
	return event.ID
}
