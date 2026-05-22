package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"testing"
)

func TestDailyPowerUpkeep_DefaultPopulation(t *testing.T) {
	got := game.DailyPowerUpkeep(8)
	want := game.DailyPowerUpkeepBase + 8/game.DailyUpkeepPopDivisor
	if got != want || got != 10 {
		t.Fatalf("game.DailyPowerUpkeep(8) = %d, want %d", got, want)
	}
}

func TestDailyFoodUpkeep_DefaultPopulation(t *testing.T) {
	got := game.DailyFoodUpkeep(8)
	want := game.DailyFoodUpkeepBase + 8/game.DailyUpkeepPopDivisor
	if got != want || got != 8 {
		t.Fatalf("game.DailyFoodUpkeep(8) = %d, want %d", got, want)
	}
}

func TestResourcesComfortable(t *testing.T) {
	if !game.ResourcesComfortable(51, 41) {
		t.Fatal("expected comfortable when above comfort floors")
	}
	if game.ResourcesComfortable(50, 41) || game.ResourcesComfortable(51, 40) {
		t.Fatal("expected stress when power or food at comfort floor")
	}
}

func TestCanGrowColonist_RequiresIntervalAndThresholds(t *testing.T) {
	if !game.CanGrowColonist(5, 8, 10, 36, 41) {
		t.Fatal("expected colonist growth on day 5 with stats above thresholds")
	}
	if game.CanGrowColonist(4, 8, 10, 36, 41) {
		t.Fatal("expected no growth off interval day")
	}
	if game.CanGrowColonist(5, 8, 10, 35, 41) {
		t.Fatal("expected no growth when food at floor")
	}
}

func TestRandomEventRollOccurs(t *testing.T) {
	if !game.RandomEventRollOccurs(0) || !game.RandomEventRollOccurs(game.RandomEventRollSkipAbove) {
		t.Fatal("expected event on low rolls")
	}
	if game.RandomEventRollOccurs(game.RandomEventRollSkipAbove + 1) {
		t.Fatal("expected skip above threshold")
	}
}
