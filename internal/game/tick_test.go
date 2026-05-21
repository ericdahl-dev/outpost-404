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
	wantPower := startPower - (8 + s.Population/2)
	if s.Power != wantPower {
		t.Fatalf("Power = %d, want %d", s.Power, wantPower)
	}
	wantFood := startFood - (6 + s.Population/2)
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d", s.Food, wantFood)
	}
	if s.Credits != startCredits+18 {
		t.Fatalf("Credits = %d, want %d", s.Credits, startCredits+18)
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

func TestNextDay_WinsAfterDay30(t *testing.T) {
	s := newTestState()
	s.Day = 30
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.NextDay()

	if !s.GameOver || !s.Won {
		t.Fatal("expected win after surviving into day 31")
	}
	if s.Day != 31 {
		t.Fatalf("Day = %d, want 31", s.Day)
	}
}
