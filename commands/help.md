---
name: blueprint-help
description: Show Blueprint commands and usage
---

# Blueprint

## The Workflow

```
/blueprint:draft      →  write blueprints (the WHAT)
/blueprint:architect   →  generate site (the ORDER)
/blueprint:build       →  ralph loop (the BUILD)
/blueprint:inspect     →  gap analysis + peer review (the CHECK)
```

## Commands

### `/blueprint:draft` — Write Blueprints

```bash
/blueprint:draft                        # interactive — asks what to build
/blueprint:draft context/refs/          # from reference materials (PRDs, docs)
/blueprint:draft --from-code            # from existing codebase (brownfield)
/blueprint:draft --filter v2            # only generate v2 blueprints
```

Decomposes your project into domains, writes `context/blueprints/blueprint-{domain}.md` files with R-numbered requirements and testable acceptance criteria.

### `/blueprint:architect` — Generate Site

```bash
/blueprint:architect                    # generates site from all blueprints
/blueprint:architect --filter v2        # only v2 blueprints
```

Reads blueprints, decomposes requirements into tasks, organizes into dependency tiers. Writes `context/frontiers/build-site.md`. No domain plans — just tasks and dependencies.

### `/blueprint:build` — Run the Loop

```bash
/blueprint:build                       # ralph loop from site
/blueprint:build --filter v2           # scope to v2
/blueprint:build --peer-review         # add Codex (GPT-5.4) review
/blueprint:build --max-iterations 30   # iteration limit
/blueprint:build --peer-review --codex-model gpt-5.4-mini
```

Auto-archives any previous cycle, then starts a Ralph Loop. Each iteration: pick unblocked task → read blueprint → implement → validate → commit.

With `--peer-review`: alternates build and review iterations, calling Codex via MCP.

### `/blueprint:inspect` — Post-Loop Inspection

```bash
/blueprint:inspect                     # inspect everything from last loop
/blueprint:inspect --filter v2         # only v2
```

Runs after build completes. Does two things:
1. **Gap analysis** — compares built code against every blueprint requirement and acceptance criterion
2. **Peer review** — finds bugs, security issues, performance problems, quality gaps

Produces a verdict: APPROVE / REVISE / REJECT with prioritized findings.

### `/blueprint:progress` — Check Progress

```bash
/blueprint:progress                    # show site progress
/blueprint:progress --filter v2
```

Shows tasks done/ready/blocked, progress bar, current tier, and next tasks.

### Maintenance (optional)

| Command | When |
|---------|------|
| `/blueprint:gap-analysis` | After a loop — compare built vs intended |
| `/blueprint:revise` | After manual code fixes — trace back to blueprints |
| `/blueprint:compact-specs` | When impl tracking files exceed ~500 lines |
| `/blueprint:archive-loop` | Manually archive a loop cycle (build does this automatically) |
| `/blueprint:next-session` | Generate a handoff document for next session |

### Legacy (advanced)

These still work but are superseded by the three main commands:

| Command | Replaced by |
|---------|-------------|
| `/blueprint init` | `/blueprint:draft` creates directories automatically |
| `/blueprint spec-from-refs` | `/blueprint:draft context/refs/` |
| `/blueprint spec-from-code` | `/blueprint:draft --from-code` |
| `/blueprint plan-from-specs` | `/blueprint:architect` (generates site directly, no domain plans) |
| `/blueprint implement` | `/blueprint:build` (one task at a time vs full loop) |
| `/blueprint spec-loop` | `/blueprint:build` |
| `/blueprint peer-review-loop` | `/blueprint:build --peer-review` |
| `/blueprint quick` | `/blueprint:draft` + `/blueprint:architect` + `/blueprint:build` |

## Skills (reference docs)

| Skill | Topic |
|-------|-------|
| `blueprint:methodology` | Core DABI lifecycle |
| `blueprint:blueprint-writing` | How to write blueprints with testable criteria |
| `blueprint:peer-review` | Cross-model review patterns |
| `blueprint:peer-review-loop` | Ralph Loop + Codex architecture |
| `blueprint:validation-first` | Every requirement must be auto-testable |
| `blueprint:convergence-monitoring` | Detecting if loop is converging or stuck |
| `blueprint:revision` | Tracing bugs back to blueprints |
| `blueprint:context-architecture` | Organizing context/ for AI agents |
| `blueprint:impl-tracking` | Progress tracking and dead ends |
| `blueprint:brownfield-adoption` | Adopting Blueprint on existing codebases |
| `blueprint:prompt-pipeline` | Designing prompt sequences |
| `blueprint:speculative-pipeline` | Staggered pipeline execution |
| `blueprint:documentation-inversion` | Agent-first documentation |
