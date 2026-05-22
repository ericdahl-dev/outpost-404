package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const AutosaveVersion = 1

// AutosaveFile is the on-disk resume document (separate from JSONL session logs).
type AutosaveFile struct {
	Version                 int            `json:"version"`
	Day                     int            `json:"day"`
	Power                   int            `json:"power"`
	Food                    int            `json:"food"`
	Morale                  int            `json:"morale"`
	Credits                 int            `json:"credits"`
	Population              int            `json:"population"`
	PopulationCap           int            `json:"population_cap"`
	BeaconParts             int            `json:"beacon_parts"`
	MaxBeaconParts          int            `json:"max_beacon_parts"`
	Buildings               map[string]SavedBuilding `json:"buildings"`
	Log                     []string       `json:"log"`
	KeyMoments              []string       `json:"key_moments,omitempty"`
	ScenarioID              string         `json:"scenario_id"`
	DifficultyID            string         `json:"difficulty_id"`
	SurvivalWinAfterDay     int            `json:"survival_win_after_day"`
	EventGateSkipAbove      int            `json:"event_gate_skip_above"`
	DailyCreditsIncomeDelta int            `json:"daily_credits_income_delta"`
	Seed                    int64          `json:"seed"`
	RNGDrawMods             []int          `json:"rng_draw_mods"`
	GameOver                bool           `json:"game_over"`
	Won                     bool           `json:"won"`
	Message                 string         `json:"message,omitempty"`
	WarningLevels           map[string]int `json:"warning_levels,omitempty"`
}

type SavedBuilding struct {
	Level   int  `json:"level"`
	Damaged bool `json:"damaged,omitempty"`
}

// DefaultAutosavePath returns the canonical autosave location under the user cache dir.
func DefaultAutosavePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "outpost-404", "autosave.json"), nil
}

// AutosaveExists reports whether path points to a readable autosave file.
func AutosaveExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func (s *State) toAutosaveFile() AutosaveFile {
	buildings := make(map[string]SavedBuilding, len(s.Buildings))
	for id, b := range s.Buildings {
		buildings[id] = SavedBuilding{Level: b.Level, Damaged: b.Damaged}
	}
	warn := make(map[string]int, len(s.WarningLevels))
	for id, sev := range s.WarningLevels {
		warn[id] = int(sev)
	}
	return AutosaveFile{
		Version:                 AutosaveVersion,
		Day:                     s.Day,
		Power:                   s.Power,
		Food:                    s.Food,
		Morale:                  s.Morale,
		Credits:                 s.Credits,
		Population:              s.Population,
		PopulationCap:           s.PopulationCap,
		BeaconParts:             s.BeaconParts,
		MaxBeaconParts:          s.MaxBeaconParts,
		Buildings:               buildings,
		Log:                     append([]string(nil), s.Log...),
		KeyMoments:              append([]string(nil), s.KeyMoments...),
		ScenarioID:              s.ScenarioID,
		DifficultyID:            s.DifficultyID,
		SurvivalWinAfterDay:     s.SurvivalWinAfterDay,
		EventGateSkipAbove:      s.EventGateSkipAbove,
		DailyCreditsIncomeDelta: s.DailyCreditsIncomeDelta,
		Seed:                    s.Seed,
		RNGDrawMods:             append([]int(nil), s.rngDrawMods...),
		GameOver:                s.GameOver,
		Won:                     s.Won,
		Message:                 s.Message,
		WarningLevels:           warn,
	}
}

// SaveAutosave writes state to path (creates parent dirs).
func SaveAutosave(s *State, path string) error {
	if path == "" {
		return fmt.Errorf("autosave path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir autosave dir: %w", err)
	}
	data, err := json.MarshalIndent(s.toAutosaveFile(), "", "  ")
	if err != nil {
		return fmt.Errorf("marshal autosave: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write autosave: %w", err)
	}
	return nil
}

// LoadAutosave restores a colony from JSON at path.
func LoadAutosave(path string, content Content, profiles RunProfiles) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, fmt.Errorf("no autosave at %s", path)
		}
		return State{}, fmt.Errorf("read autosave: %w", err)
	}
	var file AutosaveFile
	if err := json.Unmarshal(data, &file); err != nil {
		return State{}, fmt.Errorf("parse autosave: %w", err)
	}
	if file.Version != AutosaveVersion {
		return State{}, fmt.Errorf("unsupported autosave version %d", file.Version)
	}
	return stateFromAutosave(file, content, profiles)
}

func stateFromAutosave(file AutosaveFile, content Content, profiles RunProfiles) (State, error) {
	buildings := make(map[string]Building, len(file.Buildings))
	for id, b := range file.Buildings {
		buildings[id] = Building{DefID: id, Level: b.Level, Damaged: b.Damaged}
	}
	warn := make(map[string]WarningSeverity, len(file.WarningLevels))
	for id, sev := range file.WarningLevels {
		warn[id] = WarningSeverity(sev)
	}
	s := State{
		Day:                     file.Day,
		Power:                   file.Power,
		Food:                    file.Food,
		Morale:                  file.Morale,
		Credits:                 file.Credits,
		Population:              file.Population,
		PopulationCap:           file.PopulationCap,
		BeaconParts:             file.BeaconParts,
		MaxBeaconParts:          file.MaxBeaconParts,
		Buildings:               buildings,
		Log:                     file.Log,
		KeyMoments:              file.KeyMoments,
		ScenarioID:              file.ScenarioID,
		DifficultyID:            file.DifficultyID,
		SurvivalWinAfterDay:     file.SurvivalWinAfterDay,
		EventGateSkipAbove:      file.EventGateSkipAbove,
		DailyCreditsIncomeDelta: file.DailyCreditsIncomeDelta,
		Content:                 content,
		Seed:                    file.Seed,
		rngDrawMods:             append([]int(nil), file.RNGDrawMods...),
		GameOver:                file.GameOver,
		Won:                     file.Won,
		Message:                 file.Message,
		WarningLevels:           warn,
	}
	s.rng = replayRNG(s.Seed, s.rngDrawMods)
	if s.ScenarioID == "" {
		s.ScenarioID = "standard"
	}
	if s.DifficultyID == "" {
		s.DifficultyID = "normal"
	}
	_ = profiles
	return s, nil
}

// PersistAutosave writes to the default cache path.
func (s *State) PersistAutosave() error {
	path, err := DefaultAutosavePath()
	if err != nil {
		return err
	}
	return SaveAutosave(s, path)
}

// RemoveAutosave deletes the file at path (no error if missing).
func RemoveAutosave(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
