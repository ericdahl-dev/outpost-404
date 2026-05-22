package main

import (
	"fmt"
	"os"

	"github.com/ericdahl/outpost-404/internal/game"
)

func runSimulate(content game.Content, scriptPath string, seedFlag int64, seedsFlag string) error {
	script, err := game.LoadSimScript(scriptPath)
	if err != nil {
		return err
	}

	sweep, err := game.ParseSeedList(seedsFlag)
	if err != nil {
		return err
	}

	if len(sweep) > 0 {
		return runSimSweep(content, script, sweep)
	}

	seed, err := game.ResolveSimSeed(script.SeedSet, script.Seed, seedFlag)
	if err != nil {
		return err
	}
	final, err := game.Simulate(content, seed, script.Actions)
	if err != nil {
		return err
	}
	fmt.Println(game.FormatSimOutcome(seed, final))
	return nil
}

func runSimSweep(content game.Content, script game.SimScript, seeds []int64) error {
	wins := 0
	for _, seed := range seeds {
		final, err := game.Simulate(content, seed, script.Actions)
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
