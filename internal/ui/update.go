package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		m.LogViewport.Width = LogViewportWidth(msg.Width)
		m.LogViewport.Height = LogViewportHeight(msg.Height)
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
		m.BuildList.SetSize(buildListWidth(msg.Width), m.BuildList.Height())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if !m.Started {
			return updateNewGame(m, msg)
		}

		switch msg.String() {
		case "?":
			if m.Screen == screenHelp {
				m.Screen = screenMain
			} else {
				m.Screen = screenHelp
			}
		case "esc":
			m.Screen = screenMain
		}

		if m.State.GameOver {
			switch msg.String() {
			case "r":
				m.State.EndSession()
				m.Started = false
				m.Screen = screenNewGame
				m.State = game.State{}
			}
			return m, nil
		}

		switch m.Screen {
		case screenMain:
			return updateMain(m, msg)
		case screenBuild:
			return updateBuild(m, msg)
		}
	}
	return m, nil
}

func updateNewGame(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		if m.ScenarioIndex > 0 {
			m.ScenarioIndex--
		}
	case "right", "l":
		if m.ScenarioIndex < len(m.Profiles.Scenarios)-1 {
			m.ScenarioIndex++
		}
	case "up", "k":
		if m.DifficultyIndex > 0 {
			m.DifficultyIndex--
		}
	case "down", "j":
		if m.DifficultyIndex < len(m.Profiles.Difficulties)-1 {
			m.DifficultyIndex++
		}
	case "enter":
		m.startRun()
	}
	return m, nil
}

func updateMain(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.LogViewport, cmd = m.LogViewport.Update(msg)

	switch msg.String() {
	case "b":
		m.Screen = screenBuild
		m.BuildList = newBuildList(m.State, m.TermWidth)
	case "n", " ":
		m.State.NextDay()
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	case "r":
		m.State.Repair()
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	case "t":
		m.State.Trade()
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	case "s":
		m.State.WorkOnBeacon()
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	}
	return m, cmd
}

func updateBuild(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "b", "esc":
		m.Screen = screenMain
		return m, nil
	}

	var cmd tea.Cmd
	m.BuildList, cmd = m.BuildList.Update(msg)

	if msg.String() == "enter" {
		if id, ok := selectedBuildingID(m.BuildList); ok {
			m.State.Build(id)
			m.BuildList = refreshBuildList(m.BuildList, m.State)
			m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
		}
	}
	return m, cmd
}
