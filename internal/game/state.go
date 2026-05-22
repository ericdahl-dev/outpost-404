package game

import "fmt"

func NewState(content Content) State {
	return NewStateWithSeed(content, RandomSeed())
}

func (s *State) AddLog(message string) {
	s.Log = append([]string{message}, s.Log...)
	if len(s.Log) > maxColonyLogLines {
		s.Log = s.Log[:maxColonyLogLines]
	}
}

func (s State) BuildingLevel(id string) int {
	if b, ok := s.Buildings[id]; ok {
		return b.Level
	}
	return 0
}

func (s State) FindBuilding(id string) (BuildingDef, bool) {
	return s.Content.FindBuilding(id)
}

func (s *State) Clamp() {
	s.Power = clamp(s.Power, 0, 100)
	s.Food = clamp(s.Food, 0, 100)
	s.Morale = clamp(s.Morale, 0, 100)
	if s.Credits < 0 {
		s.Credits = 0
	}
	if s.Population < 0 {
		s.Population = 0
	}
	if s.Population > s.PopulationCap {
		s.Population = s.PopulationCap
	}
}

func (s *State) CheckEnd() {
	if s.GameOver {
		return
	}
	before := s.snapshot()

	if s.BeaconParts >= s.MaxBeaconParts {
		s.GameOver = true
		s.Won = true
		s.Message = "Signal Beacon complete. Rescue is inbound. Outpost 404 survives."
		s.AddLogKind(LogMilestone, "The Signal Beacon reached full charge. Outpost 404 is no longer alone.")
	} else if s.Power <= 0 || s.Food <= 0 || s.Morale <= 0 || s.Population <= 0 {
		s.GameOver = true
		s.Won = false
		cause := collapseCause(*s)
		s.Message = fmt.Sprintf("Outpost collapse on day %d. Final stats: power %d, food %d, morale %d, population %d.", s.Day, s.Power, s.Food, s.Morale, s.Population)
		s.AddLogKind(LogMilestone, fmt.Sprintf("Outpost collapsed — %s.", cause))
	} else if s.Day > s.survivalWinAfterDay() {
		target := s.survivalWinAfterDay()
		s.GameOver = true
		s.Won = true
		s.Message = fmt.Sprintf("You survived %d days. Outpost 404 is stable enough to become permanent.", target)
		s.AddLogKind(LogMilestone, fmt.Sprintf("Day %d: the colony endured. Survival victory secured.", target))
	}

	if s.GameOver {
		summary := s.SessionSummary()
		s.recordAction("game_end", map[string]any{
			"won":              s.Won,
			"message":          s.Message,
			"session_summary":  summary,
		}, before, s.snapshot())
	}
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
