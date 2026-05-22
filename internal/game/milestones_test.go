package game

import (
	"strings"
	"testing"
)

func milestoneLines(s State) []string {
	var out []string
	for _, line := range s.Log {
		if strings.HasPrefix(line, "* ") {
			out = append(out, line)
		}
	}
	return out
}

func countMilestoneContaining(s State, sub string) int {
	n := 0
	for _, line := range milestoneLines(s) {
		if strings.Contains(line, sub) {
			n++
		}
	}
	return n
}

func TestSyncMilestones_Day15OnNextDay(t *testing.T) {
	s := newTestState()
	s.Day = 14
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.NextDay()

	if s.Day != 15 {
		t.Fatalf("Day = %d, want 15", s.Day)
	}
	if count := countMilestoneContaining(s, "Day 15"); count != 1 {
		t.Fatalf("Day 15 milestone count = %d, want 1; milestones %v", count, milestoneLines(s))
	}
}

func TestSyncMilestones_Day30OnNextDay(t *testing.T) {
	s := newTestState()
	s.Day = 29
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.NextDay()

	if s.Day != 30 {
		t.Fatalf("Day = %d, want 30", s.Day)
	}
	if count := countMilestoneContaining(s, "Day 30"); count != 1 {
		t.Fatalf("Day 30 milestone count = %d, want 1; milestones %v", count, milestoneLines(s))
	}
}

func TestSyncMilestones_Day15NoDuplicateOnSecondNextDay(t *testing.T) {
	s := newTestState()
	s.Day = 14
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.NextDay()
	before := countMilestoneContaining(s, "Day 15")
	s.Power = 50
	s.Food = 50
	s.Morale = 50
	s.NextDay()
	after := countMilestoneContaining(s, "Day 15")
	if before != 1 || after != 1 {
		t.Fatalf("Day 15 milestones before=%d after=%d, want 1 each", before, after)
	}
}

func TestSyncMilestones_BeaconPartAddsMilestoneNoDuplicateWithoutNewPart(t *testing.T) {
	s := newTestState()
	s.Power = 50
	s.Credits = 100

	s.WorkOnBeacon()
	if countMilestoneContaining(s, "Signal Beacon part completed") != 1 {
		t.Fatalf("expected part-completed milestone, got %v", milestoneLines(s))
	}

	before := len(s.Log)
	s.Power = 50
	s.Food = 50
	s.Morale = 50
	s.NextDay()
	for _, line := range s.Log {
		if strings.Contains(line, "Signal Beacon part completed") {
			if len(s.Log) != before {
				t.Fatal("unexpected duplicate beacon part milestone without new part")
			}
		}
	}
	if countMilestoneContaining(s, "Signal Beacon part completed") != 1 {
		t.Fatalf("beacon part milestone count changed without new part")
	}
}

func TestSyncMilestones_BeaconEmphasisAtThreeOfFive(t *testing.T) {
	s := newTestState()
	s.Power = 50
	s.Credits = 200
	s.BeaconParts = 2

	s.WorkOnBeacon()

	if s.BeaconParts != 3 {
		t.Fatalf("BeaconParts = %d, want 3", s.BeaconParts)
	}
	if countMilestoneContaining(s, "three-fifths") != 1 {
		t.Fatalf("expected three-fifths emphasis, got %v", milestoneLines(s))
	}
}

func TestSyncMilestones_SurvivalImminentAtTargetDay(t *testing.T) {
	s := newTestState()
	s.Day = SurvivalWinAfterDay
	s.Power = 50
	s.Food = 50
	s.Morale = 50

	s.syncMilestones()

	if countMilestoneContaining(s, "day 46") != 1 {
		t.Fatalf("expected survival imminent milestone, got %v", milestoneLines(s))
	}
}

func TestNewRun_BeaconRush_EmphasisAtTwoOfThree(t *testing.T) {
	content := testContent()
	profiles, err := LoadEmbeddedRunProfiles()
	if err != nil {
		t.Skip("embedded run profiles unavailable")
	}
	s := NewRun(content, profiles, 1, "beacon_rush", "normal")
	s.Power = 50
	s.Credits = 200
	s.BeaconParts = 1

	s.WorkOnBeacon()

	if s.BeaconParts != 2 || s.MaxBeaconParts != 3 {
		t.Fatalf("beacon %d/%d, want 2/3", s.BeaconParts, s.MaxBeaconParts)
	}
	if countMilestoneContaining(s, "two-thirds") != 1 {
		t.Fatalf("expected two-thirds emphasis for beacon rush, got %v", milestoneLines(s))
	}
}
