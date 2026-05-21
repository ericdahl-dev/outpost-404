package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

func main() {
	logFlag := flag.String("log", "", "JSONL session log path (empty: default cache dir; off: disable)")
	seedFlag := flag.Int64("seed", 0, "RNG seed for random events (0: random each run)")
	replayFlag := flag.String("replay", "", "Replay a JSONL session log and verify snapshots (no TUI)")
	flag.Parse()

	content, err := game.LoadContent("data")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load game content: %v\n", err)
		os.Exit(1)
	}

	if *replayFlag != "" {
		if err := runReplay(content, *replayFlag); err != nil {
			fmt.Fprintf(os.Stderr, "replay failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	var state game.State
	if *seedFlag != 0 {
		state = game.NewStateWithSeed(content, *seedFlag)
		fmt.Fprintf(os.Stderr, "seed: %d\n", *seedFlag)
	} else {
		state = game.NewState(content)
	}

	model := ui.NewModel(state)

	if err := attachSessionLog(&model.State, resolveLogPath(*logFlag)); err != nil {
		fmt.Fprintf(os.Stderr, "session logging disabled: %v\n", err)
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	final, err := program.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "outpost crashed: %v\n", err)
		os.Exit(1)
	}
	if m, ok := final.(ui.Model); ok {
		m.State.EndSession()
	}
}

func runReplay(content game.Content, path string) error {
	entries, err := game.LoadSessionLog(path)
	if err != nil {
		return err
	}
	final, err := game.ReplaySession(content, entries)
	if err != nil {
		return err
	}
	fmt.Printf("replay ok: day=%d won=%v game_over=%v beacon=%d/%d power=%d food=%d morale=%d credits=%d\n",
		final.Day, final.Won, final.GameOver, final.BeaconParts, final.MaxBeaconParts,
		final.Power, final.Food, final.Morale, final.Credits)
	return nil
}

func resolveLogPath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("OUTPOST_LOG")
}

func attachSessionLog(state *game.State, path string) error {
	if path == "off" {
		return nil
	}
	logger, err := game.AttachSessionLog(state, path)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "session log: %s\n", logger.Path)
	return nil
}
