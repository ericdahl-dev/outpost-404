package game

import "testing"

func TestBuild_UpgradesBuildingAndDeductsScaledCost(t *testing.T) {
	s := newTestState()
	startCredits := s.Credits

	s.Build("solar_array")

	if got := s.BuildingLevel("solar_array"); got != 1 {
		t.Fatalf("BuildingLevel(solar_array) = %d, want 1", got)
	}
	if want := startCredits - 70; s.Credits != want {
		t.Fatalf("Credits = %d, want %d", s.Credits, want)
	}
	if s.Power != 83 {
		t.Fatalf("Power = %d, want 83 (65 + 18)", s.Power)
	}
}

func TestBuild_SecondLevelCostsDoubleBaseCost(t *testing.T) {
	s := newTestState()
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 1}
	s.Credits = 200

	s.Build("solar_array")

	if got := s.BuildingLevel("solar_array"); got != 2 {
		t.Fatalf("BuildingLevel(solar_array) = %d, want 2", got)
	}
	if s.Credits != 60 {
		t.Fatalf("Credits = %d, want 60 (200 - 140)", s.Credits)
	}
}

func TestBuild_RejectsWhenNotEnoughCredits(t *testing.T) {
	s := newTestState()
	s.Credits = 10

	s.Build("solar_array")

	if got := s.BuildingLevel("solar_array"); got != 0 {
		t.Fatalf("BuildingLevel(solar_array) = %d, want 0", got)
	}
	if s.Credits != 10 {
		t.Fatalf("Credits = %d, want unchanged 10", s.Credits)
	}
}

func TestBuild_RejectsAtMaxLevel(t *testing.T) {
	s := newTestState()
	s.Buildings["solar_array"] = Building{DefID: "solar_array", Level: 3}
	s.Credits = 500

	s.Build("solar_array")

	if got := s.BuildingLevel("solar_array"); got != 3 {
		t.Fatalf("BuildingLevel(solar_array) = %d, want 3", got)
	}
	if s.Credits != 500 {
		t.Fatalf("Credits = %d, want unchanged 500", s.Credits)
	}
}

func TestBuild_NoOpWhenGameOver(t *testing.T) {
	s := newTestState()
	s.GameOver = true
	s.Credits = 500

	s.Build("solar_array")

	if got := s.BuildingLevel("solar_array"); got != 0 {
		t.Fatalf("BuildingLevel(solar_array) = %d, want 0 after game over", got)
	}
}

func TestWorkOnBeacon_CompletesPartAndSpendsResources(t *testing.T) {
	s := newTestState()
	s.Power = 50
	s.Credits = 100

	s.WorkOnBeacon()

	if s.BeaconParts != 1 {
		t.Fatalf("BeaconParts = %d, want 1", s.BeaconParts)
	}
	if s.Power != 38 {
		t.Fatalf("Power = %d, want 38", s.Power)
	}
	if s.Credits != 50 {
		t.Fatalf("Credits = %d, want 50", s.Credits)
	}
	if s.Morale != 75 {
		t.Fatalf("Morale = %d, want 75", s.Morale)
	}
}

func TestWorkOnBeacon_WinsWhenBeaconComplete(t *testing.T) {
	s := newTestState()
	s.BeaconParts = 4
	s.Power = 50
	s.Credits = 100

	s.WorkOnBeacon()

	if !s.GameOver || !s.Won {
		t.Fatal("expected win after fifth beacon part")
	}
	if s.BeaconParts != 5 {
		t.Fatalf("BeaconParts = %d, want 5", s.BeaconParts)
	}
}

func TestWorkOnBeacon_RejectsWithoutPowerAndCredits(t *testing.T) {
	s := newTestState()
	s.Power = 10
	s.Credits = 100

	s.WorkOnBeacon()

	if s.BeaconParts != 0 {
		t.Fatalf("BeaconParts = %d, want 0", s.BeaconParts)
	}
}

func TestRepair_RejectsWhenGameOver(t *testing.T) {
	s := newTestState()
	s.GameOver = true
	s.Credits = 100

	s.Repair()

	if s.Credits != 100 {
		t.Fatalf("Credits = %d, want unchanged 100", s.Credits)
	}
}

func TestRepair_RejectsInsufficientCredits(t *testing.T) {
	s := newTestState()
	s.Credits = 10

	s.Repair()

	if s.Credits != 10 {
		t.Fatalf("Credits = %d, want unchanged 10", s.Credits)
	}
	if len(s.Log) == 0 || s.Log[0] == "" {
		t.Fatal("expected rejection log line")
	}
}

func TestRepair_TriggersGameOverWhenFoodDepleted(t *testing.T) {
	s := newTestState()
	s.Food = 0
	s.Power = 100
	s.Morale = 100
	s.Credits = 100

	s.Repair()

	if !s.GameOver {
		t.Fatal("expected GameOver when food is 0")
	}
	if s.Won {
		t.Fatal("expected loss, not win")
	}
}

func TestTrade_TriggersGameOverWhenMoraleDepleted(t *testing.T) {
	s := newTestState()
	s.Morale = 3
	s.Food = 100
	s.Power = 100
	s.Credits = 0

	s.Trade()

	if !s.GameOver {
		t.Fatal("expected GameOver when morale reaches 0")
	}
}
