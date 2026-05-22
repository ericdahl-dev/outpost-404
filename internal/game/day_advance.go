package game

const colonistArrivalLog = "A new colonist joined after hearing your beacon tests."

func (s *State) advanceDay() {
	s.Day++
	s.applyBuildingProduction()
	s.applyDailyUpkeep()
	s.applyMoraleDrift()
	s.tryColonistArrival()
}

func (s *State) applyDailyUpkeep() {
	s.Power -= DailyPowerUpkeep(s.Population)
	s.Food -= DailyFoodUpkeep(s.Population)
	s.Credits += s.dailyCreditsIncome()
}

func (s *State) applyMoraleDrift() {
	if ResourcesComfortable(s.Power, s.Food) {
		s.Morale += ComfortMoraleGain
	} else {
		s.Morale -= StressMoraleLoss
	}
}

func (s *State) tryColonistArrival() {
	if !CanGrowColonist(s.Day, s.Population, s.PopulationCap, s.Food, s.Morale) {
		return
	}
	s.Population++
	s.AddLogKind(LogGain, colonistArrivalLog)
}
