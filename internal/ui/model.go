package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/ericdahl/outpost-404/internal/game"
)

type screen int

const (
	screenTitle screen = iota
	screenNewGame
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
	SessionLogPath      string
	AutosavePath        string
	CanContinue         bool
	AwaitingOverwrite   bool
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
		Screen:          screenTitle,
		ScenarioIndex:   0,
		DifficultyIndex: diffIdx,
		TermWidth:       defaultTermWidth,
		TermHeight:      defaultTermHeight,
	}
	m.LogViewport = newLogViewport(m.TermWidth, m.TermHeight)
	if path, err := game.DefaultAutosavePath(); err == nil {
		m.AutosavePath = path
		m.CanContinue = game.AutosaveExists(path)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}
