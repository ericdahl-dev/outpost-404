package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ericdahl/outpost-404/internal/game"
)

func (m Model) View() string {
	if !m.Started {
		return m.newGameView()
	}
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
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
	if m.State.GameOver {
		heading := badStyle.Render("OUTPOST COLLAPSED")
		if m.State.Won {
			heading = goodStyle.Render("MISSION COMPLETE")
		}
		boxW := boxWidth(m.TermWidth, 72)
		body := m.State.Message
		if summary := m.State.SessionSummary(); len(summary) > 0 {
			body += "\n\n" + strings.Join(summary, "\n")
		}
		return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(boxW).Render(heading+"\n\n"+body+"\n\nPress r to restart or q to quit.")
	}

	leftLines := []string{
		fmt.Sprintf("Day: %d / %d", m.State.Day, m.State.SurvivalWinTarget()),
		bar("Power", m.State.Power, layout.LeftWidth),
		bar("Food", m.State.Food, layout.LeftWidth),
		bar("Morale", m.State.Morale, layout.LeftWidth),
		fmt.Sprintf("Credits: %d", m.State.Credits),
		fmt.Sprintf("Population: %d / %d", m.State.Population, m.State.PopulationCap),
		fmt.Sprintf("Signal Beacon: %d / %d", m.State.BeaconParts, m.State.MaxBeaconParts),
	}
	for _, w := range game.ActiveWarnings(m.State) {
		line := "! " + w.Message
		if w.Severity >= game.SeverityCritical {
			leftLines = append(leftLines, badStyle.Render(line))
		} else {
			leftLines = append(leftLines, warnStyle.Render(line))
		}
	}
	left := boxStyle.Width(layout.LeftWidth).Render(strings.Join(leftLines, "\n"))

	buildings := []string{"Facilities"}
	nameW := facilityNameWidth(layout.MiddleWidth)
	for _, def := range m.State.Content.Buildings {
		lvl := m.State.BuildingLevel(def.ID)
		line := fmt.Sprintf("%-*s Lv. %d/%d", nameW, def.Name, lvl, def.MaxLevel)
		if b, ok := m.State.Buildings[def.ID]; ok && b.Damaged {
			line += " " + badStyle.Render("DAMAGED")
		}
		if lvl > 0 {
			if note := game.FormatDailyProductionNote(def); note != "" {
				line += "  " + mutedStyle.Render(note)
			}
		}
		buildings = append(buildings, line)
	}
	middle := boxStyle.Width(layout.MiddleWidth).Render(strings.Join(buildings, "\n"))

	logPanel := boxStyle.Width(m.LogViewport.Width + 2).Render(
		"Event Log\n" + mutedStyle.Render("↑↓ scroll") + "\n" + m.LogViewport.View(),
	)

	var panels string
	if layout.Stacked {
		panels = lipgloss.JoinVertical(lipgloss.Left, left, middle, logPanel)
	} else {
		panels = lipgloss.JoinHorizontal(lipgloss.Top, left, middle, logPanel)
	}

	return titleStyle.Render("Outpost 404") + "\n" + mutedStyle.Render("A tiny terminal base builder by default, a future colony sim by design.") + "\n\n" + panels + "\n\n" + mainActions(layout.CompactKeys)
}

func facilityNameWidth(panelWidth int) int {
	w := panelWidth - 12
	return clamp(8, w, 14)
}

func mainActions(compact bool) string {
	if compact {
		return "[B] Build  [R] Repair  [T] Trade  [S] Beacon  [N] Next  [?] Help  [Q] Quit"
	}
	return "[B] Build/Upgrade  [R] Repair  [T] Trade  [S] Signal Beacon  [N/Space] Next Day  [?] Help  [Q] Quit"
}

func (m Model) buildView() string {
	hint := mutedStyle.Render("j/k or arrows · enter build · esc back")
	body := boxStyle.Width(m.BuildList.Width() + 2).Render(m.BuildList.View())
	return titleStyle.Render("Outpost 404") + "\n\n" + hint + "\n\n" + body
}

func (m Model) helpView() string {
	help := []string{
		fmt.Sprintf("Goal: survive %d days or complete %d Signal Beacon parts.", m.State.SurvivalWinTarget(), m.State.MaxBeaconParts),
		"",
		"Every day consumes power and food. Low resources damage morale.",
		"Build facilities to stabilize the colony, then spend power and credits on the beacon.",
		"Random events can help or hurt. Keep reserves.",
		"",
		"Keys: b build, r repair, t trade, s beacon, n/space next day, esc back, q quit.",
		"Event log scrolls with arrow keys on the main screen.",
	}
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(layout.BoxWidth).Render(strings.Join(help, "\n"))
}

func bar(label string, value int, panelWidth int) string {
	width := clamp(8, panelWidth-12, 12)
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
