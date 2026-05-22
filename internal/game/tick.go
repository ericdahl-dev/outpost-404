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
	s.Power -= 6 + s.Population/2
	s.Food -= 4 + s.Population/2
	s.Credits += 18

	if s.Power > 50 && s.Food > 40 {
		s.Morale += 2
	} else {
		s.Morale -= 3
	}

	if s.Day%5 == 0 && s.Population < s.PopulationCap && s.Food > 35 && s.Morale > 40 {
		s.Population++
		s.AddLog("A new colonist joined after hearing your beacon tests.")
	}
}

// applyBuildingProduction grants per-day output from building dailyEffects (JSON order).
// Runs at the start of advanceDay, before resource upkeep and morale drift.
func (s *State) applyBuildingProduction() {
	for _, def := range s.Content.Buildings {
		lvl := s.BuildingLevel(def.ID)
		if lvl <= 0 || len(def.DailyEffects) == 0 {
			continue
		}
		s.applyEffects(def.DailyEffects, lvl)
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
