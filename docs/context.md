# Outpost 404 Context

**Domain glossary (player-facing terms):** [CONTEXT.md](../CONTEXT.md). This file is product/stack direction; implementation plans live in [gameplay-depth-plan.md](gameplay-depth-plan.md) and [tui-graphics-plan.md](tui-graphics-plan.md).

## Origin

Public landing page: [ericdahl-dev.github.io/outpost-404](https://ericdahl-dev.github.io/outpost-404/) (static site in `docs/site/`, deployed via GitHub Pages).

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
- buildings/upgrades with JSON `dailyEffects` (per-level daily production before upkeep)
- daily progression, upkeep, morale drift, colonist growth
- random events (data-driven; replay pins `event_id` in session logs)
- signal beacon progress
- win/loss: survive **45 days** (win on day 46) or complete the beacon
- help screen
- Bubbles `list` build menu and `viewport` event log
- JSON content for facilities and events
- headless `-simulate`, `-replay`, and `-seed`; JSONL session logs
- balance baselines and CI coverage floor on `internal/game` (see [balance.md](balance.md))

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
- more buildings with `dailyEffects`; richer event weights and chains
- difficulty presets
- achievements
- save files
- Steam-like terminal polish eventually, but keep scope small now
