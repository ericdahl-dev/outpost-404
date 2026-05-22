# TUI graphics improvement plan

**PRD:** [#78 PRD: Gameplay depth and TUI graphics (v0.2)](https://github.com/ericdahl-dev/outpost-404/issues/78)  
**Epic:** [#74 Epic: TUI graphics improvement plan](https://github.com/ericdahl-dev/outpost-404/issues/74)  
**Milestone:** [TUI graphics pass](https://github.com/ericdahl-dev/outpost-404/milestone/1)

Gameplay work: [gameplay-depth-plan.md](gameplay-depth-plan.md) · [milestone: Gameplay depth v0.2](https://github.com/ericdahl-dev/outpost-404/milestone/2)

## Goal

Make Outpost 404 feel less like a stats dashboard and more like a living terminal game while keeping the Charm TUI style.

## Design principle

NASA mission console + survival colony + broken terminal logs.

Not cute. Not overly busy. Not a spreadsheet. Not a roguelike map yet.

## Implementation order

| Order | Issue | Title |
| --- | --- | --- |
| 1 | [#60](https://github.com/ericdahl-dev/outpost-404/issues/60) | ASCII outpost schematic panel |
| 2 | [#69](https://github.com/ericdahl-dev/outpost-404/issues/69) | Improve resource and warning display |
| 3 | [#72](https://github.com/ericdahl-dev/outpost-404/issues/72) | Event log prefixes and formatting |
| 4 | [#66](https://github.com/ericdahl-dev/outpost-404/issues/66) | Stronger win/loss presentation |
| 5 | [#71](https://github.com/ericdahl-dev/outpost-404/issues/71) | Title screen |
| 6 | [#70](https://github.com/ericdahl-dev/outpost-404/issues/70) | TUI responsive layout (outpost-centered) |
| 7 | [#73](https://github.com/ericdahl-dev/outpost-404/issues/73) | Milestone event ASCII vignettes |
| 8 | [#75](https://github.com/ericdahl-dev/outpost-404/issues/75) | Split UI rendering into focused modules |
| Later | [#76](https://github.com/ericdahl-dev/outpost-404/issues/76) | Optional TUI motion (Harmonica) |
| Later | [#77](https://github.com/ericdahl-dev/outpost-404/issues/77) | Terminal visual themes |

## Facility abbreviations (schematic)

| Abbr | Building |
| --- | --- |
| SA | Solar Array |
| HY | Hydroponics |
| HB | Habitat |
| WS | Workshop |
| MB | Med Bay |
| RT | Radio Tower |
| BE | Signal Beacon |

## Visual status examples

- `[SA²]` built level 2
- `[SA★]` max level
- `[SA!]` damaged
- `[SA×]` offline (low power)
- `[BE 3/5]` beacon progress

## Terminal targets

- Playable: 80×24
- Improved: 100×30
- Ideal: 120×35+
