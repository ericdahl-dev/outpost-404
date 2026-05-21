package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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

// NewStateWithSeed creates a colony with deterministic random events.
func NewStateWithSeed(content Content, seed int64) State {
	s := newBareState(content)
	s.Seed = seed
	s.rng = rand.New(rand.NewSource(seed))
	return s
}

// StateFromSnapshot restores stats and re-seeds RNG for replay.
func StateFromSnapshot(content Content, snap Snapshot, seed int64) State {
	buildings := make(map[string]Building, len(snap.Buildings))
	for id, level := range snap.Buildings {
		buildings[id] = Building{DefID: id, Level: level}
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

// ApplySimAction runs one action through the public game API.
func (s *State) ApplySimAction(a SimAction) {
	switch a.Type {
	case "build":
		s.Build(a.BuildingID)
	case "repair":
		s.Repair()
	case "trade":
		s.Trade()
	case "beacon":
		s.WorkOnBeacon()
	case "next_day":
		s.NextDay()
	}
}

// Simulate runs a scripted session without the TUI.
func Simulate(content Content, seed int64, actions []SimAction) (State, error) {
	s := NewStateWithSeed(content, seed)
	for _, a := range actions {
		if s.GameOver {
			break
		}
		s.ApplySimAction(a)
	}
	return s, nil
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
		switch entry.Type {
		case "session_start", "game_end":
			continue
		}

		action, ok := simActionFromEntry(entry)
		if !ok {
			return s, fmt.Errorf("unsupported log entry type %q", entry.Type)
		}
		s.ApplySimAction(action)
		if entry.After == nil {
			continue
		}
		if diff := snapshotDiff(*entry.After, s.snapshot()); diff != "" {
			return s, fmt.Errorf("replay mismatch after %s (day %d): %s", entry.Type, entry.Day, diff)
		}
	}
	return s, nil
}

func simActionFromEntry(entry LogEntry) (SimAction, bool) {
	switch entry.Type {
	case "build":
		id, _ := entry.Detail["building_id"].(string)
		return SimAction{Type: "build", BuildingID: id}, true
	case "repair":
		return SimAction{Type: "repair"}, true
	case "trade":
		return SimAction{Type: "trade"}, true
	case "beacon":
		return SimAction{Type: "beacon"}, true
	case "next_day":
		return SimAction{Type: "next_day"}, true
	default:
		return SimAction{}, false
	}
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
