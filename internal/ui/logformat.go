package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colony log prefix map (game AddLogKind → viewport styling). See docs/tui-log.md.
const (
	logPrefixDanger    = '!'
	logPrefixGain      = '+'
	logPrefixTrade     = '$'
	logPrefixMilestone = '*'
	logPrefixEvent     = '>'
	logPrefixSystem    = '·'
)

// FormatColonyLogLine styles a single log line for the event log viewport.
func FormatColonyLogLine(line string) string {
	if line == "" {
		return line
	}
	if len(line) < 2 || line[1] != ' ' {
		return line
	}
	switch line[0] {
	case logPrefixDanger:
		return badStyle.Render(line)
	case logPrefixGain:
		return goodStyle.Render(line)
	case logPrefixTrade:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(line)
	case logPrefixMilestone:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Bold(true).Render(line)
	case logPrefixEvent:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Render(line)
	case logPrefixSystem:
		return warnStyle.Render(line)
	default:
		return line
	}
}

// FormatColonyLogLines renders all log lines for the Bubbles viewport.
func FormatColonyLogLines(lines []string) string {
	if len(lines) == 0 {
		return mutedStyle.Render("Quiet shift. No new entries.")
	}
	formatted := make([]string, len(lines))
	for i, line := range lines {
		formatted[i] = FormatColonyLogLine(line)
	}
	return strings.Join(formatted, "\n")
}
