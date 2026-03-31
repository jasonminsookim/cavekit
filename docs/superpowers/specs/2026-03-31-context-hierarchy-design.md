# Context Hierarchy System — Design Spec

## Overview

A context hierarchy system for Blueprint that organizes project knowledge as a directed acyclic graph (DAG) of documents. Agents enter at the root and traverse only the subgraph relevant to their current task. The hierarchy operates at two levels: Blueprint adopts it internally, and Blueprint scaffolds it for user projects.

The core principle: **agents should only read what they need.** Documents are organized so that index files act as DAG hub nodes — an agent reads the index, identifies relevant edges, and follows only those to leaf documents. No agent ever loads the full tree.

---

## Context Directory Structure

### The 4-Tier Information Flow

```
refs/ (what IS)  -->  blueprints/ (what MUST BE)  -->  plans/ (HOW)  -->  impl/ (what WAS DONE)
     Tier 1                  Tier 2                     Tier 3              Tier 4
```

Each tier consumes the previous tier's output. Cross-references between tiers create the DAG edges that agents traverse.

### Directory Layout

```
context/
├── CLAUDE.md                              # Root entry node: describes all 4 tiers
├── refs/                                  # Tier 1: Source material (read-only input)
│   ├── CLAUDE.md                          # "Source of truth. Organized by source. Read-only."
│   └── {source}/                          # Subdirs per source (e.g., verse-lang-spec/, prd/)
│       └── ...
├── blueprints/                            # Tier 2: WHAT to build
│   ├── CLAUDE.md                          # "Start at blueprint-overview.md. R-numbered reqs."
│   ├── blueprint-overview.md              # Index node (DAG hub)
│   ├── blueprint-{domain}.md              # Leaf — simple domain (single file)
│   └── {domain}/                          # Complex domain gets a subdirectory
│       ├── blueprint-{domain}.md          # Domain index (becomes hub node)
│       └── blueprint-{domain}-{sub}.md    # Sub-domain leaves
├── plans/                                 # Tier 3: HOW to build (task graphs)
│   ├── CLAUDE.md                          # "Start at plan-overview.md. Task dependency tiers."
│   ├── plan-overview.md                   # Index node
│   ├── build-site.md                      # Primary build site
│   ├── build-site-{feature}.md            # Feature-specific build sites
│   └── {domain}/                          # Complex plans get subdirectories
│       └── plan-{domain}-{area}.md
├── impl/                                  # Tier 4: What WAS DONE
│   ├── CLAUDE.md                          # "Start at impl-overview.md. Update after every session."
│   ├── impl-overview.md                   # Index node
│   ├── impl-{domain}.md                   # Per-domain tracking
│   ├── impl-review-findings.md            # Codex review findings ledger
│   ├── dead-ends.md                       # Failed approaches (shared across domains)
│   └── archive/                           # Compacted/archived tracking
```

### Nesting Rule

A domain stays flat (single file) by default. When a blueprint, plan, or impl file covers multiple independent concerns that could be understood separately, it becomes an index file pointing to a subdirectory. The original file stays in place as the index — no reference breakage.

**Trigger:** Cohesion, not line count. If a file has sections that an agent working on one section would never need to read the others, decompose it. The agent that creates or modifies the file (during `/bp:draft` or `/bp:build`) is responsible for judging when to decompose.

**Example:** `blueprint-type-system.md` covers effects lattice, tagged values, and inference rules. These are independently understandable. Decompose:

```
blueprints/
├── blueprint-type-system.md                        # Now an index
└── type-system/
    ├── blueprint-type-system-effects.md
    ├── blueprint-type-system-tagged.md
    └── blueprint-type-system-inference.md
```

---

## CLAUDE.md Hierarchy

### Scope: Full Repository

`CLAUDE.md` files extend beyond `context/` into the source code tree. They form the connective tissue between code and the context DAG.

```
project/
├── CLAUDE.md                          # Project root: build/test commands,
│                                      #   "context/ has the full hierarchy"
├── context/
│   ├── CLAUDE.md                      # Root context node: 4 tiers described
│   ├── refs/CLAUDE.md                 # Tier 1 conventions
│   ├── blueprints/CLAUDE.md           # Tier 2 conventions
│   ├── plans/CLAUDE.md                # Tier 3 conventions
│   └── impl/CLAUDE.md                 # Tier 4 conventions
│
├── src/
│   ├── CLAUDE.md                      # Source code conventions
│   ├── auth/
│   │   ├── CLAUDE.md                  # "implements blueprint-auth.md R1-R3"
│   │   └── ...
│   ├── api/
│   │   ├── CLAUDE.md                  # "implements blueprint-api.md R1-R5"
│   │   └── ...
│   └── parser/
│       ├── CLAUDE.md                  # "implements blueprint-grammar.md R1-R4,
│       │                              #   see plans/build-site.md T-012 through T-018"
│       └── ...
│
├── tests/
│   ├── CLAUDE.md                      # Test conventions, how to run
│   └── ...
│
└── scripts/
    ├── CLAUDE.md                      # Utility script conventions
    └── ...
```

### Loading Behavior

When an agent works in `src/auth/`, it loads hierarchically:
1. `project/CLAUDE.md` — project-level conventions
2. `project/src/CLAUDE.md` — source code conventions
3. `project/src/auth/CLAUDE.md` — **"implements blueprint-auth.md R1-R3"**

The third file is the bridge to the context DAG. The agent knows which blueprint to load if it needs requirements context — without loading the entire `context/blueprints/` directory.

### CLAUDE.md Design Principles

- **Minimal** — 3-10 lines for source-tree files. Never duplicate blueprint content.
- **Connective** — each one names the blueprint requirements and plan tasks it relates to.
- **Contextual** — includes module-specific conventions (error handling patterns, test fixture locations, naming rules).
- **Never auto-generated content that misleads** — if `/bp:build` generates a CLAUDE.md, it only includes mappings it is certain about (tasks it completed, files it created).

### CLAUDE.md Content Templates

**`context/CLAUDE.md`:**
```markdown
# Context Hierarchy

This project uses Blueprint's context hierarchy.

## Tiers
- refs/ — Source material (Tier 1: what IS). Read-only.
- blueprints/ — Requirements (Tier 2: what MUST BE). Start at blueprint-overview.md.
- plans/ — Task graphs (Tier 3: HOW). Start at plan-overview.md.
- impl/ — Progress tracking (Tier 4: what WAS DONE). Start at impl-overview.md.

## Navigation
Start at the overview file in whichever tier is relevant to your task.
Only load domain-specific files when the overview points you there.
```

**`context/blueprints/CLAUDE.md`:**
```markdown
# Blueprints

Blueprints define WHAT needs implementing. They are implementation-agnostic.

## Conventions
- Start with blueprint-overview.md for the domain index
- R-numbered requirements (R1, R2, R3...)
- Every requirement has testable acceptance criteria
- Never prescribe HOW — that belongs in plans/
- Cross-reference related domains
- Decompose into subdirectories when a domain covers multiple independent concerns
```

**`context/plans/CLAUDE.md`:**
```markdown
# Plans

Plans define HOW to implement blueprints. They contain task dependency graphs.

## Conventions
- Start with plan-overview.md for the build site index
- Build sites use T-numbered tasks organized into dependency tiers
- Each task references blueprint requirements by ID
- build-site.md is the primary build site
- build-site-{feature}.md for feature-specific sites
```

**`context/impl/CLAUDE.md`:**
```markdown
# Implementation Tracking

Impls record what was built, what is pending, what failed.

## Conventions
- Start with impl-overview.md for current status across all domains
- impl-{domain}.md for per-domain tracking
- dead-ends.md for failed approaches (critical — prevents retrying failures)
- archive/ for compacted history
- Update after every implementation session
```

**`context/refs/CLAUDE.md`:**
```markdown
# Reference Materials

Source of truth that blueprints are derived from. Read-only.

## Conventions
- Organized by source in subdirectories (e.g., prd/, old-code-docs/, api-spec/)
- Agents read but never modify these files
- Blueprints reference specific files and sections via file:line
```

**`src/{module}/CLAUDE.md` (example):**
```markdown
# Auth Module

Implements:
- blueprint-auth.md R1 (User Registration)
- blueprint-auth.md R2 (Login/Sessions)
- blueprint-auth.md R3 (OAuth)

Build tasks: T-004, T-005, T-006 (build-site.md)
```

---

## Progressive Disclosure: The DAG Traversal

### How Agents Navigate

Agents follow a discovery pattern that mirrors a DAG traversal:

1. **Enter at root** — read `context/CLAUDE.md` to understand the 4 tiers
2. **Select tier** — based on current task, navigate to the relevant tier's `CLAUDE.md`
3. **Read index** — the tier's overview file is the DAG hub, listing all domains with one-line summaries
4. **Follow edges** — read only the domain files relevant to the current task
5. **Cross-reference** — if a domain references another domain, follow that edge only if needed
6. **Nest deeper** — if a domain has subdirectories, its root file is the sub-index; spider from there

### Index File Format

Every overview file follows the same format — a table of domains with summary and status:

```markdown
# Blueprint Overview

| Domain | File | Summary | Status |
|--------|------|---------|--------|
| Authentication | blueprint-auth.md | Registration, login, sessions, OAuth | DRAFT |
| Data Models | blueprint-data-models.md | Core entities, relationships, validation | DRAFT |
| Type System | blueprint-type-system.md | Effects lattice, tagged values (see type-system/) | DRAFT |
```

An agent reads this table, identifies "I need Authentication," and loads only `blueprint-auth.md`. It never touches `blueprint-data-models.md` or the `type-system/` subdirectory.

### Cross-Reference Edges

Cross-references between documents are the lateral edges in the DAG. They follow a standard format:

```markdown
**Dependencies:** blueprint-auth.md R2 (session tokens required for API auth)
**See also:** blueprint-api.md R4 (rate limiting uses auth identity)
```

An agent follows these edges only when the cross-referenced content is needed for the current task.

---

## `/bp:init` Command

### Purpose

Project bootstrapping. Creates the full context hierarchy with all `CLAUDE.md` files. Run once at the start of a project.

### Behavior

1. **Scan existing project structure** — detect `src/`, `tests/`, `scripts/`, and any other top-level directories
2. **Create context directories** (if missing):
   - `context/refs/`
   - `context/blueprints/`
   - `context/plans/`
   - `context/impl/`
   - `context/impl/archive/`
3. **Create `CLAUDE.md` files** (using templates above):
   - `context/CLAUDE.md`
   - `context/refs/CLAUDE.md`
   - `context/blueprints/CLAUDE.md`
   - `context/plans/CLAUDE.md`
   - `context/impl/CLAUDE.md`
   - One per detected source directory (`src/CLAUDE.md`, `tests/CLAUDE.md`, `scripts/CLAUDE.md`, etc.)
4. **Create empty index files**:
   - `context/blueprints/blueprint-overview.md`
   - `context/plans/plan-overview.md`
   - `context/impl/impl-overview.md`
5. **Detect legacy layout** — if `context/sites/` exists, offer migration to `context/plans/`
6. **Commit** the scaffolding

### Properties

- **Idempotent** — creates only what's missing. Safe to re-run.
- **Non-destructive** — never overwrites existing files.
- **No questions** — does not ask what you're building (that's `/bp:draft`).

---

## `/bp:build` CLAUDE.md Updates

### When

After `BUILD COMPLETE` — all tasks done, before the completion report.

### How

1. Read the build site to get task-to-blueprint-requirement mappings
2. Read git diff to identify which source files were created/modified during the build
3. For each source directory that was touched:
   - If no `CLAUDE.md` exists: create one with blueprint/plan references derived from the tasks that touched those files
   - If `CLAUDE.md` exists: append any new blueprint references not already listed
4. Update `context/impl/impl-overview.md` with current domain statuses
5. Update `context/plans/plan-overview.md` with build site completion status

### Constraints

- Only writes mappings it is certain about — tasks it completed and files it created
- Never removes existing content from a `CLAUDE.md`
- Source-tree `CLAUDE.md` files are kept minimal (references only, no duplicated content)

---

## Backpropagation via CLAUDE.md

### How `/bp:revise` Uses the Hierarchy

When a bug is found, the source-tree `CLAUDE.md` files provide the reverse traversal path:

```
Bug in src/auth/login.ts
    |
    v
src/auth/CLAUDE.md says "implements blueprint-auth.md R2"
    |
    v
blueprint-auth.md R2 — check acceptance criteria
    |
    |-- Criteria missing?  --> update blueprint (spec gap)
    |-- Criteria wrong?    --> fix blueprint (spec bug)
    |-- Criteria present but code violates? --> fix code (impl bug)
    |
    v
If blueprint changed --> propagate to plans/ --> flag affected tasks
```

### Forward Propagation

When a blueprint changes via `/bp:revise`:
1. Scan all `src/*/CLAUDE.md` files for references to the changed requirement
2. Flag those modules as potentially affected
3. If new requirements are added, they appear as unimplemented — no source-tree `CLAUDE.md` references them yet

---

## Backward Compatibility

### `sites/` to `plans/` Migration

**Resolution order** for all Blueprint commands that need build sites or plans:
1. Look in `context/plans/`
2. If not found, fall back to `context/sites/`
3. If found in `sites/`, use it — no auto-migration, no breakage

**Migration via `/bp:init`:** When run on an existing project with `context/sites/`, offers:
```
Found context/sites/ (legacy layout). Migrate to context/plans/?
This moves build-site files and updates internal references.
[Y/n]
```

If declined, the fallback is permanent. The system works with either layout.

**File naming:** `build-site-*.md` filenames stay unchanged inside `plans/`. The directory changed, the file convention didn't.

### Projects Without `/bp:init`

Projects that never ran `/bp:init` continue working. The `CLAUDE.md` files and index files are additive — their absence doesn't break anything. Blueprint commands work the same as before, just without progressive disclosure guidance.

---

## Blueprint Internal Adoption

Blueprint itself (the plugin repo) adopts this hierarchy for its own `context/` directory:

| Current | New |
|---------|-----|
| `context/blueprints/` | `context/blueprints/` (unchanged) |
| `context/sites/` | `context/plans/` (renamed) |
| `context/impl/` | `context/impl/` (unchanged) |
| (none) | `context/refs/` (new — for project reference materials) |
| (none) | `context/CLAUDE.md` + per-subdir `CLAUDE.md` files (new) |
| (none) | Index files: `blueprint-overview.md`, `plan-overview.md`, `impl-overview.md` (formalized) |
| `references/` | `references/` (stays — plugin infrastructure, not project context) |

---

## Summary of Decisions

| Decision | Choice |
|----------|--------|
| Blueprint naming | Keep `blueprints/` — the brand |
| References split | `references/` = plugin infra, `context/refs/` = project source material |
| Plans directory | Rename `sites/` to `plans/`, backward compat with fallback |
| Prompts directory | None — slash commands are the pipeline |
| CLAUDE.md scope | Full repo — context/ directories AND source tree |
| Scaffolding command | New `/bp:init` — creates all directories, CLAUDE.md files, index files |
| Build-time updates | `/bp:build` generates/updates source-tree CLAUDE.md at build completion |
| Progressive disclosure | Mandatory index files as DAG hub nodes, agents spider from there |
| Nesting trigger | Cohesion-based — decompose when independent concerns share a file |
| Backpropagation | `/bp:revise` traverses CLAUDE.md edges in reverse to trace bugs to specs |
| Backward compatibility | `sites/` fallback permanent, `/bp:init` optional, no forced migration |
