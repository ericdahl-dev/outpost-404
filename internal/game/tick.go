package game

type nextDayOutcome struct {
	eventID          string
	populationGrowth bool
}

// nextDayWithDetail advances the calendar and resolves the day's event.
// Live play passes nil replayDetail (random event). Replay passes logged detail with event_id.
func (s *State) nextDayWithDetail(replayDetail map[string]any) nextDayOutcome {
	popBefore := s.Population
	s.advanceDay()

	var eventID string
	if replayDetail != nil {
		if id, _ := replayDetail["event_id"].(string); id != "" {
			eventID = id
			s.applyEventByID(id)
		}
	}
	if eventID == "" {
		eventID = s.TriggerRandomEvent()
	}
	s.Clamp()

	return nextDayOutcome{
		eventID:          eventID,
		populationGrowth: s.Population > popBefore,
	}
}

func (s *State) NextDay() {
	if s.GameOver {
		return
	}

	s.doAction("next_day", nil, func(detail map[string]any) {
		outcome := s.nextDayWithDetail(nil)
		if outcome.eventID != "" {
			detail["event_id"] = outcome.eventID
		}
		if outcome.populationGrowth {
			detail["population_growth"] = true
		}
	})
}

// replayNextDay applies a logged day advance using the recorded event_id when present.
func (s *State) replayNextDay(detail map[string]any) {
	if s.GameOver {
		return
	}
	s.nextDayWithDetail(detail)
}

// applyBuildingProduction grants per-day output from building dailyEffects (JSON order).
// Runs at the start of advanceDay, before resource upkeep and morale drift.
func (s *State) applyBuildingProduction() {
	for _, def := range s.Content.Buildings {
		b, ok := s.Buildings[def.ID]
		if !ok || b.Level <= 0 || len(def.DailyEffects) == 0 {
			continue
		}
		if effects := dailyEffectsScaled(def, b); len(effects) > 0 {
			s.applyEffects(effects, 1)
		}
	}
}

func (s *State) TriggerRandomEvent() string {
	s.ensureRNG()
	if !s.randomEventRollOccurs(s.rng.Intn(RandomEventRollSides)) {
		return ""
	}

	candidates := eligibleEventsForState(*s, s.Content.Events)
	event, ok := s.pickRandomEligibleEvent(candidates)
	if !ok {
		return ""
	}
	s.applyEvent(event)
	return event.ID
}
