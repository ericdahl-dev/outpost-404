package game

import (
	"fmt"
	"strings"
)

// LogKind prefixes colony log lines for display (see CONTEXT.md).
type LogKind string

const (
	LogPlain     LogKind = ""
	LogDanger    LogKind = "!"
	LogGain      LogKind = "+"
	LogTrade     LogKind = "$"
	LogMilestone LogKind = "*"
	LogEvent     LogKind = ">"
	LogSystem    LogKind = "·"
)

const maxColonyLogLines = 12
const maxKeyMoments = 6

var effectSummaryOrder = []string{
	"power", "food", "morale", "credits", "population", "populationCap",
}

func resourceLabel(key string) string {
	switch key {
	case "power":
		return "Power"
	case "food":
		return "Food"
	case "morale":
		return "Morale"
	case "credits":
		return "Credits"
	case "population":
		return "Population"
	case "populationCap":
		return "Population cap"
	default:
		return key
	}
}

func formatResourceDelta(key string, amount int) string {
	sign := "+"
	if amount < 0 {
		sign = ""
	}
	return fmt.Sprintf("%s %s%d", resourceLabel(key), sign, amount)
}

func formatEffectSummary(effects map[string]int) string {
	if len(effects) == 0 {
		return ""
	}
	var parts []string
	for _, key := range effectSummaryOrder {
		amount, ok := effects[key]
		if !ok || amount == 0 {
			continue
		}
		parts = append(parts, formatResourceDelta(key, amount))
	}
	return strings.Join(parts, ", ")
}

func formatEventLogLine(event EventDef) string {
	flavor := strings.TrimSpace(event.Description)
	impact := formatEffectSummary(event.Effects)
	switch {
	case flavor != "" && impact != "":
		return event.Title + ": " + flavor + " " + impact + "."
	case flavor != "":
		return event.Title + ": " + flavor
	case impact != "":
		return event.Title + ": " + impact + "."
	default:
		return event.Title
	}
}

func (s *State) AddLogKind(kind LogKind, message string) {
	line := message
	if kind != LogPlain && message != "" && !strings.HasPrefix(message, string(kind)+" ") {
		line = string(kind) + " " + message
	}
	s.AddLog(line)
	if kind == LogMilestone {
		s.recordKeyMoment(line)
	}
}

func (s *State) recordKeyMoment(line string) {
	if line == "" {
		return
	}
	s.KeyMoments = append(s.KeyMoments, line)
	if len(s.KeyMoments) > maxKeyMoments {
		s.KeyMoments = s.KeyMoments[len(s.KeyMoments)-maxKeyMoments:]
	}
}

// SessionSummary returns end-of-run highlights for the game-over screen.
func (s State) SessionSummary() []string {
	lines := []string{"--- Run summary ---"}
	if len(s.KeyMoments) > 0 {
		lines = append(lines, "Key moments:")
		lines = append(lines, s.KeyMoments...)
	}
	lines = append(lines, fmt.Sprintf(
		"Final: day %d · beacon %d/%d · pop %d/%d",
		s.Day, s.BeaconParts, s.MaxBeaconParts, s.Population, s.PopulationCap,
	))
	return lines
}

// emitQuietBeat logs a telemetry line on days when no random event fires.
func (s *State) emitQuietBeat() {
	beats := s.Content.QuietBeats
	if len(beats) == 0 {
		return
	}
	s.ensureRNG()
	beat := beats[s.rngIntn(len(beats))]
	s.AddLogKind(LogSystem, beat)
}

func collapseCause(s State) string {
	switch {
	case s.Power <= 0:
		return "life support failed (power exhausted)"
	case s.Food <= 0:
		return "food reserves depleted"
	case s.Morale <= 0:
		return "morale collapsed"
	case s.Population <= 0:
		return "population lost"
	default:
		return "critical systems failed"
	}
}
