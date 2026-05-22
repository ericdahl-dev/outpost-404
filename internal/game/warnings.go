package game

// Warning severity increases with resource pressure (higher = worse).
type WarningSeverity int

const (
	SeverityCaution WarningSeverity = iota + 1
	SeverityUrgent
	SeverityCritical
)

// Warning IDs for edge-trigger tracking.
const (
	WarningFoodLow       = "food_low"
	WarningPowerLow      = "power_low"
	WarningMoraleLow     = "morale_low"
	WarningPopulationLow = "population_low"
	WarningBeaconBehind  = "beacon_behind"
)

// Thresholds use current resources after the day resolves (not projected).
const (
	WarningFoodCriticalAt       = 15
	WarningFoodUrgentAt         = 25
	WarningPowerCriticalAt      = 15
	WarningPowerUrgentAt        = 25
	WarningMoraleCriticalAt     = 15
	WarningMoraleUrgentAt       = 25
	WarningPopulationCriticalAt = 2
	WarningBeaconLateDay        = 30
	WarningBeaconPartsBehind    = 2
)

// Warning is a derived lose-state alert for UI and log escalation.
type Warning struct {
	ID       string
	Severity WarningSeverity
	Message  string
}

// ActiveWarnings returns current alerts from colony vitals (no persistence).
func ActiveWarnings(s State) []Warning {
	var out []Warning
	if s.GameOver {
		return out
	}

	if s.Food <= WarningFoodCriticalAt {
		out = append(out, Warning{
			ID:       WarningFoodLow,
			Severity: SeverityCritical,
			Message:  "Food critically low — starvation collapse is imminent.",
		})
	} else if s.Food <= WarningFoodUrgentAt {
		out = append(out, Warning{
			ID:       WarningFoodLow,
			Severity: SeverityUrgent,
			Message:  "Food reserves unstable — next shortages will hit morale and population.",
		})
	}

	if s.Power <= WarningPowerCriticalAt {
		out = append(out, Warning{
			ID:       WarningPowerLow,
			Severity: SeverityCritical,
			Message:  "Power reserves critical — life support may fail this cycle.",
		})
	} else if s.Power <= WarningPowerUrgentAt {
		out = append(out, Warning{
			ID:       WarningPowerLow,
			Severity: SeverityUrgent,
			Message:  "Power reserves unstable — hydro and repairs will struggle.",
		})
	}

	if s.Morale <= WarningMoraleCriticalAt {
		out = append(out, Warning{
			ID:       WarningMoraleLow,
			Severity: SeverityCritical,
			Message:  "Morale near collapse — colonists may refuse work or leave.",
		})
	} else if s.Morale <= WarningMoraleUrgentAt {
		out = append(out, Warning{
			ID:       WarningMoraleLow,
			Severity: SeverityUrgent,
			Message:  "Morale fading — comfort losses will compound quickly.",
		})
	}

	if s.Population <= WarningPopulationCriticalAt {
		out = append(out, Warning{
			ID:       WarningPopulationLow,
			Severity: SeverityCritical,
			Message:  "Population at risk — the outpost cannot sustain further losses.",
		})
	}

	if s.Day >= WarningBeaconLateDay && s.BeaconParts < WarningBeaconPartsBehind && s.BeaconParts < s.MaxBeaconParts {
		out = append(out, Warning{
			ID:       WarningBeaconBehind,
			Severity: SeverityUrgent,
			Message:  "Beacon progress behind schedule — rescue window is closing.",
		})
	}

	return out
}

// syncWarnings edge-triggers colony log lines when a warning appears or worsens.
func (s *State) syncWarnings() {
	if s.WarningLevels == nil {
		s.WarningLevels = make(map[string]WarningSeverity)
	}
	for _, w := range ActiveWarnings(*s) {
		prev, seen := s.WarningLevels[w.ID]
		if !seen || w.Severity > prev {
			s.AddLog("! " + w.Message)
		}
		s.WarningLevels[w.ID] = w.Severity
	}
	for id := range s.WarningLevels {
		if !warningActive(id, s) {
			delete(s.WarningLevels, id)
		}
	}
}

func warningActive(id string, s *State) bool {
	for _, w := range ActiveWarnings(*s) {
		if w.ID == id {
			return true
		}
	}
	return false
}
