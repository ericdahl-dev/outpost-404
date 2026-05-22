package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) View() string {
	switch m.Screen {
	case screenBuild:
		return m.buildView()
	case screenHelp:
		return m.helpView()
	default:
		return m.mainView()
	}
}

func (m Model) mainView() string {
	if m.State.GameOver {
		heading := badStyle.Render("OUTPOST COLLAPSED")
		if m.State.Won {
			heading = goodStyle.Render("MISSION COMPLETE")
		}
		return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(72).Render(heading+"\n\n"+m.State.Message+"\n\nPress r to restart or q to quit.")
	}

	left := boxStyle.Width(36).Render(strings.Join([]string{
		fmt.Sprintf("Day: %d / %d", m.State.Day, game.SurvivalWinAfterDay),
		bar("Power", m.State.Power),
		bar("Food", m.State.Food),
		bar("Morale", m.State.Morale),
		fmt.Sprintf("Credits: %d", m.State.Credits),
		fmt.Sprintf("Population: %d / %d", m.State.Population, m.State.PopulationCap),
		fmt.Sprintf("Signal Beacon: %d / %d", m.State.BeaconParts, m.State.MaxBeaconParts),
	}, "\n"))

	buildings := []string{"Facilities"}
	for _, def := range m.State.Content.Buildings {
		buildings = append(buildings, fmt.Sprintf("%-14s Lv. %d/%d", def.Name, m.State.BuildingLevel(def.ID), def.MaxLevel))
	}
	middle := boxStyle.Width(34).Render(strings.Join(buildings, "\n"))

	logPanel := boxStyle.Width(m.LogViewport.Width + 2).Render(
		"Event Log\n" + mutedStyle.Render("↑↓ scroll") + "\n" + m.LogViewport.View(),
	)

	actions := "[B] Build/Upgrade  [R] Repair  [T] Trade  [S] Signal Beacon  [N/Space] Next Day  [?] Help  [Q] Quit"
	return titleStyle.Render("Outpost 404") + "\n" + mutedStyle.Render("A tiny terminal base builder by default, a future colony sim by design.") + "\n\n" + lipgloss.JoinHorizontal(lipgloss.Top, left, middle, logPanel) + "\n\n" + actions
}

func (m Model) buildView() string {
	hint := mutedStyle.Render("j/k or arrows · enter build · esc back")
	body := boxStyle.Width(m.BuildList.Width() + 2).Render(m.BuildList.View())
	return titleStyle.Render("Outpost 404") + "\n\n" + hint + "\n\n" + body
}

func (m Model) helpView() string {
	help := []string{
		fmt.Sprintf("Goal: survive %d days or complete 5 Signal Beacon parts.", game.SurvivalWinAfterDay),
		"",
		"Every day consumes power and food. Low resources damage morale.",
		"Build facilities to stabilize the colony, then spend power and credits on the beacon.",
		"Random events can help or hurt. Keep reserves.",
		"",
		"Keys: b build, r repair, t trade, s beacon, n/space next day, esc back, q quit.",
		"Event log scrolls with arrow keys on the main screen.",
	}
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(80).Render(strings.Join(help, "\n"))
}

func bar(label string, value int) string {
	width := 12
	filled := value * width / 100
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	cells := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	rendered := fmt.Sprintf("%-7s %s %3d%%", label+":", cells, value)
	if value <= 20 {
		return badStyle.Render(rendered)
	}
	if value <= 40 {
		return warnStyle.Render(rendered)
	}
	return rendered
}
