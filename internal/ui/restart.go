package ui

import "github.com/ericdahl/outpost-404/internal/game"

func (m *Model) restartSameSeed() {
	seed := m.State.Seed
	scenarioID := m.State.ScenarioID
	difficultyID := m.State.DifficultyID
	m.State.EndSession()
	m.State = game.NewRun(m.Content, m.Profiles, seed, scenarioID, difficultyID)
	m.Started = true
	m.Screen = screenMain
	m.BuildList = newBuildList(m.State, m.TermWidth)
	m.LogViewport = syncLogViewport(m.LogViewport, m.State.Log)
	if m.AutosavePath != "" {
		_ = game.RemoveAutosave(m.AutosavePath)
		m.CanContinue = false
	}
	m.attachSessionLogIfConfigured()
}
