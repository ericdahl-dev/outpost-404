package game

import "testing"

func TestStatusBadges_FoodCritical(t *testing.T) {
	s := newTestState()
	s.Food = WarningFoodCriticalAt
	badges := StatusBadges(s)
	if !containsBadge(badges, "FOOD CRITICAL") {
		t.Fatalf("expected FOOD CRITICAL, got %v", badgeLabels(badges))
	}
}

func TestStatusBadges_PowerLow(t *testing.T) {
	s := newTestState()
	s.Power = WarningPowerUrgentAt
	badges := StatusBadges(s)
	if !containsBadge(badges, "POWER LOW") {
		t.Fatalf("expected POWER LOW, got %v", badgeLabels(badges))
	}
}

func TestStatusBadges_FacilityDamaged(t *testing.T) {
	s := newTestState()
	s.Buildings["workshop"] = Building{DefID: "workshop", Level: 1, Damaged: true}
	badges := StatusBadges(s)
	if !containsBadgePrefix(badges, "WORKSHOP DAMAGED") {
		t.Fatalf("expected workshop damaged badge, got %v", badgeLabels(badges))
	}
}

func TestStatusBadges_BeaconReady(t *testing.T) {
	s := newTestState()
	s.Power = 20
	s.Credits = 60
	badges := StatusBadges(s)
	if !containsBadge(badges, "BEACON READY") {
		t.Fatalf("expected BEACON READY, got %v", badgeLabels(badges))
	}
}

func TestStatusBadges_TraderSignal(t *testing.T) {
	s := newTestState()
	s.Food = MinFoodToTrade + 1
	badges := StatusBadges(s)
	if !containsBadge(badges, "TRADER SIGNAL") {
		t.Fatalf("expected TRADER SIGNAL, got %v", badgeLabels(badges))
	}
}

func TestStatusBadges_GameOverEmpty(t *testing.T) {
	s := newTestState()
	s.GameOver = true
	if len(StatusBadges(s)) != 0 {
		t.Fatalf("expected no badges when game over, got %v", badgeLabels(StatusBadges(s)))
	}
}

func containsBadge(badges []StatusBadge, label string) bool {
	for _, b := range badges {
		if b.Label == label {
			return true
		}
	}
	return false
}

func containsBadgePrefix(badges []StatusBadge, prefix string) bool {
	for _, b := range badges {
		if len(b.Label) >= len(prefix) && b.Label[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func badgeLabels(badges []StatusBadge) []string {
	out := make([]string, len(badges))
	for i, b := range badges {
		out[i] = b.Label
	}
	return out
}
