package game

import "fmt"

func NewState(content Content) State {
	return NewStateWithSeed(content, RandomSeed())
}

func newBareState(content Content) State {
	s := State{
		Day:            1,
		Power:          65,
		Food:           60,
		Morale:         70,
		Credits:        180,
		Population:     8,
		PopulationCap:  10,
		BeaconParts:    0,
		MaxBeaconParts: 5,
		Buildings:      map[string]Building{},
		Content:        content,
		Log:            []string{},
	}
	s.AddLog("Welcome to Outpost 404. Keep the systems online and finish the Signal Beacon.")
	s.AddLog("Survive 30 days or complete 5 beacon parts to win.")
	return s
}

func (s *State) AddLog(message string) {
	s.Log = append([]string{message}, s.Log...)
	if len(s.Log) > 8 {
		s.Log = s.Log[:8]
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
	} else if s.Power <= 0 || s.Food <= 0 || s.Morale <= 0 || s.Population <= 0 {
		s.GameOver = true
		s.Won = false
		s.Message = fmt.Sprintf("Outpost collapse on day %d. Final stats: power %d, food %d, morale %d, population %d.", s.Day, s.Power, s.Food, s.Morale, s.Population)
	} else if s.Day > 30 {
		s.GameOver = true
		s.Won = true
		s.Message = "You survived 30 days. Outpost 404 is stable enough to become permanent."
	}

	if s.GameOver {
		s.recordAction("game_end", map[string]any{
			"won":     s.Won,
			"message": s.Message,
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
