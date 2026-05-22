package game

// Floors for the survival_30 reference script end state (#31).
const (
	SurvivalMinEndPower = 15
	SurvivalMinEndFood  = 10
)

func snapshotVitalsDepleted(s Snapshot) bool {
	return s.Power <= 0 || s.Food <= 0 || s.Morale <= 0 || s.Population <= 0
}

func snapshotMeetsSurvivalEndMargins(s Snapshot) bool {
	return s.Power >= SurvivalMinEndPower &&
		s.Food >= SurvivalMinEndFood &&
		s.Morale > 0 &&
		s.Population > 0
}
