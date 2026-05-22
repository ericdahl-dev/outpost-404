# Gameplay depth plan (v0.2)

**PRD:** [#78 PRD: Gameplay depth and TUI graphics (v0.2)](https://github.com/ericdahl-dev/outpost-404/issues/78)  
**Epic:** [#79 Epic: Gameplay depth v0.2](https://github.com/ericdahl-dev/outpost-404/issues/79)  
**Milestone:** [Gameplay depth v0.2](https://github.com/ericdahl-dev/outpost-404/milestone/2)

## Goal

Deepen pacing, tradeoffs, and replay fidelity before new infrastructure (maps, networking, heavy meta).

**Domain glossary:** [CONTEXT.md](../CONTEXT.md) (terms and relationships; grill session 2026-05-21).

## Locked decisions (gameplay)

| Topic | Decision |
| --- | --- |
| **Standard run** | Win/loss rules frozen (`standard` scenario); other scenarios override only when selected. CI baselines stay on **Standard run**. |
| **Scenarios** | Hybrid: `data/scenarios.json` profiles + Go eligibility. Ship `standard`, **First Landing** (+starting resources), **Dust Season**, **Silent Colony**, **Beacon Rush**. **Difficulty** easy/normal/hard in same slice. Silent Colony gates via event `requiresBuilding`, not Go switches. |
| **Events (#58)** | Two-stage: tunable **event gate**, then weighted pick among eligible events. Add `weight`, `maxDay`; gate tunable per scenario/difficulty. |
| **Daily effects (#59)** | `dailyEffects` = per-day resource deltas only. **Building modifiers** (e.g. workshop build discount) are separate fields. Extend existing five buildings (E2); no Med Bay in this slice. |
| **Warnings (#64)** | Go thresholds; current resources (not projected). Edge-trigger **Colony log** lines; TUI badges in #69. |
| **Damage (#63)** | **Damaged** → half daily effects; one-time build effects unchanged. **Repair** targets one `building_id`; formula cost in Go. **Offline** deferred. Migrate `scripts/*.json` off generic repair buff. |
| **Colony log (#65)** | Keep `State.Log` as `[]string`; game sets **log kind** prefixes (`!`, `+`, `$`, `*`). |
| **Save (#61)** | `autosave.json` separate from optional JSONL **session log**. Write after each `next_day` and on quit. Rich document: state, log, seed, scenario, difficulty, damaged map, **RNG step count**. |
| **Validation (#67–68)** | After feature slices; fail fast on bad content. |

## Implementation order

| Order | Issue | Title |
| --- | --- | --- |
| 1 | [#58](https://github.com/ericdahl-dev/outpost-404/issues/58) | Weighted random events |
| 2 | [#59](https://github.com/ericdahl-dev/outpost-404/issues/59) | Expand daily building effects |
| 3 | [#64](https://github.com/ericdahl-dev/outpost-404/issues/64) | Meaningful lose-state warnings |
| 4 | [#63](https://github.com/ericdahl-dev/outpost-404/issues/63) | Building damage and repair |
| 5 | [#65](https://github.com/ericdahl-dev/outpost-404/issues/65) | Event log storytelling |
| 6 | [#62](https://github.com/ericdahl-dev/outpost-404/issues/62) | Scenarios and difficulty settings |
| 7 | [#61](https://github.com/ericdahl-dev/outpost-404/issues/61) | Save/load support |
| 8 | [#67](https://github.com/ericdahl-dev/outpost-404/issues/67) | JSON content validation at startup |
| 9 | [#68](https://github.com/ericdahl-dev/outpost-404/issues/68) | Contributor content documentation |

## Coordinate with TUI

- Warning **UI:** [#69](https://github.com/ericdahl-dev/outpost-404/issues/69) (with #64 logic)
- Log **formatting:** [#72](https://github.com/ericdahl-dev/outpost-404/issues/72) (with #65 copy)
- **End/summary screens:** [#66](https://github.com/ericdahl-dev/outpost-404/issues/66)
- **Schematic damage styling:** [#60](https://github.com/ericdahl-dev/outpost-404/issues/60) after #63

See [tui-graphics-plan.md](tui-graphics-plan.md) and [milestone: TUI graphics pass](https://github.com/ericdahl-dev/outpost-404/milestone/1).
