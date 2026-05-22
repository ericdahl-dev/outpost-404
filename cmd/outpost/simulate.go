package main

import (
	"fmt"
	"os"

	"github.com/ericdahl/outpost-404/internal/game"
)

func runSimulate(content game.Content, scriptPath string, seedFlag int64, seedsFlag string, scenarioFlag, difficultyFlag string) error {
	profiles, err := loadRunProfiles()
	if err != nil {
		return err
	}
	script, err := game.LoadSimScript(scriptPath)
	if err != nil {
		return err
	}

	sweep, err := game.ParseSeedList(seedsFlag)
	if err != nil {
		return err
	}

	if len(sweep) > 0 {
		return runSimSweep(content, profiles, script, sweep, scenarioFlag, difficultyFlag)
	}

	seed, err := game.ResolveSimSeed(script.SeedSet, script.Seed, seedFlag)
	if err != nil {
		return err
	}
	setup := game.RunSetup{
		Seed:         seed,
		ScenarioID:   coalesce(script.Scenario, scenarioFlag),
		DifficultyID: coalesce(script.Difficulty, difficultyFlag),
	}
	final, err := game.SimulateRun(content, setup, script.Actions)
	if err != nil {
		return err
	}
	fmt.Println(game.FormatSimOutcome(seed, final))
	return nil
}

func coalesce(scriptVal, flagVal string) string {
	if scriptVal != "" {
		return scriptVal
	}
	return flagVal
}

func runSimSweep(content game.Content, profiles game.RunProfiles, script game.SimScript, seeds []int64, scenarioFlag, difficultyFlag string) error {
	wins := 0
	setup := game.RunSetup{
		ScenarioID:   coalesce(script.Scenario, scenarioFlag),
		DifficultyID: coalesce(script.Difficulty, difficultyFlag),
	}
	for _, seed := range seeds {
		setup.Seed = seed
		final, err := game.SimulateRun(content, setup, script.Actions)
		if err != nil {
			return err
		}
		if final.Won {
			wins++
		}
		fmt.Println(game.FormatSimOutcome(seed, final))
	}
	fmt.Fprintf(os.Stderr, "sweep: %d/%d won\n", wins, len(seeds))
	return nil
}
