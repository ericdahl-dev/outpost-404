package game

import "testing"

func TestNextDay_AdvancesDayAndConsumesResources(t *testing.T) {
	s := newTestState()
	startDay := s.Day
	startPower := s.Power
	startFood := s.Food
	startCredits := s.Credits

	s.NextDay()

	if s.Day != startDay+1 {
		t.Fatalf("Day = %d, want %d", s.Day, startDay+1)
	}
	wantPower := startPower - DailyPowerUpkeep(s.Population)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d", s.Power, wantPower)
	}
	wantFood := startFood - DailyFoodUpkeep(s.Population)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d", s.Food, wantFood)
	}
	if s.Credits != startCredits+DailyCreditsIncome {
		t.Fatalf("Credits = %d, want %d", s.Credits, startCredits+DailyCreditsIncome)
	}
}

func TestNextDay_NoOpWhenGameOver(t *testing.T) {
	s := newTestState()
	s.GameOver = true
	s.Day = 5

	s.NextDay()

	if s.Day != 5 {
		t.Fatalf("Day = %d, want unchanged 5", s.Day)
	}
}

func TestNextDay_WinsAfterSurvivalTarget(t *testing.T) {
	s := newTestState()
	s.Day = SurvivalWinAfterDay
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.NextDay()

	wantDay := SurvivalWinAfterDay + 1
	if !s.GameOver || !s.Won {
		t.Fatalf("expected win after surviving into day %d", wantDay)
	}
	if s.Day != wantDay {
		t.Fatalf("Day = %d, want %d", s.Day, wantDay)
	}
}

func TestApplyEffects_PopulationEffect(t *testing.T) {
	s := newTestState()
	s.Population = 5
	s.applyEffects(map[string]int{"population": 3}, 1)
	if s.Population != 8 {
		t.Fatalf("Population = %d, want 8", s.Population)
	}
}
