package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.Started {
		if m.Screen == screenTitle {
			return m.titleView()
		}
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

	left := boxStyle.Width(layout.LeftWidth).Render(RenderResourcePanel(m.State, layout.LeftWidth))

	middle := boxStyle.Width(layout.MiddleWidth).Render(RenderOutpostPanel(m.State, layout.MiddleWidth))

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
		"Log prefixes: ! alert  + gain  $ trade  * milestone  > event  · system",
	}
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(layout.BoxWidth).Render(strings.Join(help, "\n"))
}

