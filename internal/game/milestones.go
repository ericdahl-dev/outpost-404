package game

import "fmt"

// Milestone IDs for edge-triggered colony log beats (see issue #99).
const (
	MilestoneIDDay15            = "day_15"
	MilestoneIDDay30            = "day_30"
	MilestoneIDSurvivalImminent = "survival_imminent"
	MilestoneIDBeaconOneLeft    = "beacon_one_left"
)

func beaconEmphasisPart(max int) int {
	if max >= 5 {
		return 3
	}
	if max <= 1 {
		return 0
	}
	return (max + 1) / 2
}

func milestoneBeaconEmphasisID(parts, max int) string {
	return fmt.Sprintf("beacon_emphasis_%d_%d", parts, max)
}

// syncMilestones edge-triggers milestone log lines for survival pacing and beacon progress.
func (s *State) syncMilestones() {
	if s.GameOver {
		return
	}
	if s.MilestonesSeen == nil {
		s.MilestonesSeen = make(map[string]bool)
	}
	target := s.survivalWinAfterDay()

	for _, day := range []int{15, 30} {
		if day > target {
			continue
		}
		var id string
		switch day {
		case 15:
			id = MilestoneIDDay15
		case 30:
			id = MilestoneIDDay30
		default:
			id = fmt.Sprintf("day_%d", day)
		}
		if s.Day == day {
			s.emitMilestoneOnce(id, dayMilestoneMessage(day))
		}
	}

	if s.Day == target {
		s.emitMilestoneOnce(MilestoneIDSurvivalImminent,
			fmt.Sprintf("Survival window closes tomorrow — hold vitals through day %d.", target+1))
	}

	if s.BeaconParts > 0 && s.BeaconParts == s.MaxBeaconParts-1 {
		s.emitMilestoneOnce(MilestoneIDBeaconOneLeft,
			fmt.Sprintf("One beacon part remains — finish %d/%d before collapse.", s.BeaconParts, s.MaxBeaconParts))
	}

	emph := beaconEmphasisPart(s.MaxBeaconParts)
	if emph > 0 && s.BeaconParts == emph {
		id := milestoneBeaconEmphasisID(s.BeaconParts, s.MaxBeaconParts)
		s.emitMilestoneOnce(id, beaconEmphasisMessage(s.BeaconParts, s.MaxBeaconParts))
	}
}

func (s *State) emitMilestoneOnce(id, message string) {
	if s.MilestonesSeen[id] {
		return
	}
	s.MilestonesSeen[id] = true
	s.AddLogKind(LogMilestone, message)
}

func dayMilestoneMessage(day int) string {
	switch day {
	case 15:
		return "Day 15 — two weeks holding; the colony is still standing."
	case 30:
		return "Day 30 — a month out; beacon progress or endurance will decide the run."
	default:
		return fmt.Sprintf("Day %d — the colony endures.", day)
	}
}

func beaconEmphasisMessage(parts, max int) string {
	switch {
	case max >= 5 && parts == 3:
		return "Beacon three-fifths charged — rescue signal strengthening."
	case max == 3 && parts == 2:
		return "Beacon two-thirds charged — one part left for rescue."
	default:
		return fmt.Sprintf("Beacon %d/%d — halfway to full signal.", parts, max)
	}
}
