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
| `scripts/conservative_mid.json` | Same opener through **day 11+** alive on all reference seeds (`TestEarlyBalance_*`) |
| `scripts/no_trade_survival.json` | Mid-game probe **without trade**; must **finish alive** on day 14 on all reference seeds |
| `scripts/beacon_rush.json` | Solar then beacon work; must reach **day 6** with **≥2 beacon parts** |
| `scripts/survival_45.json` | Hydro/solar L3, workshop (no repair buffs); solar daily **10** power; **45 `next_day`** → **day 46 win** on all reference seeds |

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
go run ./cmd/outpost -simulate scripts/survival_45.json -seeds 1,7,42,99,100,101,1779403310247544000
```

## Interpreting results

- **Conservative alive on day 5** — early food/power loop still viable; if this fails, early drain or build costs are too harsh.
- **No-trade reaches day 14 alive** — mid-game without trade; baseline expects survival through the full script after #16 food/trade/production tuning.
- **Conservative mid day 11+** — early hydro/solar path must not collapse on days 3–5 and should reach day 11 on every reference seed.
- **Beacon rush beacon ≥2** — rush path still buys parts before collapse; power/credit gates for `beacon` still reachable.
- **Survival 45 day 46 win** — `scripts/survival_45.json` with exactly 45 `next_day` actions; hydro/solar L3, workshop, repairs; `survival_45` baseline + `TestSurvival45_*`.
- **Survival stat integrity (#31)** — `CheckEnd` treats collapse before the day-46 win. `TestSurvival45_NoVitalHitsZeroDuringScript` walks snapshots; end floors: power ≥ `SurvivalMinEndPower` (4), food ≥ `SurvivalMinEndFood` (10) on all reference seeds.
- **Exact mismatch** — RNG or rules changed; update `baseline.go` expectations only after intentional balance work (#16+).

Win-rate sweeps (`sweep: N/M won`) are for exploration; the automated baseline uses **documented outcomes**, not “must win on all seeds.”

## Trade guard

`Trade()` is blocked when `food <= 30` (`MinFoodToTrade` in `internal/game/trade_balance.go`). Rejections log a clear message and record `ok: false`, `reason: "low_food"` in session JSONL for replay.

## Updating expectations (TDD)

1. **RED** — Change or add a baseline assertion first (or let an existing test fail after a deliberate `data/` edit).
2. **GREEN** — Tune `data/*.json` or `internal/game` until `go test ./internal/game/ -run TestBalanceBaseline` passes.
3. **REFACTOR** — Only while green; update `Expected` in `internal/game/baseline.go` to match intentional new outcomes.

Use `-simulate` sweeps to explore; lock outcomes in tests before merging balance PRs.

Do not loosen checks to greenwash a regression—change data/rules or document a deliberate baseline shift in the PR.

## Coverage floor (`internal/game`)

CI enforces **≥ 80%** statement coverage on `internal/game` only (`cmd/outpost` and `internal/ui` are excluded). The check runs after `go test -race ./...` via `scripts/check-game-coverage.sh`.

```bash
./scripts/check-game-coverage.sh
```

Override the minimum locally with `MIN_GAME_COVERAGE=85 ./scripts/check-game-coverage.sh`. Raise the default in the script and this doc together when coverage improves intentionally (#47).

## Building daily production

`data/buildings.json` may include `dailyEffects` (per-level per day). `advanceDay` applies them **before** power/food upkeep and morale drift (`applyBuildingProduction` in `internal/game/tick.go`). One-time build bonuses stay in `effects`.

## Scenarios — strategic identity

Each non-standard scenario has a situational mechanic that rewards a distinct opener.

| Scenario | Strategic hook | Good opener |
|----------|----------------|-------------|
| **Dust Season** | Extra `dailyPowerDelta: +3` drain from regolith storms; higher event frequency (`eventGateSkipAbove: 52`). Power buffer depletes 3× faster, and damaged buildings halve output. | Build solar array first, repair aggressively rather than saving credits. The `dust_season_power_buffer.json` script demonstrates this. |
| **Silent Colony** | `dailyCreditsIncomeDelta: -8` — supply runs arrive sporadically until the radio is on air. Credit flow is tight until a Radio Tower is built (which also unlocks `trader_arrives` events). | Radio Tower early, before a second hydro or solar upgrade. `silent_colony_radio_first.json` demonstrates the cost/benefit. |
| **Beacon Rush** | `maxBeaconParts: 3` and `dailyCreditsIncomeDelta: +7` — motivated supply runners. You can reach a win-beacon faster without grinding Hydro L3, because the credit income more than offsets skipping one food upgrade. | Skip Hydro/Solar L3 and divert credits to beacon work. The standard `beacon_rush.json` still validates this path. |
| **First Landing** | Onboarding cushion — extra credits and slightly better morale/food. No strategic hook; same viable openers as Standard Run. |

### Scenario simulate scripts

```bash
# Dust Season — solar-first, repair-heavy opener
outpost -simulate scripts/dust_season_power_buffer.json -seeds 1,7,42

# Silent Colony — radio tower early
outpost -simulate scripts/silent_colony_radio_first.json -seeds 1,7,42
```
