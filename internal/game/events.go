package game

// applyEvent applies event effects and appends the standard log line.
func (s *State) applyEvent(event EventDef) {
	s.applyEffects(event.Effects, 1)
	s.AddLog(event.Title + ": " + event.Description)
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
