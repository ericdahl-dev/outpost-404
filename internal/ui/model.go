package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/ericdahl/outpost-404/internal/game"
)

type screen int

const (
	screenNewGame screen = iota
	screenMain
	screenBuild
	screenHelp
)

type Model struct {
	Content         game.Content
	Profiles        game.RunProfiles
	State           game.State
	Started         bool
	ScenarioIndex   int
	DifficultyIndex int
	Screen          screen
	BuildList       list.Model
	LogViewport     viewport.Model
	TermWidth       int
	TermHeight      int
	SessionLogPath  string
}

func NewModel(content game.Content, profiles game.RunProfiles) Model {
	diffIdx := 0
	for i, d := range profiles.Difficulties {
		if d.ID == "normal" {
			diffIdx = i
			break
		}
	}
	m := Model{
		Content:         content,
		Profiles:        profiles,
		Screen:          screenNewGame,
		ScenarioIndex:   0,
		DifficultyIndex: diffIdx,
		TermWidth:       defaultTermWidth,
		TermHeight:      defaultTermHeight,
	}
	m.LogViewport = newLogViewport(m.TermWidth, m.TermHeight)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}
