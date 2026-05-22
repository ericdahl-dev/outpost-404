package game

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"

	outpostdata "github.com/ericdahl/outpost-404/data"
)

// StartingResources overrides default colony vitals at run start.
type StartingResources struct {
	Power         int `json:"power,omitempty"`
	Food          int `json:"food,omitempty"`
	Morale        int `json:"morale,omitempty"`
	Credits       int `json:"credits,omitempty"`
	Population    int `json:"population,omitempty"`
	PopulationCap int `json:"populationCap,omitempty"`
}

// RunModifiers tune pressure for a scenario or difficulty layer.
type RunModifiers struct {
	EventGateSkipAbove      int `json:"eventGateSkipAbove,omitempty"`
	DailyCreditsIncomeDelta int `json:"dailyCreditsIncomeDelta,omitempty"`
}

// ScenarioDef is a selectable run profile from data/scenarios.json.
type ScenarioDef struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Description         string            `json:"description,omitempty"`
	Starting            StartingResources `json:"starting,omitempty"`
	MaxBeaconParts      int               `json:"maxBeaconParts,omitempty"`
	SurvivalWinAfterDay int               `json:"survivalWinAfterDay,omitempty"`
	Modifiers           RunModifiers      `json:"modifiers,omitempty"`
}

// DifficultyDef stacks pressure on top of the chosen scenario.
type DifficultyDef struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Modifiers   RunModifiers `json:"modifiers,omitempty"`
}

// RunProfiles holds scenarios and difficulties loaded from JSON.
type RunProfiles struct {
	Scenarios    []ScenarioDef   `json:"scenarios"`
	Difficulties []DifficultyDef `json:"difficulties"`
}

// RunSetup selects scenario and difficulty for a new run.
type RunSetup struct {
	Seed         int64
	ScenarioID   string
	DifficultyID string
}

// LoadEmbeddedRunProfiles loads scenarios.json bundled with the binary.
func LoadEmbeddedRunProfiles() (RunProfiles, error) {
	return LoadRunProfilesFS(outpostdata.Files)
}

// LoadRunProfiles loads scenarios.json from dir.
func LoadRunProfiles(dir string) (RunProfiles, error) {
	return LoadRunProfilesFS(os.DirFS(dir))
}

// LoadRunProfilesFS reads scenarios.json from fsys.
func LoadRunProfilesFS(fsys fs.FS) (RunProfiles, error) {
	profiles, err := loadJSONFS[RunProfiles](fsys, "scenarios.json")
	if err != nil {
		return RunProfiles{}, err
	}
	if err := validateRunProfiles(profiles); err != nil {
		return RunProfiles{}, err
	}
	return profiles, nil
}

func validateRunProfiles(p RunProfiles) error {
	if len(p.Scenarios) < 3 {
		return fmt.Errorf("scenarios.json: need at least 3 scenarios, got %d", len(p.Scenarios))
	}
	if len(p.Difficulties) < 3 {
		return fmt.Errorf("scenarios.json: need at least 3 difficulties, got %d", len(p.Difficulties))
	}
	seen := map[string]bool{}
	for _, sc := range p.Scenarios {
		if sc.ID == "" {
			return fmt.Errorf("scenario missing id")
		}
		if seen[sc.ID] {
			return fmt.Errorf("duplicate scenario id %q", sc.ID)
		}
		seen[sc.ID] = true
	}
	seen = map[string]bool{}
	for _, d := range p.Difficulties {
		if d.ID == "" {
			return fmt.Errorf("difficulty missing id")
		}
		if seen[d.ID] {
			return fmt.Errorf("duplicate difficulty id %q", d.ID)
		}
		seen[d.ID] = true
	}
	return nil
}

func (p RunProfiles) FindScenario(id string) (ScenarioDef, bool) {
	for _, sc := range p.Scenarios {
		if sc.ID == id {
			return sc, true
		}
	}
	return ScenarioDef{}, false
}

func (p RunProfiles) FindDifficulty(id string) (DifficultyDef, bool) {
	for _, d := range p.Difficulties {
		if d.ID == id {
			return d, true
		}
	}
	return DifficultyDef{}, false
}

func mergeRunModifiers(base, overlay RunModifiers) RunModifiers {
	out := base
	if overlay.EventGateSkipAbove != 0 {
		out.EventGateSkipAbove = overlay.EventGateSkipAbove
	}
	out.DailyCreditsIncomeDelta += overlay.DailyCreditsIncomeDelta
	return out
}

func defaultStartingResources() StartingResources {
	return StartingResources{
		Power: 65, Food: 60, Morale: 70, Credits: 180,
		Population: 8, PopulationCap: 10,
	}
}

func applyStartingResources(dst *StartingResources, override StartingResources) {
	if override.Power > 0 {
		dst.Power = override.Power
	}
	if override.Food > 0 {
		dst.Food = override.Food
	}
	if override.Morale > 0 {
		dst.Morale = override.Morale
	}
	if override.Credits > 0 {
		dst.Credits = override.Credits
	}
	if override.Population > 0 {
		dst.Population = override.Population
	}
	if override.PopulationCap > 0 {
		dst.PopulationCap = override.PopulationCap
	}
}

// NewRun starts a colony with scenario + difficulty applied (defaults: standard, normal).
func NewRun(content Content, profiles RunProfiles, seed int64, scenarioID, difficultyID string) State {
	if scenarioID == "" {
		scenarioID = "standard"
	}
	if difficultyID == "" {
		difficultyID = "normal"
	}
	sc, ok := profiles.FindScenario(scenarioID)
	if !ok {
		sc, _ = profiles.FindScenario("standard")
		scenarioID = "standard"
	}
	diff, ok := profiles.FindDifficulty(difficultyID)
	if !ok {
		diff, _ = profiles.FindDifficulty("normal")
		difficultyID = "normal"
	}

	start := defaultStartingResources()
	applyStartingResources(&start, sc.Starting)

	s := State{
		Day:            1,
		Power:          start.Power,
		Food:           start.Food,
		Morale:         start.Morale,
		Credits:        start.Credits,
		Population:     start.Population,
		PopulationCap:  start.PopulationCap,
		BeaconParts:    0,
		MaxBeaconParts: 5,
		Buildings:      map[string]Building{},
		Content:        content,
		Log:            []string{},
		ScenarioID:     scenarioID,
		DifficultyID:   difficultyID,
		SurvivalWinAfterDay: SurvivalWinAfterDay,
		EventGateSkipAbove:  RandomEventRollSkipAbove,
	}
	if sc.MaxBeaconParts > 0 {
		s.MaxBeaconParts = sc.MaxBeaconParts
	}
	if sc.SurvivalWinAfterDay > 0 {
		s.SurvivalWinAfterDay = sc.SurvivalWinAfterDay
	}

	mods := mergeRunModifiers(RunModifiers{EventGateSkipAbove: RandomEventRollSkipAbove}, sc.Modifiers)
	mods = mergeRunModifiers(mods, diff.Modifiers)
	if mods.EventGateSkipAbove != 0 {
		s.EventGateSkipAbove = mods.EventGateSkipAbove
	}
	s.DailyCreditsIncomeDelta = mods.DailyCreditsIncomeDelta

	s.Seed = seed
	if seed != 0 {
		s.rng = rand.New(rand.NewSource(seed))
	}

	s.AddLogKind(LogMilestone, "Welcome to Outpost 404. Keep the systems online and finish the Signal Beacon.")
	s.AddLogKind(LogPlain, fmt.Sprintf("Scenario: %s · %s. Survive %d days or complete %d beacon parts.",
		sc.Name, diff.Name, s.SurvivalWinAfterDay, s.MaxBeaconParts))
	return s
}

// SurvivalWinTarget is the last day before survival victory (default 45).
func (s State) SurvivalWinTarget() int {
	return s.survivalWinAfterDay()
}

func (s State) survivalWinAfterDay() int {
	if s.SurvivalWinAfterDay > 0 {
		return s.SurvivalWinAfterDay
	}
	return SurvivalWinAfterDay
}
