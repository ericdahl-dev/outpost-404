package game

// DailyDeltas is the net change to vitals if the colony advances one day now (no random event).
type DailyDeltas struct {
	Power  int
	Food   int
	Morale int
}

// ProjectedDailyDeltas estimates per-day power, food, and morale change from upkeep,
// building dailyEffects (including damage halving), and morale drift.
func (s State) ProjectedDailyDeltas() DailyDeltas {
	var prodPower, prodFood, prodMorale int
	for _, def := range s.Content.Buildings {
		b, ok := s.Buildings[def.ID]
		if !ok || b.Level <= 0 {
			continue
		}
		effects := dailyEffectsScaled(def, b)
		prodPower += effects["power"]
		prodFood += effects["food"]
		prodMorale += effects["morale"]
	}
	pop := s.Population
	drift := -StressMoraleLoss
	if ResourcesComfortable(s.Power, s.Food) {
		drift = ComfortMoraleGain
	}
	return DailyDeltas{
		Power:  prodPower - DailyPowerUpkeep(pop),
		Food:   prodFood - DailyFoodUpkeep(pop),
		Morale: prodMorale + drift,
	}
}
