package game

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SimScript is a headless action script, optionally with a default seed.
type SimScript struct {
	Seed    int64
	SeedSet bool
	Actions []SimAction
}

// LoadSimScript reads a JSON array of actions or {"seed": N, "actions": [...]}.
func LoadSimScript(path string) (SimScript, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SimScript{}, fmt.Errorf("read sim script: %w", err)
	}
	trim := strings.TrimSpace(string(data))
	if trim == "" {
		return SimScript{}, fmt.Errorf("sim script is empty")
	}

	var wrapped struct {
		Seed    *int64      `json:"seed"`
		Actions []SimAction `json:"actions"`
	}
	if trim[0] == '{' {
		if err := json.Unmarshal(data, &wrapped); err != nil {
			return SimScript{}, fmt.Errorf("parse sim script: %w", err)
		}
		if len(wrapped.Actions) == 0 {
			return SimScript{}, fmt.Errorf("sim script has no actions")
		}
		out := SimScript{Actions: wrapped.Actions}
		if wrapped.Seed != nil {
			out.Seed = *wrapped.Seed
			out.SeedSet = true
		}
		return out, nil
	}

	var actions []SimAction
	if err := json.Unmarshal(data, &actions); err != nil {
		return SimScript{}, fmt.Errorf("parse sim script: %w", err)
	}
	if len(actions) == 0 {
		return SimScript{}, fmt.Errorf("sim script has no actions")
	}
	return SimScript{Actions: actions}, nil
}

// ResolveSimSeed picks the seed for a run: -seed flag overrides script seed.
func ResolveSimSeed(scriptSeedSet bool, scriptSeed, flagSeed int64) (int64, error) {
	if flagSeed != 0 {
		return flagSeed, nil
	}
	if scriptSeedSet {
		return scriptSeed, nil
	}
	return 0, fmt.Errorf("seed required: pass -seed or include \"seed\" in the script JSON")
}

// FormatSimOutcome is one line suitable for logs and CLI output.
func FormatSimOutcome(seed int64, s State) string {
	return fmt.Sprintf("seed=%d day=%d won=%v game_over=%v beacon=%d/%d power=%d food=%d morale=%d credits=%d",
		seed, s.Day, s.Won, s.GameOver, s.BeaconParts, s.MaxBeaconParts,
		s.Power, s.Food, s.Morale, s.Credits)
}

// SimulateSeeds runs the same script across multiple seeds.
func SimulateSeeds(content Content, seeds []int64, actions []SimAction) ([]State, error) {
	if len(seeds) == 0 {
		return nil, fmt.Errorf("no seeds to simulate")
	}
	out := make([]State, 0, len(seeds))
	for _, seed := range seeds {
		s, err := Simulate(content, seed, actions)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// ParseSeedList parses comma-separated int64 seeds (e.g. "1,42,99").
func ParseSeedList(s string) ([]int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	seeds := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid seed %q", p)
		}
		seeds = append(seeds, n)
	}
	if len(seeds) == 0 {
		return nil, fmt.Errorf("no seeds in list")
	}
	return seeds, nil
}
