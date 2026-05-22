package ui

import (
	"fmt"
	"strings"

	"github.com/ericdahl/outpost-404/internal/game"
)

const (
	winArtFull = `       /\
      /  \
     /____\
       ||
       ||
   SIGNAL SENT`

	lossArtFull = `    .-.
   ( x )
    '-'`

	winArtCompact  = "  /\\\n ====\n  ||"
	lossArtCompact = " .-.\n( x )"
)

// RenderEndScreen draws the win or loss presentation for game over.
func RenderEndScreen(r game.RunReport, termWidth int) string {
	boxW := boxWidth(termWidth, 72)
	var blocks []string

	if r.Won {
		blocks = append(blocks, endTitle("MISSION COMPLETE", true))
		blocks = append(blocks, pickArt(winArtFull, winArtCompact, termWidth))
		blocks = append(blocks, "", mutedStyle.Render("Outpost 404 is no longer alone."))
	} else {
		blocks = append(blocks, endTitle("OUTPOST COLLAPSED", false))
		blocks = append(blocks, pickArt(lossArtFull, lossArtCompact, termWidth))
		blocks = append(blocks, "", "Last transmission:", endCauseLine(r))
	}

	blocks = append(blocks, "")
	blocks = append(blocks, endStatsBlock(r)...)
	if risk := vitalRiskLine(r); risk != "" {
		blocks = append(blocks, risk)
	}
	blocks = append(blocks,
		fmt.Sprintf("Run: %s · %s · seed %d", r.ScenarioName, r.DifficultyName, r.Seed),
	)
	if len(r.KeyMoments) > 0 {
		blocks = append(blocks, "", "Key moments:")
		show := r.KeyMoments
		if len(show) > 3 {
			show = show[len(show)-3:]
		}
		for _, m := range show {
			blocks = append(blocks, "  · "+strings.TrimPrefix(m, "* "))
		}
	}
	blocks = append(blocks, "", endScreenActions(termWidth < 88))

	body := strings.Join(blocks, "\n")
	return titleStyle.Render("Outpost 404") + "\n\n" + boxStyle.Width(boxW).Render(body)
}

func endTitle(text string, won bool) string {
	if won {
		return goodStyle.Render(text) + "\n" + text
	}
	return badStyle.Render(text) + "\n" + text
}

func pickArt(full, compact string, termWidth int) string {
	if termWidth < 56 {
		return compact
	}
	return full
}

func endCauseLine(r game.RunReport) string {
	line := r.Cause
	if r.Day > 0 {
		line = fmt.Sprintf("%s on Day %d.", strings.TrimSuffix(r.Cause, "."), r.Day)
	}
	if r.Won {
		return line
	}
	return badStyle.Render(line) + "\n" + line
}

func endStatsBlock(r game.RunReport) []string {
	return []string{
		"--- Final status ---",
		fmt.Sprintf("Day %-4d  Beacon %d/%d  Facilities %d built",
			r.Day, r.BeaconParts, r.MaxBeaconParts, r.BuiltCount),
		fmt.Sprintf("Power %-3d  Food %-3d  Morale %-3d  Credits %-4d",
			r.Power, r.Food, r.Morale, r.Credits),
		fmt.Sprintf("Population %d/%d", r.Population, r.PopulationCap),
	}
}

func vitalRiskLine(r game.RunReport) string {
	var parts []string
	if r.LowestPower > 0 && r.LowestPower <= game.WarningPowerUrgentAt {
		parts = append(parts, fmt.Sprintf("lowest power %d%%", r.LowestPower))
	}
	if r.LowestFood > 0 && r.LowestFood <= game.WarningFoodUrgentAt {
		parts = append(parts, fmt.Sprintf("lowest food %d%%", r.LowestFood))
	}
	if r.LowestMorale > 0 && r.LowestMorale <= game.WarningMoraleUrgentAt {
		parts = append(parts, fmt.Sprintf("lowest morale %d%%", r.LowestMorale))
	}
	if len(parts) == 0 {
		return ""
	}
	return "Risks: " + strings.Join(parts, " · ")
}

func endScreenActions(compact bool) string {
	if compact {
		return "[R] New run  [S] Same seed  [Q] Quit"
	}
	return "[R] New run (title)  [S] Replay same seed  [Q] Quit"
}
