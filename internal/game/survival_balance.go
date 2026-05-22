package game

// Floors for the survival_45 reference script end state (#31).
const (
	SurvivalMinEndPower = 4 // longer 45-day path drains more power on reference seeds
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
