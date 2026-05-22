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
- Keep code easy to understand and self-documenting: clear names, obvious control flow, and comments only where intent is not evident from the code itself.

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

### Test-driven development (required)

Use **TDD** for game logic and balance changes: one behavior per cycle, not a batch of tests then a batch of code.

**Workflow (vertical slices):**

1. **RED** — Write one test that asserts observable behavior through a public API (`NextDay`, `Build`, `Simulate`, `CheckBaselineOutcome`, etc.). Run it; it must fail for the right reason.
2. **GREEN** — Add the smallest change in `internal/game` (or `data/` when tuning JSON) that makes that test pass.
3. **REFACTOR** — Clean up only while green; run `go test ./internal/game/...` after each step.

**Rules:**

- Do not write all tests upfront, then all implementation (horizontal slices produce brittle tests).
- Tests describe **behavior** (player-visible outcomes, replay fidelity, baseline end states), not private helpers or implementation shape.
- Prefer `game.Simulate` and `TestBalanceBaseline` for multi-step flows; unit-test single transitions when that is clearer.
- Balance tuning (#16+): failing baseline test first, then adjust `data/` or rules, then update `baseline.go` expectations only when the new balance is intentional (see `docs/balance.md`).
- No tests for states the type system or invariants already rule out.

**Before claiming done:** `go test ./...` must pass with evidence from the terminal, not assumed.

## Suggested next task

Extend `dailyEffects` on more buildings or add scenarios; use `game.Simulate` / `SimulateWithSnapshots` and `-replay` for regression checks.

## Agent skills

### Issue tracker

Issues live in GitHub Issues (`ericdahl-dev/outpost-404`), managed via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

Default canonical label vocabulary (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repo — one `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.

## Learned User Preferences

- Never commit or push work directly to `main`; use a feature branch and open a PR.

## Learned Workspace Facts

- CI and releases follow the local `git-green` repo pattern: GitHub Actions (build, `go test -race`, golangci-lint), GoReleaser on `v*` tags, Homebrew formula in `ericdahl-dev/homebrew-tap` (needs `HOMEBREW_TAP_GITHUB_TOKEN`).
- Shipped binaries embed `data/` JSON; a checkout’s `./data` directory overrides embedded content when present.
- Headless balance tooling: `go run ./cmd/outpost -simulate scripts/*.json`, `-replay` on JSONL session logs, plus `game.Simulate` / `TestBalanceBaseline` in tests.
