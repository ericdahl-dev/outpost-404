# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`docs/context.md`** — game pitch, Charm stack map, tone, MVP scope
- **`docs/balance.md`** — reference scripts, baseline seeds, coverage policy, TDD tuning workflow
- **`docs/adr/`** — ADRs that touch the area you're about to work in (when present)

If `docs/adr/` is empty or missing entries for your area, proceed without blocking. The producer skill (`/grill-with-docs`) can add ADRs when decisions crystallize.

## File structure (this repo)

```
/
├── docs/
│   ├── context.md          # design context
│   ├── balance.md          # balance baselines and CI coverage
│   └── adr/                # architecture decisions (as added)
├── AGENTS.md
├── data/
├── scripts/                # headless balance scripts
├── internal/game/
├── internal/ui/
└── cmd/outpost/
```

Other single-context repos may use root `CONTEXT.md` instead of `docs/context.md`; read whichever exists.

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `docs/context.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't documented yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/grill-with-docs`).

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 (event-sourced orders) — but worth reopening because…_
