package ui

import (
	"fmt"
	"strings"

	"github.com/ericdahl/outpost-404/internal/game"
)

// RenderStatusHeader is the top title strip: game name, day, run profile, status badges.
func RenderStatusHeader(s game.State, profiles game.RunProfiles, _ int) string {
	scName := s.ScenarioID
	if sc, ok := profiles.FindScenario(s.ScenarioID); ok {
		scName = sc.Name
	}
	diffName := s.DifficultyID
	if d, ok := profiles.FindDifficulty(s.DifficultyID); ok {
		diffName = d.Name
	}
	line := fmt.Sprintf("Outpost 404 — Day %d / %d — %s · %s",
		s.Day, s.SurvivalWinTarget(), scName, diffName)
	if strip := formatStatusStrip(game.StatusBadges(s)); strip != "" {
		line += " — " + strings.TrimPrefix(strip, "Status: ")
	}
	return titleStyle.Render("Outpost 404") + "\n" + line
}
