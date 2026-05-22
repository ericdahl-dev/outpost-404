package game

import "testing"

func TestActiveWarnings_FoodCritical(t *testing.T) {
	s := newTestState()
	s.Food = WarningFoodCriticalAt
	got := ActiveWarnings(s)
	if !containsWarning(got, WarningFoodLow) {
		t.Fatalf("expected %s, got %v", WarningFoodLow, got)
	}
}

func TestActiveWarnings_NoWarningWhenResourcesHealthy(t *testing.T) {
	s := newTestState()
	s.Power = 60
	s.Food = 60
	s.Morale = 60
	if len(ActiveWarnings(s)) != 0 {
		t.Fatalf("expected no warnings, got %v", ActiveWarnings(s))
	}
}

func TestSyncWarnings_LogsOnFirstEscalationOnly(t *testing.T) {
	s := newTestState()
	s.Food = WarningFoodCriticalAt
	s.syncWarnings()
	var warned bool
	for _, line := range s.Log {
		if len(line) > 0 && line[0] == '!' {
			warned = true
			break
		}
	}
	if !warned {
		t.Fatalf("expected warning in log, got %v", s.Log)
	}
	before := len(s.Log)
	s.syncWarnings()
	if len(s.Log) != before {
		t.Fatal("expected no duplicate warning log while unchanged")
	}
}

func containsWarning(ws []Warning, id string) bool {
	for _, w := range ws {
		if w.ID == id {
			return true
		}
	}
	return false
}
