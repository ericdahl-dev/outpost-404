# Outpost 404

A tiny terminal base builder built with Go and the Charm stack.

**Site:** [ericdahl-dev.github.io/outpost-404](https://ericdahl-dev.github.io/outpost-404/)

You are the systems operator for a remote survival colony. Keep power online, food growing, morale stable, and the logs quiet while completing a Signal Beacon.

## Why this stack

This project uses the [Charm stack](https://charm.land/) for a terminal-native game:

- **Go** — single binary, fast builds
- **Bubble Tea** — app/game loop and input
- **Lip Gloss** — styling and panel layout
- **Bubbles** — `list` for build/upgrade, `viewport` for the scrollable event log
- **JSON data** — buildings and events tuned without UI code

See `docs/context.md` for the full Charm map (planned: Glamour, Huh, Harmonica).

## Install

### Homebrew

```bash
brew install ericdahl-dev/tap/outpost-404
```

### Go

```bash
go install github.com/ericdahl-dev/outpost-404/cmd/outpost@latest
```

## Run it

From the repo (uses `./data` when present, otherwise embedded content):

```bash
go mod tidy
go run ./cmd/outpost
```

## Build it

```bash
go build -o outpost ./cmd/outpost
./outpost -version
```

## CI and releases

Matches [git-green](https://github.com/ericdahl-dev/git-green): GitHub Actions runs **build**, **test** (`-race`), **golangci-lint**, and an **≥ 80%** coverage floor on `internal/game` (`scripts/check-game-coverage.sh`) on every push/PR; pushing a `v*` tag runs **GoReleaser** (multi-platform binaries + `ericdahl-dev/homebrew-tap` formula).

Repo secrets required for releases:

- `GITHUB_TOKEN` (provided by Actions)
- `HOMEBREW_TAP_GITHUB_TOKEN` — PAT with write access to `ericdahl-dev/homebrew-tap`

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Test it

Game rules are covered by unit tests in `test/internal/game` (mirrors `internal/game`; CI requires ≥ 80% statement coverage on that package; see [docs/balance.md](docs/balance.md)):

```bash
go test ./...          # game rules + cmd/outpost CLI contract tests
./scripts/check-game-coverage.sh
```

## Session logs

Each run can write a **JSONL** session log for balance and play analysis. Logs include `session_start`, every player action (`build`, `repair`, `trade`, `beacon`, `next_day`), random `event_id` when one fires, and `game_end` with before/after stat snapshots.

By default logs go to your OS cache dir, e.g. `~/Library/Caches/outpost-404/sessions/` on macOS. The path is printed to stderr when a session starts.

```bash
# default location
go run ./cmd/outpost

# custom file
go run ./cmd/outpost -log ./logs/my-run.jsonl

# disable logging
go run ./cmd/outpost -log off

# or via env (used when -log is not set)
OUTPOST_LOG=./logs/run.jsonl go run ./cmd/outpost
```

Example line (pretty-printed):

```json
{"ts":"2026-05-21T12:00:00Z","session_id":"1716292800000000000","type":"session_start","day":1,"snapshot":{"day":1,"power":65},"detail":{"seed":1716292800000000000}}
```

`session_start` includes a `seed` so runs are reproducible.

## Replay and simulation

**Replay** reapplies a JSONL log headlessly and checks each step against recorded snapshots:

```bash
go run ./cmd/outpost -replay ./logs/my-run.jsonl
```

**Fixed seed** for interactive play (same random events if you repeat the same actions):

```bash
go run ./cmd/outpost -seed 42
```

**Headless simulate** runs a JSON script without the TUI (for balance tuning and regression checks):

```bash
go run ./cmd/outpost -simulate scripts/conservative.json
go run ./cmd/outpost -simulate scripts/conservative.json -seed 99
go run ./cmd/outpost -simulate scripts/conservative.json -seeds 1,42,99,100,101
```

Script formats:

- Array of actions: `[{"type":"build","building_id":"solar_array"},{"type":"next_day"}]`
- Object with optional seed: `{"seed":42,"actions":[...]}`

Action `type` values: `build`, `repair`, `trade`, `beacon`, `next_day`. Use `-seed` to override the script seed; use `-seeds` for a comma-separated sweep (prints one outcome line per seed and a win count on stderr).

Example output:

```text
seed=42 day=5 won=false game_over=true beacon=0/5 power=40 food=0 morale=55 credits=30
```

You can also call `game.Simulate` from Go tests (see `internal/game/replay.go`).

**Balance baseline** — four reference scripts (`conservative.json`, `no_trade_survival.json`, `beacon_rush.json`, `survival_45.json`) run across fixed seeds in `go test`; see [docs/balance.md](docs/balance.md) for seeds, viability rules, and how to update expectations after tuning.

Logs recorded before seeds were added cannot be replayed; record a new session with the current build.

## Controls

| Key | Action |
| --- | --- |
| `b` | Open build/upgrade menu |
| `j` / `k` or arrows | Move in build menu; scroll event log on main screen |
| `enter` | Build or upgrade selected facility |
| `r` | Repair systems |
| `t` | Trade food for credits |
| `s` | Work on Signal Beacon |
| `n` or `space` | Advance to next day |
| `?` | Toggle help |
| `esc` | Return to main screen |
| `q` | Quit |

## Game objective

Win by doing either of these:

- **Survive 45 days** — reach day 46 with power, food, morale, and population still above zero. Collapse is checked before the survival win, so hitting zero on the final day still loses.
- **Signal Beacon** — complete 5 beacon parts before collapse.

Lose if power, food, morale, or population hits zero on any day.

The survival path is tuned separately from beacon rush; see `scripts/survival_45.json` and [docs/balance.md](docs/balance.md).

## Project structure

```text
cmd/outpost/              # executable entrypoint
test/cmd/outpost/         # CLI contract tests (mirrors cmd/outpost)
internal/game/            # game state, rules, actions, daily tick, content loading
test/internal/game/       # game rule tests (mirrors internal/game)
internal/ui/              # Bubble Tea model, Bubbles widgets, Lip Gloss styles
data/buildings.json       # data-driven building definitions (incl. dailyEffects)
data/events.json          # data-driven random events
scripts/                  # balance reference scripts + check-game-coverage.sh
docs/site/                # GitHub Pages landing (deployed on push to main)
CONTEXT.md                # domain glossary (player-facing terms)
docs/context.md           # design context and direction
docs/gameplay-depth-plan.md  # v0.2 gameplay decisions and issue order
docs/balance.md           # baseline seeds, scripts, coverage policy
AGENTS.md                 # coding-agent instructions
```

## Roadmap (v0.2)

Product spec: [GitHub #78 PRD](https://github.com/ericdahl-dev/outpost-404/issues/78).

| Track | Plan | Milestone |
| --- | --- | --- |
| Gameplay depth | [docs/gameplay-depth-plan.md](docs/gameplay-depth-plan.md) | [Gameplay depth v0.2](https://github.com/ericdahl-dev/outpost-404/milestone/2) |
| TUI graphics | [docs/tui-graphics-plan.md](docs/tui-graphics-plan.md) | [TUI graphics pass](https://github.com/ericdahl-dev/outpost-404/milestone/1) |

Priority: weighted events and daily effects first, then ASCII outpost schematic, then save/scenarios. Later: Glamour/Huh (see [docs/context.md](docs/context.md)).
