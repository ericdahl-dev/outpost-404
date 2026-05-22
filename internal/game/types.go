package game

import "math/rand"

type Resource string

const (
	ResourcePower   Resource = "power"
	ResourceFood    Resource = "food"
	ResourceMorale  Resource = "morale"
	ResourceCredits Resource = "credits"
)

type BuildingDef struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Cost         int            `json:"cost"`
	MaxLevel     int            `json:"maxLevel"`
	Effects      map[string]int `json:"effects"`
	DailyEffects map[string]int `json:"dailyEffects,omitempty"`
}

type EventDef struct {
	ID                string         `json:"id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Effects           map[string]int `json:"effects"`
	MinDay            int            `json:"minDay"`
	MaxDay            int            `json:"maxDay,omitempty"`
	Weight            int            `json:"weight,omitempty"`
	DamageBuilding     string         `json:"damageBuilding,omitempty"`
	DamageRandomBuilt  bool           `json:"damageRandomBuilt,omitempty"`
	RequiresBuilding   string         `json:"requiresBuilding,omitempty"`
}

type Content struct {
	Buildings []BuildingDef
	Events    []EventDef
}

func (c Content) FindBuilding(id string) (BuildingDef, bool) {
	for _, b := range c.Buildings {
		if b.ID == id {
			return b, true
		}
	}
	return BuildingDef{}, false
}

type Building struct {
	DefID   string
	Level   int
	Damaged bool
}

type State struct {
	Day            int
	Power          int
	Food           int
	Morale         int
	Credits        int
	Population     int
	PopulationCap  int
	BeaconParts    int
	MaxBeaconParts int
	Buildings      map[string]Building
	Log            []string
	KeyMoments     []string
	ScenarioID     string
	DifficultyID   string
	SurvivalWinAfterDay   int
	EventGateSkipAbove    int
	DailyCreditsIncomeDelta int
	Content        Content
	SessionLog     Recorder
	Seed           int64
	rngDrawMods    []int
	rng            *rand.Rand
	GameOver       bool
	Won            bool
	Message        string
	WarningLevels  map[string]WarningSeverity
	MilestonesSeen map[string]bool
	MinPowerSeen   int
	MinFoodSeen    int
	MinMoraleSeen  int
}
