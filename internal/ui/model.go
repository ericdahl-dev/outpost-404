package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/ericdahl/outpost-404/internal/game"
)

type screen int

const (
	screenMain screen = iota
	screenBuild
	screenHelp
)

type Model struct {
	State      game.State
	Screen     screen
	BuildList  list.Model
	LogViewport viewport.Model
	TermWidth  int
	TermHeight int
}

func NewModel(state game.State) Model {
	m := Model{
		State:      state,
		Screen:     screenMain,
		TermWidth:  defaultTermWidth,
		TermHeight: defaultTermHeight,
	}
	m.BuildList = newBuildList(state, m.TermWidth)
	m.LogViewport = syncLogViewport(newLogViewport(m.TermWidth, m.TermHeight), state.Log)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}
