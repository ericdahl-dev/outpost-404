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

func (s *State) TriggerRandomEvent() string {
	s.ensureRNG()
	if !RandomEventRollOccurs(s.rng.Intn(RandomEventRollSides)) {
		return ""
	}

	candidates := eligibleEvents(s.Content.Events, s.Day)
	event, ok := s.pickRandomEligibleEvent(candidates)
	if !ok {
		return ""
	}
	s.applyEvent(event)
	return event.ID
}
