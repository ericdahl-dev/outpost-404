package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const titleTagline = "Remote colony operator console — keep the outpost alive and finish the Signal Beacon."

func (m Model) titleView() string {
	layout := MainLayoutFor(m.TermWidth, m.TermHeight)
	w := layout.BoxWidth
	if w > 72 {
		w = 72
	}

	lines := []string{
		"",
		titleBanner(),
		"",
		mutedStyle.Render(titleTagline),
		"",
		"───",
		"",
	}
	if m.CanContinue {
		lines = append(lines, goodStyle.Render("[ C ] Continue saved run"))
	}
	lines = append(lines,
		"[ S ] Start new run",
		"[ Q ] Quit",
		"",
		mutedStyle.Render("Press a highlighted key"),
	)
	body := boxStyle.Width(w).Render(strings.Join(lines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(" OUTPOST 404 "),
		"",
		body,
	)
}

func titleBanner() string {
	return mutedStyle.Render(`     ___       _        _     ___   ___  
    / _ \ ___ | |_ __ _| |_  / _ \ / _ \ 
   | | | / _ \| __/ _` + "`" + ` | __| | | | | | | |
   | |_| | (_) | || (_| | |_  | |_| | |_| |
    \___/ \___/ \__\__,_|\__|  \___/ \___/ `)
}
