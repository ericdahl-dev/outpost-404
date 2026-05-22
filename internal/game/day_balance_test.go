package game

import "testing"

func TestDailyPowerUpkeep_DefaultPopulation(t *testing.T) {
	got := DailyPowerUpkeep(8)
	want := DailyPowerUpkeepBase + 8/DailyUpkeepPopDivisor
	if got != want || got != 10 {
		t.Fatalf("DailyPowerUpkeep(8) = %d, want %d", got, want)
	}
}

func TestDailyFoodUpkeep_DefaultPopulation(t *testing.T) {
	got := DailyFoodUpkeep(8)
	want := DailyFoodUpkeepBase + 8/DailyUpkeepPopDivisor
	if got != want || got != 8 {
		t.Fatalf("DailyFoodUpkeep(8) = %d, want %d", got, want)
	}
}

func TestResourcesComfortable(t *testing.T) {
	if !ResourcesComfortable(51, 41) {
		t.Fatal("expected comfortable when above comfort floors")
	}
	if ResourcesComfortable(50, 41) || ResourcesComfortable(51, 40) {
		t.Fatal("expected stress when power or food at comfort floor")
	}
}

func TestCanGrowColonist_RequiresIntervalAndThresholds(t *testing.T) {
	if !CanGrowColonist(5, 8, 10, 36, 41) {
		t.Fatal("expected colonist growth on day 5 with stats above thresholds")
	}
	if CanGrowColonist(4, 8, 10, 36, 41) {
		t.Fatal("expected no growth off interval day")
	}
	if CanGrowColonist(5, 8, 10, 35, 41) {
		t.Fatal("expected no growth when food at floor")
	}
}

func TestRandomEventRollOccurs(t *testing.T) {
	if !RandomEventRollOccurs(0) || !RandomEventRollOccurs(RandomEventRollSkipAbove) {
		t.Fatal("expected event on low rolls")
	}
	if RandomEventRollOccurs(RandomEventRollSkipAbove + 1) {
		t.Fatal("expected skip above threshold")
	}
}
