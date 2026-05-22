package game

import "testing"

func TestApplyDailyUpkeep_SubtractsResourcesAndAddsCredits(t *testing.T) {
	s := newTestState()
	s.Power = 40
	s.Food = 40
	s.Credits = 100

	s.applyDailyUpkeep()

	wantPower := 40 - DailyPowerUpkeep(s.Population)
	wantFood := 40 - DailyFoodUpkeep(s.Population)
	if s.Power != wantPower || s.Food != wantFood {
		t.Fatalf("power=%d food=%d, want power=%d food=%d", s.Power, s.Food, wantPower, wantFood)
	}
	if s.Credits != 100+DailyCreditsIncome {
		t.Fatalf("Credits = %d, want %d", s.Credits, 100+DailyCreditsIncome)
	}
}

func TestApplyMoraleDrift_ComfortBonus(t *testing.T) {
	s := newTestState()
	s.Morale = 50
	s.Power = ComfortPowerMin + 1
	s.Food = ComfortFoodMin + 1

	s.applyMoraleDrift()

	if s.Morale != 50+ComfortMoraleGain {
		t.Fatalf("Morale = %d, want %d", s.Morale, 50+ComfortMoraleGain)
	}
}

func TestApplyMoraleDrift_StressPenalty(t *testing.T) {
	s := newTestState()
	s.Morale = 50
	s.Power = ComfortPowerMin
	s.Food = ComfortFoodMin + 1

	s.applyMoraleDrift()

	if s.Morale != 50-StressMoraleLoss {
		t.Fatalf("Morale = %d, want %d", s.Morale, 50-StressMoraleLoss)
	}
}

func TestTryColonistArrival_OnIntervalDay(t *testing.T) {
	s := newTestState()
	s.Day = ColonistArrivalDayModulo
	s.Population = 8
	s.Food = ColonistFoodMin + 1
	s.Morale = ColonistMoraleMin + 1

	s.tryColonistArrival()

	if s.Population != 9 {
		t.Fatalf("Population = %d, want 9", s.Population)
	}
	wantLog := string(LogGain) + " " + colonistArrivalLog
	if len(s.Log) == 0 || s.Log[0] != wantLog {
		t.Fatalf("log = %q, want %q", s.Log[0], wantLog)
	}
}

func TestTryColonistArrival_SkipsWhenNotEligible(t *testing.T) {
	s := newTestState()
	s.Day = 4
	before := s.Population

	s.tryColonistArrival()

	if s.Population != before {
		t.Fatalf("Population = %d, want unchanged %d", s.Population, before)
	}
}
