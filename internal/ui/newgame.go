package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) newGameView() string {
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
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
		mutedStyle.Render("[←/→] scenario  [↑/↓] difficulty  [Enter] start  [Q] quit"),
	}
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(layout.BoxWidth).Render(strings.Join(lines, "\n"))
}

func (m *Model) startRun() {
	seed := game.RandomSeed()
	sc := m.Profiles.Scenarios[m.ScenarioIndex].ID
	diffID := m.Profiles.Difficulties[m.DifficultyIndex].ID
	m.State = game.NewRun(m.Content, m.Profiles, seed, sc, diffID)
	m.Started = true
	m.Screen = screenMain
	m.BuildList = newBuildList(m.State, m.TermWidth)
	m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
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
