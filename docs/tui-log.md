# Colony log prefixes

Game logic sets a single-character prefix via `AddLogKind` in `internal/game/colony_log.go`. The UI colors each line in `internal/ui/logformat.go` for the Bubbles viewport.

| Prefix | Kind | Meaning | Examples |
| --- | --- | --- | --- |
| `!` | danger | Warnings, damage | `! Food critically low…` |
| `+` | gain | Builds, repairs, arrivals | `+ Built Solar Array level 1.` |
| `$` | trade | Trade action | `$ Traded surplus rations…` |
| `*` | milestone | Beacon progress, win/loss | `* Signal Beacon part completed: 2/5.` |
| `>` | event | Random day events | `> Solar Storm: … Power -9.` |
| `·` | system | Refusals, errors, no-ops | `· Not enough credits for Hydroponics.` |
| _(none)_ | plain | Intro / scenario copy | `Survive 45 days or complete 5 beacon parts.` |

Prefixes are stored in save/autosave strings so replay and resume stay consistent.
