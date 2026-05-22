package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ericdahl/outpost-404/internal/game"
	"github.com/ericdahl/outpost-404/internal/ui"
)

var version = "dev"

func main() {
	logFlag := flag.String("log", "", "JSONL session log path (empty: default cache dir; off: disable)")
	seedFlag := flag.Int64("seed", 0, "RNG seed for random events (0: random each run)")
	scenarioFlag := flag.String("scenario", "standard", "Scenario id for -simulate (standard, first_landing, dust_season, silent_colony, beacon_rush)")
	difficultyFlag := flag.String("difficulty", "normal", "Difficulty id for -simulate (easy, normal, hard)")
	replayFlag := flag.String("replay", "", "Replay a JSONL session log and verify snapshots (no TUI)")
	simulateFlag := flag.String("simulate", "", "Run a JSON sim script headlessly (no TUI)")
	seedsFlag := flag.String("seeds", "", "Comma-separated seeds for -simulate sweep (overrides -seed)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	content, err := loadContent()
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

	if *simulateFlag != "" {
		if err := runSimulate(content, *simulateFlag, *seedFlag, *seedsFlag, *scenarioFlag, *difficultyFlag); err != nil {
			fmt.Fprintf(os.Stderr, "simulate failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	profiles, err := loadRunProfiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load scenarios: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(content, profiles)
	model.SessionLogPath = resolveLogPath(*logFlag)

	program := tea.NewProgram(model, tea.WithAltScreen())

	final, err := program.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "outpost crashed: %v\n", err)
		os.Exit(1)
	}
	if m, ok := final.(ui.Model); ok {
		if m.Started {
			m.State.EndSession()
		}
	}
}

func loadRunProfiles() (game.RunProfiles, error) {
	if _, err := os.Stat("data/scenarios.json"); err == nil {
		return game.LoadRunProfiles("data")
	}
	return game.LoadEmbeddedRunProfiles()
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

func loadContent() (game.Content, error) {
	if _, err := os.Stat("data/buildings.json"); err == nil {
		return game.LoadContent("data")
	}
	return game.LoadEmbeddedContent()
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
