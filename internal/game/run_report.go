package game

import "fmt"

// RunReport is structured end-of-run data for the game-over screen.
type RunReport struct {
	Won            bool
	Cause          string
	Tagline        string
	Day            int
	Power          int
	Food           int
	Morale         int
	Credits        int
	Population     int
	PopulationCap  int
	BeaconParts    int
	MaxBeaconParts int
	LowestPower    int
	LowestFood     int
	LowestMorale   int
	ScenarioName   string
	DifficultyName string
	ScenarioID     string
	DifficultyID   string
	Seed           int64
	BuiltCount     int
	KeyMoments     []string
}

// EndCause returns a short headline for win or loss.
func (s State) EndCause() string {
	if s.Won {
		if s.BeaconParts >= s.MaxBeaconParts {
			return "Signal Beacon complete — rescue inbound"
		}
		target := s.survivalWinAfterDay()
		return fmt.Sprintf("Survived %d days — colony stabilized", target)
	}
	switch {
	case s.Power <= 0:
		return "Power reserves exhausted"
	case s.Food <= 0:
		return "Food reserves depleted"
	case s.Morale <= 0:
		return "Morale collapsed"
	case s.Population <= 0:
		return "Population lost"
	default:
		return "Critical systems failed"
	}
}

func (s *State) initVitalLows() {
	s.MinPowerSeen = s.Power
	s.MinFoodSeen = s.Food
	s.MinMoraleSeen = s.Morale
}

func (s *State) recordVitalLows() {
	if s.GameOver {
		return
	}
	if s.Power < s.MinPowerSeen {
		s.MinPowerSeen = s.Power
	}
	if s.Food < s.MinFoodSeen {
		s.MinFoodSeen = s.Food
	}
	if s.Morale < s.MinMoraleSeen {
		s.MinMoraleSeen = s.Morale
	}
}

// BuildRunReport assembles display fields for the end screen.
func BuildRunReport(s State, profiles RunProfiles) RunReport {
	scName := s.ScenarioID
	if sc, ok := profiles.FindScenario(s.ScenarioID); ok {
		scName = sc.Name
	}
	diffName := s.DifficultyID
	if d, ok := profiles.FindDifficulty(s.DifficultyID); ok {
		diffName = d.Name
	}
	built := 0
	for _, b := range s.Buildings {
		if b.Level > 0 {
			built++
		}
	}
	moments := append([]string(nil), s.KeyMoments...)
	return RunReport{
		Won:            s.Won,
		Cause:          s.EndCause(),
		Tagline:        s.Message,
		Day:            s.Day,
		Power:          s.Power,
		Food:           s.Food,
		Morale:         s.Morale,
		Credits:        s.Credits,
		Population:     s.Population,
		PopulationCap:  s.PopulationCap,
		BeaconParts:    s.BeaconParts,
		MaxBeaconParts: s.MaxBeaconParts,
		LowestPower:    s.MinPowerSeen,
		LowestFood:     s.MinFoodSeen,
		LowestMorale:   s.MinMoraleSeen,
		ScenarioName:   scName,
		DifficultyName: diffName,
		ScenarioID:     s.ScenarioID,
		DifficultyID:   s.DifficultyID,
		Seed:           s.Seed,
		BuiltCount:     built,
		KeyMoments:     moments,
	}
}
