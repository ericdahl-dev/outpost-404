package game

func (s *State) NextDay() {
	if s.GameOver {
		return
	}

	before := s.snapshot()
	popBefore := s.Population

	s.Day++
	s.Power -= 8 + s.Population/2
	s.Food -= 6 + s.Population/2
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

	eventID := s.TriggerRandomEvent()
	s.Clamp()

	detail := map[string]any{}
	if eventID != "" {
		detail["event_id"] = eventID
	}
	if s.Population > popBefore {
		detail["population_growth"] = true
	}
	s.recordAction("next_day", detail, before, s.snapshot())
	s.CheckEnd()
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
