package game

import "testing"

func TestClamp_KeepsCoreResourcesInRange(t *testing.T) {
	s := newTestState()
	s.Power = 150
	s.Food = -5
	s.Morale = 200
	s.Credits = -10
	s.Population = 99
	s.PopulationCap = 10

	s.Clamp()

	if s.Power != 100 {
		t.Fatalf("Power = %d, want 100", s.Power)
	}
	if s.Food != 0 {
		t.Fatalf("Food = %d, want 0", s.Food)
	}
	if s.Morale != 100 {
		t.Fatalf("Morale = %d, want 100", s.Morale)
	}
	if s.Credits != 0 {
		t.Fatalf("Credits = %d, want 0", s.Credits)
	}
	if s.Population != 10 {
		t.Fatalf("Population = %d, want capped at 10", s.Population)
	}
}

func TestCheckEnd_LossWhenCriticalResourceDepleted(t *testing.T) {
	s := newTestState()
	s.Power = 0

	s.CheckEnd()

	if !s.GameOver || s.Won {
		t.Fatal("expected loss when power is 0")
	}
}

func TestCheckEnd_WinAfterSurvivingTargetDays(t *testing.T) {
	s := newTestState()
	s.Day = SurvivalWinAfterDay + 1
	s.Power = 20
	s.Food = 20
	s.Morale = 50
	s.Population = 8

	s.CheckEnd()

	if !s.GameOver || !s.Won {
		t.Fatalf("expected win after day %d with positive vitals", s.Day)
	}
}

func TestNewState_StartsPlayableWithBeaconGoal(t *testing.T) {
	s := newTestState()

	if s.Day != 1 {
		t.Fatalf("Day = %d, want 1", s.Day)
	}
	if s.MaxBeaconParts != 5 {
		t.Fatalf("MaxBeaconParts = %d, want 5", s.MaxBeaconParts)
	}
	if s.GameOver {
		t.Fatal("new game should not be over")
	}
}
