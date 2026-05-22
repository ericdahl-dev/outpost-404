package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ericdahl/outpost-404/internal/game"
)

// Schematic facility abbreviations (docs/tui-graphics-plan.md).
const (
	abbrSolar     = "SA"
	abbrHydro     = "HY"
	abbrHabitat   = "HB"
	abbrWorkshop  = "WS"
	abbrMedBay    = "MB"
	abbrRadio     = "RT"
	abbrBeacon    = "BE"
)

const powerOfflineThreshold = 20

type facilityDisplayStatus int

const (
	facilityUnbuilt facilityDisplayStatus = iota
	facilityBuilt
	facilityMax
	facilityDamaged
	facilityOffline
	facilityBeaconProgress
	facilityBeaconComplete
)

type facilitySpec struct {
	level       int
	maxLevel    int
	damaged     bool
	offline     bool
	beaconParts int
	maxBeacon   int
	isBeacon    bool
}

type schematicSlot struct {
	id   string
	abbr string
}

var (
	schematicFullLines = []string{
		"      {RT}      ",
		"       │       ",
		" {SA}──{HB}──{WS} ",
		"  │    │    │  ",
		" {HY}──{MB}──{BE} ",
	}
	schematicCompactOrder = []schematicSlot{
		{id: "radio_tower", abbr: abbrRadio},
		{id: "solar_array", abbr: abbrSolar},
		{id: "habitat", abbr: abbrHabitat},
		{id: "workshop", abbr: abbrWorkshop},
		{id: "hydroponics", abbr: abbrHydro},
		{id: "med_bay", abbr: abbrMedBay},
		{id: "signal_beacon", abbr: abbrBeacon},
	}
	schematicSlotByAbbr = map[string]schematicSlot{
		abbrRadio:   {id: "radio_tower", abbr: abbrRadio},
		abbrSolar:   {id: "solar_array", abbr: abbrSolar},
		abbrHabitat: {id: "habitat", abbr: abbrHabitat},
		abbrWorkshop: {id: "workshop", abbr: abbrWorkshop},
		abbrHydro:   {id: "hydroponics", abbr: abbrHydro},
		abbrMedBay:  {id: "med_bay", abbr: abbrMedBay},
		abbrBeacon:  {id: "signal_beacon", abbr: abbrBeacon},
	}
)

// RenderOutpostPanel draws the fixed-layout ASCII colony schematic for the main screen.
func RenderOutpostPanel(s game.State, panelWidth int) string {
	title := lipgloss.NewStyle().Bold(true).Render("Outpost")
	var body string
	if panelWidth < 34 {
		body = renderOutpostCompact(s)
	} else {
		body = renderOutpostFull(s)
	}
	return title + "\n" + body
}

func renderOutpostFull(s game.State) string {
	lines := make([]string, len(schematicFullLines))
	for i, line := range schematicFullLines {
		lines[i] = substituteSchematicLine(s, line)
	}
	return strings.Join(lines, "\n")
}

func substituteSchematicLine(s game.State, line string) string {
	out := line
	for abbr, slot := range schematicSlotByAbbr {
		token, status := facilityTokenForState(s, slot)
		out = strings.ReplaceAll(out, "{"+abbr+"}", styleFacilityToken(token, status))
	}
	return out
}

func renderOutpostCompact(s game.State) string {
	lines := make([]string, 0, len(schematicCompactOrder))
	for _, slot := range schematicCompactOrder {
		token, status := facilityTokenForState(s, slot)
		lines = append(lines, styleFacilityToken(token, status))
	}
	return strings.Join(lines, "\n")
}

func facilityTokenForState(s game.State, slot schematicSlot) (string, facilityDisplayStatus) {
	if slot.id == "signal_beacon" {
		spec := facilitySpec{
			beaconParts: s.BeaconParts,
			maxBeacon:   s.MaxBeaconParts,
			isBeacon:    true,
		}
		if s.MaxBeaconParts > 0 && s.BeaconParts >= s.MaxBeaconParts {
			return formatFacilityToken(slot.abbr, spec), facilityBeaconComplete
		}
		if s.BeaconParts > 0 {
			return formatFacilityToken(slot.abbr, spec), facilityBeaconProgress
		}
		return formatFacilityToken(slot.abbr, spec), facilityUnbuilt
	}

	level := s.BuildingLevel(slot.id)
	def, known := s.FindBuilding(slot.id)
	if !known || level <= 0 {
		return formatFacilityToken(slot.abbr, facilitySpec{maxLevel: def.MaxLevel}), facilityUnbuilt
	}

	b := s.Buildings[slot.id]
	spec := facilitySpec{
		level:    level,
		maxLevel: def.MaxLevel,
		damaged:  b.Damaged,
		offline:  buildingOfflineForPower(s, def, b),
	}
	raw := formatFacilityToken(slot.abbr, spec)
	switch {
	case b.Damaged:
		return raw, facilityDamaged
	case spec.offline:
		return raw, facilityOffline
	case level >= def.MaxLevel:
		return raw, facilityMax
	default:
		return raw, facilityBuilt
	}
}

func buildingOfflineForPower(s game.State, def game.BuildingDef, b game.Building) bool {
	if b.Level <= 0 || b.Damaged || len(def.DailyEffects) == 0 {
		return false
	}
	return s.Power <= powerOfflineThreshold
}

func formatFacilityToken(abbr string, spec facilitySpec) string {
	if spec.isBeacon {
		if spec.maxBeacon > 0 && spec.beaconParts >= spec.maxBeacon {
			return "[" + abbr + "★]"
		}
		return "[" + abbr + " " + strconv.Itoa(spec.beaconParts) + "/" + strconv.Itoa(spec.maxBeacon) + "]"
	}
	if spec.damaged {
		return "[" + abbr + "!]"
	}
	if spec.offline {
		return "[" + abbr + "×]"
	}
	if spec.level <= 0 {
		return "[" + abbr + "]"
	}
	if spec.maxLevel > 0 && spec.level >= spec.maxLevel {
		return "[" + abbr + "★]"
	}
	if suffix := levelSuperscript(spec.level); suffix != "" {
		return "[" + abbr + suffix + "]"
	}
	return "[" + abbr + "]"
}

func levelSuperscript(level int) string {
	switch level {
	case 2:
		return "²"
	case 3:
		return "³"
	default:
		if level > 1 {
			return strconv.Itoa(level)
		}
		return ""
	}
}

func styleFacilityToken(raw string, status facilityDisplayStatus) string {
	switch status {
	case facilityUnbuilt:
		return mutedStyle.Render(raw)
	case facilityDamaged:
		return warnStyle.Render(raw)
	case facilityOffline:
		return mutedStyle.Render(raw)
	case facilityMax, facilityBeaconComplete:
		return goodStyle.Render(raw)
	case facilityBeaconProgress:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Render(raw)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true).Render(raw)
	}
}
