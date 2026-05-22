package game

// Daily economy balance (game layer only — not UI).
const (
	DailyPowerUpkeepBase  = 6
	DailyFoodUpkeepBase   = 4
	DailyUpkeepPopDivisor = 2
	DailyCreditsIncome    = 18

	ComfortPowerMin   = 50
	ComfortFoodMin    = 40
	ComfortMoraleGain = 2
	StressMoraleLoss  = 3

	ColonistArrivalDayModulo = 5
	ColonistFoodMin          = 35
	ColonistMoraleMin        = 40

	RandomEventRollSides     = 100
	RandomEventRollSkipAbove = 45 // roll > this skips the event (see RandomEventRollOccurs)
)

func DailyPowerUpkeep(population int) int {
	return DailyPowerUpkeepBase + population/DailyUpkeepPopDivisor
}

func DailyFoodUpkeep(population int) int {
	return DailyFoodUpkeepBase + population/DailyUpkeepPopDivisor
}

func ResourcesComfortable(power, food int) bool {
	return power > ComfortPowerMin && food > ComfortFoodMin
}

func CanGrowColonist(day, population, populationCap, food, morale int) bool {
	return day%ColonistArrivalDayModulo == 0 &&
		population < populationCap &&
		food > ColonistFoodMin &&
		morale > ColonistMoraleMin
}

func RandomEventRollOccurs(roll int) bool {
	return roll <= RandomEventRollSkipAbove
}
