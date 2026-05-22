package game

import "fmt"

// applyEvent applies event effects and appends the standard log line.
func (s *State) applyEvent(event EventDef) {
	s.applyEffects(event.Effects, 1)
	s.applyEventDamage(event)
	s.AddLogKind(LogPlain, formatEventLogLine(event))
}

func (s *State) applyEventDamage(event EventDef) {
	switch {
	case event.DamageBuilding != "":
		if s.BuildingLevel(event.DamageBuilding) > 0 {
			s.damageBuilding(event.DamageBuilding)
			if def, ok := s.FindBuilding(event.DamageBuilding); ok {
				s.AddLogKind(LogDanger, fmt.Sprintf("%s damaged; daily output halved until repaired.", def.Name))
			}
		}
	case event.DamageRandomBuilt:
		s.damageRandomBuiltFacility()
	}
}

func (s *State) findEventByID(id string) (EventDef, bool) {
	for _, event := range s.Content.Events {
		if event.ID == id {
			return event, true
		}
	}
	return EventDef{}, false
}

// applyEventByID applies a known event by ID (no-op when unknown).
func (s *State) applyEventByID(id string) {
	if event, ok := s.findEventByID(id); ok {
		s.applyEvent(event)
	}
}

// eligibleEvents returns events that can fire on the given day.
func eligibleEvents(events []EventDef, day int) []EventDef {
	out := make([]EventDef, 0)
	for _, event := range events {
		if event.MinDay > day {
			continue
		}
		if event.MaxDay > 0 && day > event.MaxDay {
			continue
		}
		out = append(out, event)
	}
	return out
}

func eventWeight(e EventDef) int {
	if e.Weight <= 0 {
		return 1
	}
	return e.Weight
}

func (s *State) pickRandomEligibleEvent(candidates []EventDef) (EventDef, bool) {
	if len(candidates) == 0 {
		return EventDef{}, false
	}
	total := 0
	for _, e := range candidates {
		total += eventWeight(e)
	}
	if total <= 0 {
		return EventDef{}, false
	}
	s.ensureRNG()
	r := s.rng.Intn(total)
	for _, e := range candidates {
		w := eventWeight(e)
		if r < w {
			return e, true
		}
		r -= w
	}
	return candidates[len(candidates)-1], true
}
