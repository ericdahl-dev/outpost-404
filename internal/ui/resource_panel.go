package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ericdahl/outpost-404/internal/game"
)

// RenderResourcePanel draws colony vitals, projected daily deltas, and status badges.
func RenderResourcePanel(s game.State, panelWidth int) string {
	title := lipgloss.NewStyle().Bold(true).Render("Colony")
	deltas := s.ProjectedDailyDeltas()
	barW := resourceBarWidth(panelWidth)

	lines := []string{
		resourceBarLine("Power", s.Power, deltas.Power, barW, game.WarningPowerCriticalAt, game.WarningPowerUrgentAt),
		resourceBarLine("Food", s.Food, deltas.Food, barW, game.WarningFoodCriticalAt, game.WarningFoodUrgentAt),
		resourceBarLine("Morale", s.Morale, deltas.Morale, barW, game.WarningMoraleCriticalAt, game.WarningMoraleUrgentAt),
		fmt.Sprintf("Credits %-4d  Pop %d/%d  Beacon %d/%d",
			s.Credits, s.Population, s.PopulationCap, s.BeaconParts, s.MaxBeaconParts),
	}
	return title + "\n" + strings.Join(lines, "\n")
}

func resourceBarWidth(panelWidth int) int {
	return clamp(10, panelWidth-22, 16)
}

func resourceBarLine(label string, value, delta, barWidth, criticalAt, urgentAt int) string {
	bar := renderResourceBar(value, barWidth, urgentAt, criticalAt)
	deltaStr := formatDailyDelta(delta)
	tag := resourceSeverityTag(value, criticalAt, urgentAt)
	line := fmt.Sprintf("%-6s %s %3d%%  %s", label, bar, value, deltaStr)
	if tag != "" {
		line += "  " + tag
	}
	return styleResourceLine(line, value, criticalAt, urgentAt)
}

func renderResourceBar(value, width, urgentAt, criticalAt int) string {
	if width < 1 {
		width = 1
	}
	filled := value * width / 100
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	cells := make([]rune, width)
	for i := 0; i < width; i++ {
		switch {
		case i < filled:
			cells[i] = '█'
		case barThresholdIndex(i, width, urgentAt) || barThresholdIndex(i, width, criticalAt):
			cells[i] = '┊'
		default:
			cells[i] = '░'
		}
	}
	return string(cells)
}

func barThresholdIndex(cell, width, threshold int) bool {
	if threshold <= 0 || threshold >= 100 {
		return false
	}
	idx := threshold * width / 100
	return cell == idx
}

func formatDailyDelta(delta int) string {
	if delta > 0 {
		return fmt.Sprintf("+%d/day", delta)
	}
	if delta < 0 {
		return fmt.Sprintf("%d/day", delta)
	}
	return "±0/day"
}

func resourceSeverityTag(value, criticalAt, urgentAt int) string {
	if value <= criticalAt {
		return "CRITICAL"
	}
	if value <= urgentAt {
		return "LOW"
	}
	return ""
}

func styleResourceLine(line string, value, criticalAt, urgentAt int) string {
	if value <= criticalAt {
		return badStyle.Render(line)
	}
	if value <= urgentAt {
		return warnStyle.Render(line)
	}
	return line
}

func formatStatusStrip(badges []game.StatusBadge) string {
	if len(badges) == 0 {
		return ""
	}
	parts := make([]string, len(badges))
	for i, b := range badges {
		parts[i] = styleStatusBadge(b)
	}
	return "Status: " + strings.Join(parts, " · ")
}

func styleStatusBadge(b game.StatusBadge) string {
	switch b.Severity {
	case game.SeverityCritical:
		return badStyle.Render(b.Label)
	case game.SeverityUrgent:
		return warnStyle.Render(b.Label)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(b.Label)
	}
}
