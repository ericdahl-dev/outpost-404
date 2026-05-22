package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) newGameView() string {
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
	if m.AwaitingOverwrite {
		lines := []string{
			"Overwrite saved run?",
			"",
			"Starting a new run will replace your autosave.",
			"",
			mutedStyle.Render("[Y] overwrite and start  [N/Esc] cancel"),
		}
		return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(layout.BoxWidth).Render(strings.Join(lines, "\n"))
	}

	sc := m.Profiles.Scenarios[m.ScenarioIndex]
	diff := m.Profiles.Difficulties[m.DifficultyIndex]
	lines := []string{
		"New Run Setup",
		"",
		fmt.Sprintf("Scenario: %s", sc.Name),
		mutedStyle.Render(sc.Description),
		"",
		fmt.Sprintf("Difficulty: %s", diff.Name),
		mutedStyle.Render(diff.Description),
		"",
	}
	if m.CanContinue {
		lines = append(lines, warnStyle.Render("Saved run found — [C] continue"))
	}
	lines = append(lines, mutedStyle.Render("[←/→] scenario  [↑/↓] difficulty  [Enter] new run  [Q] quit"))
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(layout.BoxWidth).Render(strings.Join(lines, "\n"))
}

// ContinueFromAutosave loads the default autosave and skips setup (CLI -continue).
func (m *Model) ContinueFromAutosave() error {
	return m.continueRun()
}

func (m *Model) continueRun() error {
	if m.AutosavePath == "" {
		return fmt.Errorf("autosave path not configured")
	}
	s, err := game.LoadAutosave(m.AutosavePath, m.Content, m.Profiles)
	if err != nil {
		return err
	}
	m.State = s
	m.Started = true
	m.Screen = screenMain
	m.BuildList = newBuildList(m.State, m.TermWidth)
	m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	m.attachSessionLogIfConfigured()
	return nil
}

func (m *Model) startRun() {
	seed := game.RandomSeed()
	sc := m.Profiles.Scenarios[m.ScenarioIndex].ID
	diffID := m.Profiles.Difficulties[m.DifficultyIndex].ID
	m.State = game.NewRun(m.Content, m.Profiles, seed, sc, diffID)
	m.Started = true
	m.Screen = screenMain
	m.AwaitingOverwrite = false
	m.BuildList = newBuildList(m.State, m.TermWidth)
	m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	if m.AutosavePath != "" {
		_ = game.RemoveAutosave(m.AutosavePath)
		m.CanContinue = false
	}
	m.attachSessionLogIfConfigured()
}

func (m *Model) attachSessionLogIfConfigured() {
	if m.SessionLogPath == "off" {
		return
	}
	logger, err := game.AttachSessionLog(&m.State, m.SessionLogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "session logging disabled: %v\n", err)
		return
	}
	m.State.LogSessionStart()
	fmt.Fprintf(os.Stderr, "session log: %s\n", logger.Path)
}
