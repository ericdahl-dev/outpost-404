# Balance baseline

Headless scripts and fixed seeds catch accidental balance regressions before tuning `data/events.json` or building costs.

## Reference seeds

These RNG seeds are checked on every baseline run (see `game.ReferenceSeeds`):

| Seed | Note |
| --- | --- |
| `1`, `7`, `42`, `99`, `100`, `101` | Small spread for event RNG |
| `1779403310247544000` | From a real `session_start` log (JSON float-safe) |

Override a script’s embedded seed with `-seed` or sweep with `-seeds` (comma-separated).

## Reference strategies

| Script | Intent |
| --- | --- |
| `scripts/conservative.json` | Early hydro + solar, one repair; must **finish alive** on day 5 on all reference seeds |
| `scripts/no_trade_survival.json` | Mid-game probe **without trade**; must reach **day 8+** (may collapse after) |
| `scripts/beacon_rush.json` | Solar then beacon work; must reach **day 6** with **≥2 beacon parts** |

Exact end states (day, `game_over`, beacon parts) are locked in `internal/game/baseline_test.go`. Changing JSON balance or events without updating expectations should fail CI.

## Running checks

```bash
go test ./internal/game/ -run TestBalanceBaseline
```

CLI sweeps (after `-simulate` is available):

```bash
go run ./cmd/outpost -simulate scripts/conservative.json -seeds 1,7,42,99,100,101,1779403310247544000
go run ./cmd/outpost -simulate scripts/no_trade_survival.json -seeds 1,7,42,99,100,101,1779403310247544000
go run ./cmd/outpost -simulate scripts/beacon_rush.json -seeds 1,7,42,99,100,101,1779403310247544000
```

## Interpreting results

- **Conservative alive on day 5** — early food/power loop still viable; if this fails, early drain or build costs are too harsh.
- **No-trade reaches day 8+** — mid-game without trade crutch; failures often mean food drain or events are too punishing before day 10.
- **Beacon rush beacon ≥2** — rush path still buys parts before collapse; power/credit gates for `beacon` still reachable.
- **Exact mismatch** — RNG or rules changed; update `baseline.go` expectations only after intentional balance work (#16+).

Win-rate sweeps (`sweep: N/M won`) are for exploration; the automated baseline uses **documented outcomes**, not “must win on all seeds.”

## Updating expectations

1. Run the failing strategy with `-simulate` and note per-seed lines.
2. Adjust `Expected` in `internal/game/baseline.go` for that strategy.
3. Re-run `go test ./internal/game/ -run TestBalanceBaseline`.

Do not loosen checks to greenwash a regression—change data/rules or document a deliberate baseline shift in the PR.
