# AGENTS.md

## Project goal

Build **Outpost 404**, a Go/Bubble Tea terminal base builder. The priority is making a polished terminal-native game, not a generic dashboard.

## Development principles

- Keep game logic in `internal/game` and UI concerns in `internal/ui`.
- Do not put balance constants directly into view code.
- Prefer small, testable functions for state transitions.
- Keep buildings and events data-driven where possible.
- Avoid adding persistence, networking, or external APIs until the core loop is fun.
- Preserve simple keyboard controls.
- Favor readable code over clever abstractions.

## Charm stack expectations

Maximize [Charm](https://charm.land/) usage in `internal/ui`; keep `internal/game` Charm-free.

- **Bubble Tea** — state update flow, `tea.WindowSizeMsg` for layout
- **Lip Gloss** — layout and styles (no balance constants in styles)
- **Bubbles** — `list` (build menu), `viewport` (event log); add more widgets when they replace custom UI
- **Planned** — Glamour (markdown help/events), Huh (settings), Harmonica (motion)
- Keep the UI responsive at common terminal widths, especially 100+ columns

## Game design expectations

This should feel like a game, not a spreadsheet. Add pressure, tradeoffs, and consequences.

Good mechanics:

- scarce resources
- clear risk/reward choices
- random events with narrative flavor
- win/loss pressure
- visible progress toward Signal Beacon completion

Avoid:

- too many resources too early
- hidden rules that feel unfair
- tedious menu depth
- unbounded scope creep

## Testing guidance

Prioritize tests for `internal/game`:

- building cost and level behavior
- resource clamping
- daily resource consumption
- win/loss conditions
- beacon progress

UI tests are optional until the game loop stabilizes.

## Suggested next task

Balance the 30-day win path; use `game.Simulate` and `-replay` for regression checks.
