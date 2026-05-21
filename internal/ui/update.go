package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		m.LogViewport.Width = logViewportWidth(msg.Width)
		m.LogViewport.Height = logViewportHeight(msg.Height)
		m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
		m.BuildList.SetSize(buildListWidth(msg.Width), m.BuildList.Height())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
				hadLog := m.State.SessionLog != nil
				m.State.EndSession()
				m.State = game.NewState(m.State.Content)
				m.Screen = screenMain
				m.BuildList = newBuildList(m.State, m.TermWidth)
				m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
				if hadLog {
					if logger, err := game.AttachSessionLog(&m.State, ""); err != nil {
						fmt.Fprintf(os.Stderr, "session logging disabled: %v\n", err)
					} else {
						fmt.Fprintf(os.Stderr, "session log: %s\n", logger.Path)
					}
				}
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
