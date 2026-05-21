# Outpost 404

A tiny terminal base builder built with Go and the Charm stack.

You are the systems operator for a remote survival colony. Keep power online, food growing, morale stable, and the logs quiet while completing a Signal Beacon.

## Why this stack

This project uses the [Charm stack](https://charm.land/) for a terminal-native game:

- **Go** — single binary, fast builds
- **Bubble Tea** — app/game loop and input
- **Lip Gloss** — styling and panel layout
- **Bubbles** — `list` for build/upgrade, `viewport` for the scrollable event log
- **JSON data** — buildings and events tuned without UI code

See `docs/context.md` for the full Charm map (planned: Glamour, Huh, Harmonica).

## Run it

```bash
go mod tidy
go run ./cmd/outpost
```

## Build it

```bash
go build -o outpost ./cmd/outpost
./outpost
```

## Test it

Game rules are covered by unit tests in `internal/game`:

```bash
go test ./...
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

**Headless scripts** use `game.Simulate` in Go (see `internal/game/replay.go`):

```go
final, err := game.Simulate(content, 42, []game.SimAction{
    {Type: "build", BuildingID: "solar_array"},
    {Type: "next_day"},
    {Type: "beacon"},
})
```

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

- survive 30 days
- complete 5 Signal Beacon parts

Lose if power, food, morale, or population hits zero.

## Project structure

```text
cmd/outpost/main.go       # executable entrypoint
internal/game/            # game state, rules, actions, daily tick, content loading
internal/ui/              # Bubble Tea model, Bubbles widgets, Lip Gloss styles
data/buildings.json       # data-driven building definitions
data/events.json          # data-driven random events
docs/context.md           # design context and direction
AGENTS.md                 # coding-agent instructions
```

## Next good milestones

1. Add save/load.
2. Add scenarios and difficulty settings.
3. Make random events weighted instead of uniform.
4. Add facility upkeep and passive per-day production.
5. Add a map panel or ASCII base layout.
6. Expand `internal/game` tests (events, repair/trade, balance).
7. Balance resources so the 30-day win is tense but fair.
