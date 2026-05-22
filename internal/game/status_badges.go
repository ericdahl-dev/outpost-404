package game

import "strings"

// Status badge labels for the main-screen status strip (readable without color).
const (
	BadgeFoodCritical = "FOOD CRITICAL"
	BadgeFoodLow      = "FOOD LOW"
	BadgePowerLow     = "POWER LOW"
	BadgeMoraleLow    = "MORALE UNSTABLE"
	BadgeBeaconReady  = "BEACON READY"
	BadgeTraderSignal = "TRADER SIGNAL"
)

// StatusBadge is a short operator alert for the resource panel strip.
type StatusBadge struct {
	Label    string
	Severity WarningSeverity
}

// StatusBadges returns current main-screen alerts derived from colony state.
func StatusBadges(s State) []StatusBadge {
	if s.GameOver {
		return nil
	}
	var out []StatusBadge

	for _, w := range ActiveWarnings(s) {
		switch w.ID {
		case WarningFoodLow:
			if w.Severity >= SeverityCritical {
				out = append(out, StatusBadge{BadgeFoodCritical, SeverityCritical})
			} else {
				out = append(out, StatusBadge{BadgeFoodLow, w.Severity})
			}
		case WarningPowerLow:
			out = append(out, StatusBadge{BadgePowerLow, w.Severity})
		case WarningMoraleLow:
			out = append(out, StatusBadge{BadgeMoraleLow, w.Severity})
		}
	}

	if name := firstDamagedFacilityName(s); name != "" {
		out = append(out, StatusBadge{strings.ToUpper(name) + " DAMAGED", SeverityUrgent})
	}

	if beaconWorkReady(s) {
		out = append(out, StatusBadge{BadgeBeaconReady, SeverityCaution})
	}

	if tradeAvailable(s) {
		out = append(out, StatusBadge{BadgeTraderSignal, SeverityCaution})
	}

	return out
}

func firstDamagedFacilityName(s State) string {
	for _, def := range s.Content.Buildings {
		if b, ok := s.Buildings[def.ID]; ok && b.Level > 0 && b.Damaged {
			return def.Name
		}
	}
	return ""
}

func beaconWorkReady(s State) bool {
	return s.BeaconParts < s.MaxBeaconParts && s.Power >= 18 && s.Credits >= 50
}

func tradeAvailable(s State) bool {
	return s.Food > MinFoodToTrade
}
