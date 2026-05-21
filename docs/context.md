# Outpost 404 Context

## Origin

The project started as a discussion about what kind of game would work well as a TUI. A terminal city/base builder stood out because the terminal naturally supports dashboards, panels, logs, resource bars, and keyboard-driven decision making.

The chosen direction is to prioritize making a **great terminal game** over choosing the stack that collaborators might find most familiar.

## Stack decision

Use **Go + [Charm stack](https://charm.land/)**. Prefer Charm libraries over hand-rolled TUI when a component fits.

| Library | Role in Outpost 404 | Status |
| --- | --- | --- |
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | App loop, screens, keyboard input | In use |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Styles, boxes, horizontal layout | In use |
| [Bubbles](https://github.com/charmbracelet/bubbles) | `list` (build menu), `viewport` (event log) | In use |
| [Glamour](https://github.com/charmbracelet/glamour) | Markdown help / event copy from JSON | Planned |
| [Huh](https://github.com/charmbracelet/huh) | Difficulty, scenarios, new-game setup | Planned |
| [Harmonica](https://github.com/charmbracelet/harmonica) | Short UI motion (day tick, beacon pulse) | Planned |

**Not targeted for this repo:** Wish (SSH apps), gum (shell scripts), Log (dev logging only).

**Boundaries:**

- `internal/game` has zero Charm imports.
- `internal/ui` owns all Bubble Tea / Bubbles / Lip Gloss code.
- Add a Bubbles widget when it replaces non-trivial custom UI; do not add deps for show.

TypeScript/Ink was considered because collaborators may prefer TypeScript, but Go/Bubble Tea is the better fit for a terminal-native game with single-binary distribution and a clean update loop.

## Game pitch

**Outpost 404** is a terminal base builder where the player operates a remote survival colony.

The player must keep systems online, manage scarce resources, handle random events, and complete a Signal Beacon before the outpost collapses.

## Current MVP

The scaffold includes:

- resource stats: power, food, morale, credits, population
- buildings/upgrades
- daily progression
- random events
- signal beacon progress
- win/loss conditions
- help screen
- Bubbles `list` build menu and `viewport` event log
- JSON content for facilities and events

## Desired tone

Serious survival sim with subtle terminal/infrastructure humor.

Examples:

- Coffee outages affecting morale
- quiet shifts improving morale
- radio tower and telemetry language
- event log as the narrative driver

## Future ideas

- ASCII base map
- named colonists
- production/upkeep model per building
- event weights and event chains
- difficulty presets
- achievements
- save files
- Steam-like terminal polish eventually, but keep scope small now
