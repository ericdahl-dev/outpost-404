package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"testing"
)

func TestApplyDailyUpkeep_SubtractsResourcesAndAddsCredits(t *testing.T) {
	s := newTestState()
	s.Power = 40
	s.Food = 40
	s.Credits = 100

	game.ApplyDailyUpkeepForTest(&s)

	wantPower := 40 - game.DailyPowerUpkeep(s.Population)
	wantFood := 40 - game.DailyFoodUpkeep(s.Population)
	if s.Power != wantPower || s.Food != wantFood {
		t.Fatalf("power=%d food=%d, want power=%d food=%d", s.Power, s.Food, wantPower, wantFood)
	}
	if s.Credits != 100+game.DailyCreditsIncome {
		t.Fatalf("Credits = %d, want %d", s.Credits, 100+game.DailyCreditsIncome)
	}
}

func TestApplyMoraleDrift_ComfortBonus(t *testing.T) {
	s := newTestState()
	s.Morale = 50
	s.Power = game.ComfortPowerMin + 1
	s.Food = game.ComfortFoodMin + 1

	game.ApplyMoraleDriftForTest(&s)

	if s.Morale != 50+game.ComfortMoraleGain {
		t.Fatalf("Morale = %d, want %d", s.Morale, 50+game.ComfortMoraleGain)
	}
}

func TestApplyMoraleDrift_StressPenalty(t *testing.T) {
	s := newTestState()
	s.Morale = 50
	s.Power = game.ComfortPowerMin
	s.Food = game.ComfortFoodMin + 1

	game.ApplyMoraleDriftForTest(&s)

	if s.Morale != 50-game.StressMoraleLoss {
		t.Fatalf("Morale = %d, want %d", s.Morale, 50-game.StressMoraleLoss)
	}
}

func TestTryColonistArrival_OnIntervalDay(t *testing.T) {
	s := newTestState()
	s.Day = game.ColonistArrivalDayModulo
	s.Population = 8
	s.Food = game.ColonistFoodMin + 1
	s.Morale = game.ColonistMoraleMin + 1

	game.TryColonistArrivalForTest(&s)

	if s.Population != 9 {
		t.Fatalf("Population = %d, want 9", s.Population)
	}
	if len(s.Log) == 0 || s.Log[0] != game.ColonistArrivalLog {
		t.Fatalf("log = %q, want colonist arrival line", s.Log[0])
	}
}

func TestTryColonistArrival_SkipsWhenNotEligible(t *testing.T) {
	s := newTestState()
	s.Day = 4
	before := s.Population

	game.TryColonistArrivalForTest(&s)

	if s.Population != before {
		t.Fatalf("Population = %d, want unchanged %d", s.Population, before)
	}
}
