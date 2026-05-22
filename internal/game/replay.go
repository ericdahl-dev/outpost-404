package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// SimAction is one headless player action for Simulate.
type SimAction struct {
	Type       string `json:"type"`
	BuildingID string `json:"building_id,omitempty"`
}

// RandomSeed returns a seed suitable for new sessions.
func RandomSeed() int64 {
	return time.Now().UnixNano()
}

// NewStateWithSeed creates a standard/normal colony with deterministic random events.
func NewStateWithSeed(content Content, seed int64) State {
	profiles, err := LoadEmbeddedRunProfiles()
	if err != nil {
		s := NewRun(content, RunProfiles{}, seed, "standard", "normal")
		s.Seed = seed
		if seed != 0 {
			s.rng = rand.New(rand.NewSource(seed))
		}
		return s
	}
	return NewRun(content, profiles, seed, "standard", "normal")
}

// StateFromSnapshot restores stats and re-seeds RNG for replay.
func StateFromSnapshot(content Content, snap Snapshot, seed int64) State {
	buildings := make(map[string]Building, len(snap.Buildings))
	for id, level := range snap.Buildings {
		buildings[id] = Building{DefID: id, Level: level, Damaged: snap.Damaged[id]}
	}
	s := State{
		Day:            snap.Day,
		Power:          snap.Power,
		Food:           snap.Food,
		Morale:         snap.Morale,
		Credits:        snap.Credits,
		Population:     snap.Population,
		PopulationCap:  snap.PopulationCap,
		BeaconParts:    snap.BeaconParts,
		MaxBeaconParts: snap.MaxBeacon,
		Buildings:      buildings,
		Content:        content,
		Log:            []string{},
		Seed:           seed,
		GameOver:       snap.GameOver,
		Won:            snap.Won,
	}
	s.rng = rand.New(rand.NewSource(seed))
	return s
}

func (s *State) ensureRNG() {
	if s.rng == nil {
		if s.Seed == 0 {
			s.Seed = RandomSeed()
		}
		s.rng = rand.New(rand.NewSource(s.Seed))
	}
}

// actionRegistry maps action type strings to their State methods.
// Add new actions here; ApplySimAction derives from this map.
var actionRegistry = map[string]func(*State, SimAction){
	"build":    func(s *State, a SimAction) { s.Build(a.BuildingID) },
	"repair": func(s *State, a SimAction) {
		if a.BuildingID != "" {
			s.RepairBuilding(a.BuildingID)
		} else {
			s.Repair()
		}
	},
	"trade":    func(s *State, a SimAction) { s.Trade() },
	"beacon":   func(s *State, a SimAction) { s.WorkOnBeacon() },
	"damage":   func(s *State, a SimAction) { s.DamageBuilding(a.BuildingID) },
	"next_day": func(s *State, a SimAction) { s.NextDay() },
}

// ApplySimAction runs one action through the public game API.
func (s *State) ApplySimAction(a SimAction) {
	if fn, ok := actionRegistry[a.Type]; ok {
		fn(s, a)
	}
}

// Simulate runs a scripted session without the TUI (standard/normal).
func Simulate(content Content, seed int64, actions []SimAction) (State, error) {
	final, _, err := SimulateWithSnapshots(content, seed, actions)
	return final, err
}

// SimulateRun runs a scripted session with scenario and difficulty.
func SimulateRun(content Content, setup RunSetup, actions []SimAction) (State, error) {
	final, _, err := SimulateRunWithSnapshots(content, setup, actions)
	return final, err
}

// SimulateWithSnapshots runs a script and records a snapshot after each action (standard/normal).
func SimulateWithSnapshots(content Content, seed int64, actions []SimAction) (State, []Snapshot, error) {
	return SimulateRunWithSnapshots(content, RunSetup{Seed: seed, ScenarioID: "standard", DifficultyID: "normal"}, actions)
}

// SimulateRunWithSnapshots applies scenario/difficulty then runs actions.
func SimulateRunWithSnapshots(content Content, setup RunSetup, actions []SimAction) (State, []Snapshot, error) {
	profiles, err := LoadEmbeddedRunProfiles()
	if err != nil {
		return State{}, nil, fmt.Errorf("load run profiles: %w", err)
	}
	s := NewRun(content, profiles, setup.Seed, setup.ScenarioID, setup.DifficultyID)
	snaps := []Snapshot{s.snapshot()}
	for _, a := range actions {
		if s.GameOver {
			break
		}
		s.ApplySimAction(a)
		snaps = append(snaps, s.snapshot())
	}
	return s, snaps, nil
}

// LoadSessionLog reads JSONL session records from path.
func LoadSessionLog(path string) ([]LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open session log: %w", err)
	}
	defer func() { _ = f.Close() }()

	var entries []LogEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("parse session log line: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read session log: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("empty session log")
	}
	return entries, nil
}

// ReplaySession reapplies a logged session and verifies each step matches recorded snapshots.
func ReplaySession(content Content, entries []LogEntry) (State, error) {
	var start *LogEntry
	for i := range entries {
		if entries[i].Type == "session_start" {
			start = &entries[i]
			break
		}
	}
	if start == nil || start.Snapshot == nil {
		return State{}, fmt.Errorf("session log missing session_start snapshot")
	}

	seed, err := seedFromDetail(start.Detail)
	if err != nil {
		return State{}, fmt.Errorf("session log: %w", err)
	}

	s := StateFromSnapshot(content, *start.Snapshot, seed)
	for _, entry := range entries {
		if entry.Type == "session_start" {
			continue
		}
		if err := s.replayLogEntry(entry); err != nil {
			return s, err
		}
	}
	return s, nil
}

func (s *State) replayLogEntry(entry LogEntry) error {
	switch entry.Type {
	case "game_end":
		if entry.After != nil {
			s.restoreSnapshot(*entry.After)
		}
		if entry.Detail != nil {
			if msg, ok := entry.Detail["message"].(string); ok {
				s.Message = msg
			}
		}
		return nil
	case "next_day":
		s.replayNextDay(entry.Detail)
	default:
		if err := s.replayAction(entry); err != nil {
			return err
		}
	}
	if entry.After == nil {
		return nil
	}
	if diff := snapshotDiff(*entry.After, s.snapshot()); diff != "" {
		return fmt.Errorf("replay mismatch after %s (day %d): %s", entry.Type, entry.Day, diff)
	}
	return nil
}

func (s *State) replayAction(entry LogEntry) error {
	detail := cloneDetail(entry.Detail)
	switch entry.Type {
	case "build":
		id, _ := detail["building_id"].(string)
		s.buildWithDetail(detail, id)
	case "repair":
		id, _ := detail["building_id"].(string)
		s.repairWithDetail(detail, id)
	case "trade":
		s.tradeWithDetail(detail)
	case "beacon":
		s.beaconWithDetail(detail)
	case "damage":
		id, _ := detail["building_id"].(string)
		s.damageWithDetail(detail, id)
	default:
		return fmt.Errorf("unsupported log entry type %q", entry.Type)
	}
	return nil
}

func (s *State) restoreSnapshot(snap Snapshot) {
	buildings := make(map[string]Building, len(snap.Buildings))
	for id, level := range snap.Buildings {
		buildings[id] = Building{DefID: id, Level: level, Damaged: snap.Damaged[id]}
	}
	s.Day = snap.Day
	s.Power = snap.Power
	s.Food = snap.Food
	s.Morale = snap.Morale
	s.Credits = snap.Credits
	s.Population = snap.Population
	s.PopulationCap = snap.PopulationCap
	s.BeaconParts = snap.BeaconParts
	s.MaxBeaconParts = snap.MaxBeacon
	s.Buildings = buildings
	s.GameOver = snap.GameOver
	s.Won = snap.Won
}

func cloneDetail(detail map[string]any) map[string]any {
	if detail == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(detail))
	for k, v := range detail {
		out[k] = v
	}
	return out
}

func seedFromDetail(detail map[string]any) (int64, error) {
	if detail == nil {
		return 0, fmt.Errorf("missing seed in session_start (re-record with a current build)")
	}
	raw, ok := detail["seed"]
	if !ok {
		return 0, fmt.Errorf("missing seed in session_start (re-record with a current build)")
	}
	switch v := raw.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("invalid seed: %w", err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("invalid seed type %T", raw)
	}
}

func normalizeSnapshot(s Snapshot) Snapshot {
	s.Buildings = normalizeBuildings(s.Buildings)
	return s
}

func snapshotDiff(want, got Snapshot) string {
	want = normalizeSnapshot(want)
	got = normalizeSnapshot(got)
	var parts []string
	if want.Day != got.Day {
		parts = append(parts, fmt.Sprintf("day %d→%d", want.Day, got.Day))
	}
	if want.Power != got.Power {
		parts = append(parts, fmt.Sprintf("power %d→%d", want.Power, got.Power))
	}
	if want.Food != got.Food {
		parts = append(parts, fmt.Sprintf("food %d→%d", want.Food, got.Food))
	}
	if want.Morale != got.Morale {
		parts = append(parts, fmt.Sprintf("morale %d→%d", want.Morale, got.Morale))
	}
	if want.Credits != got.Credits {
		parts = append(parts, fmt.Sprintf("credits %d→%d", want.Credits, got.Credits))
	}
	if want.Population != got.Population {
		parts = append(parts, fmt.Sprintf("population %d→%d", want.Population, got.Population))
	}
	if want.PopulationCap != got.PopulationCap {
		parts = append(parts, fmt.Sprintf("population_cap %d→%d", want.PopulationCap, got.PopulationCap))
	}
	if want.BeaconParts != got.BeaconParts {
		parts = append(parts, fmt.Sprintf("beacon %d→%d", want.BeaconParts, got.BeaconParts))
	}
	if want.GameOver != got.GameOver {
		parts = append(parts, fmt.Sprintf("game_over %v→%v", want.GameOver, got.GameOver))
	}
	if want.Won != got.Won {
		parts = append(parts, fmt.Sprintf("won %v→%v", want.Won, got.Won))
	}
	if !mapsEqual(normalizeBuildings(want.Buildings), normalizeBuildings(got.Buildings)) {
		parts = append(parts, fmt.Sprintf("buildings %v→%v", want.Buildings, got.Buildings))
	}
	return strings.Join(parts, ", ")
}

func normalizeBuildings(m map[string]int) map[string]int {
	if len(m) == 0 {
		return map[string]int{}
	}
	return m
}

func mapsEqual(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
