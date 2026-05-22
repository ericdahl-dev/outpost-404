# Outpost 404

Terminal survival colony sim: the player keeps resources above collapse while racing to finish a Signal Beacon or outlast the survival window.

## Language

**Run**:
One playthrough from first day through win, loss, or quit.
_Avoid_: Session (reserved for JSONL replay logs), game (too generic).

**Standard run**:
The default **Scenario** profile (`standard`) with baseline win targets and rules.
_Avoid_: Default run, normal mode, vanilla.

**Beacon victory**:
Win when built beacon parts reach the run's maximum (default: 5 parts).
_Avoid_: Tech win, rescue win.

**Survival victory**:
Win when the calendar day exceeds the survival target after completing full days (default: day 46, i.e. 45 `next_day`s completed).
_Avoid_: Time win, endurance win.

**Collapse**:
Loss when power, food, morale, or population hits zero or below.
_Avoid_: Game over (also used as the ended-run flag), death.

**Scenario**:
A named ruleset selected at run start that may override win targets, starting state, event pools, and modifiers. Does not change **Standard run** rules unless the player picks that scenario.
_Avoid_: Campaign, level, map.

**First Landing**:
An onboarding **Scenario** with extra starting resources and **Standard run** win targets; not a separate balance baseline.
_Avoid_: Tutorial mode, easy scenario.

**Difficulty**:
A pressure tuning layer (easy/normal/hard) applied on top of the chosen scenario; adjusts costs, decay, or resource pressure—not separate win paths.
_Avoid_: Mode (collides with scenario).

**Scenario profile**:
The JSON-defined knobs for a **Scenario** (win targets, starting resources, modifiers). Loaded at run start into run state.
_Avoid_: Config, preset file.

**Run modifiers**:
Runtime multipliers and flags applied from the **Scenario profile** and **Difficulty** (e.g. solar output scale). Separate from per-building `dailyEffects`.
_Avoid_: Buffs, perks.

**Event eligibility**:
Rules that exclude an event from the random pool for the current day and **Run** (day range, required **Building**, scenario blocks). Implemented in game logic; criteria expressed in event JSON where possible.
_Avoid_: Event filter, gating.

**Building**:
A constructible outpost structure defined in content data and tracked by level in run state (e.g. Solar Array, Hydroponics, Signal Beacon).
_Avoid_: Facility, structure, installation.

**Event day**:
A calendar day on which the random-event gate succeeds and an event may fire; otherwise the day has no random event.
_Avoid_: Event turn, trigger day.

**Event weight**:
A positive integer on an event definition; higher values increase share of the weighted pick among eligible events on an **Event day**. Does not control how often **Event day** occurs.
_Avoid_: Probability, rarity score.

**Event gate**:
The first random roll that decides whether an **Event day** happens at all (default ~46% on a d100). Tunable per **Scenario** / **Difficulty** via **Run modifiers**.
_Avoid_: Spawn rate, event chance.

**Damaged**:
A run-state flag on a built **Building** meaning **Daily effects** from that building run at half strength until **Repair** clears it. One-time build **effects** are unchanged.
_Avoid_: Broken, disabled.

**Repair**:
Player action that clears **Damaged** on one targeted **Building** for a credit cost from a Go formula (`base + per-level step`). Not a generic resource grant.
_Avoid_: Patch, emergency fix, mend.

**Offline**:
(Future.) A **Building** that produces nothing because the colony lacks power headroom—not the same as **Damaged**.
_Avoid_: Down, unpowered (use only when **Offline** is implemented).

**Warning**:
A derived lose-state alert from current resources (not projected end-of-day), with id, severity, and copy. Evaluated in game logic from Go thresholds; not stored as persistent state.
_Avoid_: Alert, notification, danger.

**Warning escalation**:
When a **Warning** newly appears or rises in severity, the game edge-triggers one event-log line; no repeat while unchanged.
_Avoid_: Reminder, ping.

**Daily effects**:
Per-level per-day resource deltas on a **Building** definition (`dailyEffects` in content data), applied during day advance before upkeep.
_Avoid_: Production, passive income.

**Building modifier**:
An explicit rule on a **Building** definition beyond **Daily effects** (e.g. build cost multiplier), interpreted by targeted game logic—not crammed into `dailyEffects`.
_Avoid_: Perk, trait, buff.

**Session log**:
Append-only JSONL record of a play **Run** for analysis and **Replay**; separate from **Autosave**.
_Avoid_: Save file, transcript (too generic).

**Autosave**:
A versioned JSON document written after each completed day (and on quit) so **Continue** can restore the **Run** without parsing the **Session log**.
_Avoid_: Checkpoint, quicksave.

**RNG step count**:
Number of random draws consumed in the **Run** so far; stored in **Autosave** and applied on load so mid-run random events stay deterministic with the same **seed**.
_Avoid_: rng state, random offset.

**Colony log**:
The ordered `[]string` narrative shown to the player; lines may include a type prefix (`!`, `+`, `$`, `*`) set by game logic when recorded.
_Avoid_: Event log (collides with random **Event day**), history.

**Log kind**:
Category for a **Colony log** line (danger, gain, trade, milestone, etc.) that determines its prefix before display or export.
_Avoid_: Log level, severity.

## Relationships

- A **Run** belongs to exactly one **Scenario** profile and one **Difficulty**; **Standard run** is the `standard` profile unless the player picks **First Landing**, Dust Season, Silent Colony, or Beacon Rush.
- **First Landing** only bumps starting resources; win paths and **Event gate** match **Standard run**.
- **Beacon victory** and **Survival victory** are alternate win paths for the same **Run**; either can end the run successfully.
- **Collapse** ends the run unsuccessfully; it is mutually exclusive with both win paths.
- A **Scenario profile** may set **Run modifiers** and override win targets; **Event eligibility** reads scenario and building state.
- **Difficulty** stacks on the same **Run** after the **Scenario profile** is applied.
- On an **Event day**, one eligible event is chosen by **Event weight**; **Event eligibility** uses `minDay`, `maxDay`, buildings, and scenario rules.
- **Event gate** and **Event weight** are independent levers (frequency vs which event).
- **Repair** clears **Damaged** on one **Building**; events or rules may set **Damaged**.
- **Offline** is deferred; only **Damaged** ships in the first building-condition slice.
- **Warning**s are recomputed from current resources; **Warning escalation** controls log lines (TUI badges consume the same list later).
- **Daily effects** apply resource deltas; **Building modifier**s change other rules (build cost, etc.).
- **Event eligibility** for Radio Tower–gated events uses building level, not a **Building modifier** hack.
- **Autosave** holds colony state, **Colony log**, **seed**, **RNG step count**, scenario/difficulty, and **Damaged** flags; **Session log** remains optional for replay/analysis.
- **Colony log** lines carry **Log kind** prefixes from game logic (not applied later in the TUI).

## Example dialogue

> **Dev:** "Does Beacon Rush change the default survival target?"
> **Domain expert:** "No. **Standard run** still uses day 46 **Survival victory**. **Beacon Rush** is a **Scenario** whose **Scenario profile** sets `maxBeaconParts` to 3—only for that run."
>
> **Dev:** "Where does Silent Colony's radio-tower gate live?"
> **Domain expert:** "In **Event eligibility** on each event in JSON—`requiresBuilding`, not a Go-only switch. Game code evaluates the rule."
>
> **Dev:** "Should we rename Hydroponics to a facility in the glossary?"
> **Domain expert:** "No—**Building** is the term in data, code, and docs. Say 'upgrade the Hydroponics building,' not facility."
>
> **Dev:** "Can I still press R for +12 power without a damaged building?"
> **Domain expert:** "No—**Repair** only fixes **Damaged** buildings you select. Workshop **dailyEffects** are how economy compensates, not a fake repair button."
>
> **Dev:** "Should FOOD CRITICAL repeat every morning?"
> **Domain expert:** "Only on **Warning escalation**—first day it trips, not every day you're still starving."
>
> **Dev:** "Can Workshop discount live in dailyEffects as negative credits?"
> **Domain expert:** "No—that's a **Building modifier** on build cost. **Daily effects** are only per-day resource ticks."
>
> **Dev:** "Is Continue the same file as session JSONL?"
> **Domain expert:** "No—**Autosave** is one JSON file for resume. **Session log** is optional JSONL for replay tooling."
>
> **Dev:** "What if the player skips the scenario picker?"
> **Domain expert:** "They get **Standard run**—scenario id `standard` in the log, not an empty string."
>
> **Dev:** "Who adds the `!` on danger lines?"
> **Domain expert:** "Game logic when writing the **Colony log**—the TUI doesn't reinterpret **Log kind**."
>
> **Dev:** "Is First Landing the same as standard for balance tests?"
> **Domain expert:** "No—it's an onboarding **Scenario** with a cushioned start. **Standard run** stays the CI baseline."

## Flagged ambiguities

- (none)
