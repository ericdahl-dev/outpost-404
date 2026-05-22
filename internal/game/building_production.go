package game

import (
	"fmt"
	"sort"
	"strings"
)

// FormatDailyProductionNote summarizes per-day output for build menu/UI.
func FormatDailyProductionNote(def BuildingDef) string {
	if len(def.DailyEffects) == 0 {
		return ""
	}
	parts := make([]string, 0, len(def.DailyEffects))
	for _, key := range dailyEffectKeys(def.DailyEffects) {
		amount := def.DailyEffects[key]
		if amount == 0 {
			continue
		}
		parts = append(parts, fmt.Sprintf("%+d %s/lv", amount, dailyEffectLabel(key)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "Daily: " + strings.Join(parts, ", ")
}

func dailyEffectKeys(effects map[string]int) []string {
	keys := make([]string, 0, len(effects))
	for k := range effects {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func dailyEffectLabel(key string) string {
	switch Resource(key) {
	case ResourcePower:
		return "power"
	case ResourceFood:
		return "food"
	case ResourceMorale:
		return "morale"
	case ResourceCredits:
		return "credits"
	default:
		return key
	}
}
