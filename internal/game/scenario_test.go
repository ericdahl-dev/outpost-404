package game

import "testing"

func TestNewRun_BeaconRush_WinsAtThreeBeaconParts(t *testing.T) {
	profiles := testRunProfiles()
	s := NewRun(testContent(), profiles, 1, "beacon_rush", "normal")
	if s.MaxBeaconParts != 3 {
		t.Fatalf("MaxBeaconParts = %d, want 3", s.MaxBeaconParts)
	}
	s.BeaconParts = 2
	s.Power = 50
	s.Credits = 100
	s.WorkOnBeacon()
	if !s.GameOver || !s.Won {
		t.Fatal("expected beacon win at 3/3 parts")
	}
}

func TestNewRun_FirstLanding_ExtraStartingCredits(t *testing.T) {
	std := NewRun(testContent(), testRunProfiles(), 1, "standard", "normal")
	fl := NewRun(testContent(), testRunProfiles(), 1, "first_landing", "normal")
	if fl.Credits <= std.Credits {
		t.Fatalf("first_landing credits %d want > standard %d", fl.Credits, std.Credits)
	}
}

func TestEligibleEvents_RequiresBuilding_SilentColonyOnly(t *testing.T) {
	events := []EventDef{
		{ID: "open", MinDay: 1},
		{ID: "radio_only", MinDay: 1, RequiresBuilding: "radio_tower"},
	}
	std := NewRun(testContent(), testRunProfiles(), 1, "standard", "normal")
	std.Day = 5
	if len(eligibleEventsForState(std, events)) != 2 {
		t.Fatal("standard run ignores requiresBuilding gate")
	}
	silent := NewRun(testContent(), testRunProfiles(), 1, "silent_colony", "normal")
	silent.Day = 5
	got := eligibleEventsForState(silent, events)
	if len(got) != 1 || got[0].ID != "open" {
		t.Fatalf("silent without radio: %+v, want only open", got)
	}
	silent.Buildings["radio_tower"] = Building{DefID: "radio_tower", Level: 1}
	got = eligibleEventsForState(silent, events)
	if len(got) != 2 {
		t.Fatalf("silent with radio: %+v, want 2 events", got)
	}
}

func TestNewRun_Hard_LowersEventGate(t *testing.T) {
	easy := NewRun(testContent(), testRunProfiles(), 1, "standard", "easy")
	hard := NewRun(testContent(), testRunProfiles(), 1, "standard", "hard")
	if easy.EventGateSkipAbove >= hard.EventGateSkipAbove {
		t.Fatalf("easy gate %d should be < hard gate %d", easy.EventGateSkipAbove, hard.EventGateSkipAbove)
	}
}

func testRunProfiles() RunProfiles {
	p, err := LoadEmbeddedRunProfiles()
	if err != nil {
		panic(err)
	}
	return p
}
