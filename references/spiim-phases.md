# SPIIM Phases Reference

Complete reference for the five-phase SDD lifecycle: Spec, Plan, Implement, Iterate, Monitor.

---

## 1. Overview

SPIIM is the five-phase lifecycle of Spec-Driven Development. Each phase has dedicated prompts that drive it, explicit inputs and outputs, and defined roles for both the AI agent and the human engineer.

The core principle that governs all phases:

> **Idea -> Specs -> Code, Never Idea -> Code**

Specs are the first-class citizen. Never go directly from source material to code. Always go through an intermediate specification stage.

| Project Type | Idea (Source Material) | Specs | Code (Output) |
|---|---|---|---|
| **Greenfield** | Language spec, PRD, design docs | Technical specs derived from source material | Working application, tests, tooling |
| **Rewrite** | Old application source code | Implementation-agnostic specs reverse-engineered from old code | New application in target framework |

---

## 2. Phase Table

| Phase | Input | Output | AI Role | Human Role |
|-------|-------|--------|---------|------------|
| **Spec** | Old code, reference docs, research | Implementation-agnostic specs | Analyze, document, organize | Review specs for completeness |
| **Plan** | Specs + framework research | Framework-specific implementation plans | Architect, decompose, sequence | Validate architecture decisions |
| **Implement** | Plans + specs | Working code + tests + tracking docs | Build, test, validate, generate tests on changed files | Monitor progress |
| **Iterate** | Failed validations, gaps, manual fixes | Updated specs/plans via backpropagation | Diagnose, backpropagate, fix | Audit results, steer direction |
| **Monitor** | Running application, git history, worktree activity | Issues, anomalies, progress reports | Observe, scan, report | Review reports, steer direction, trigger backpropagation |

---

## 3. Phase Details

### 3.1 Spec Phase

**Purpose:** Transform source material into implementation-agnostic specifications that define WHAT needs to be built.

**Inputs:**
- Reference materials (PRDs, language specs, old code docs, design documents)
- Feature scope documents (what is in/out of scope)
- Existing codebase (for brownfield/rewrite projects)

**Outputs:**
- Domain-specific spec files (`specs/spec-{domain}.md`)
- Spec overview/index file (`specs/spec-overview.md`)
- Cross-references between related specs

**AI Role:**
- Read and analyze all reference materials
- Decompose into domain-specific specifications
- Write specs with testable acceptance criteria
- Cross-reference specs where domains interact

**Human Role:**
- Review specs for completeness and accuracy
- Verify acceptance criteria are testable
- Ensure scope is correct (not too broad, not too narrow)
- Validate domain decomposition makes sense

**Key Principles:**
- Specs are implementation-agnostic -- they describe WHAT, not HOW
- Every requirement must include testable acceptance criteria
- Specs must be hierarchical -- one index file linking to domain-specific sub-specs
- Specs must be cross-referenced -- related specs link to each other

**Spec Format Template:**
```markdown
# Spec: {Domain Name}

## Scope
{What this spec covers}

## Requirements

### R1: {Requirement Name}
**Description:** {What must be true}
**Acceptance Criteria:**
- [ ] {Testable criterion 1}
- [ ] {Testable criterion 2}
**Dependencies:** {Other specs/requirements this depends on}

### R2: ...

## Out of Scope
{Explicit exclusions}

## Cross-References
- See also: spec-{related-domain}.md
```

**Greenfield Pattern:**
- Reference material -> specs (single prompt, e.g., `001-generate-specs-from-refs.md`)
- Agent reads `context/refs/` and produces `context/specs/`

**Rewrite Pattern:**
- Old code -> reference docs -> specs (multiple prompts)
- `001`: Generate reference materials from old code
- `002`: Generate specs from reference + feature scope
- `003`: Validate specs against codebase

---

### 3.2 Plan Phase

**Purpose:** Transform implementation-agnostic specs into framework-specific implementation plans that define HOW to build.

**Inputs:**
- Specs from the Spec phase
- Framework documentation and research
- Existing implementation tracking (if any)

**Outputs:**
- Domain-specific plan files (`plans/plan-{domain}.md`)
- Feature frontier document (`plans/plan-feature-frontier.md`)
- Known issues backlog (`plans/plan-known-issues.md`)

**AI Role:**
- Read specs and research framework patterns
- Architect the implementation approach
- Decompose into tasks with dependencies
- Sequence implementation order
- Define test strategies per feature

**Human Role:**
- Validate architecture decisions
- Review framework choices
- Approve dependency ordering
- Verify test strategy coverage

**Key Principles:**
- Plans are framework-specific -- they describe HOW to implement
- Plans reference specs for the WHAT
- Plans include feature dependencies (what must be built first)
- Plans include test strategies (how each feature will be validated)
- Plans include acceptance criteria (runnable checks)

**Plan Format Template:**
```markdown
# Plan: {Domain Name}

## Framework
{FRAMEWORK} — {version and key dependencies}

## Implementation Sequence

### Task T-1: {Task Name}
**Spec Reference:** spec-{domain}.md R1
**Dependencies:** None
**Files:** {files to create/modify}
**Approach:** {How to implement}
**Tests:** {Test strategy}
**Acceptance:** {BUILD_COMMAND} passes, tests pass

### Task T-2: {Task Name}
**Spec Reference:** spec-{domain}.md R2
**Dependencies:** T-1
**Files:** {files to create/modify}
...

## Feature Frontier
| Tier | Features | Dependencies |
|------|----------|-------------|
| 1 (Foundation) | Core data models, basic routing | None |
| 2 (Core) | Business logic, API integration | Tier 1 |
| 3 (Advanced) | Performance, polish, edge cases | Tier 2 |

## Known Issues
| Priority | Issue | Workaround |
|----------|-------|------------|
| P0 | {Critical blocker} | {Temporary fix} |
| P1 | {High priority} | {Approach} |
```

**Bidirectional Flow:**
Plans and implementation tracking files update each other. The Plan phase reads `impl/` for feedback from prior implementation passes, and the Implement phase updates plans when it discovers new information. This bidirectional flow is expected and healthy -- it is how the system self-corrects.

---

### 3.3 Implement Phase

**Purpose:** Build working code from plans and specs, with full validation.

**Inputs:**
- Plans from the Plan phase
- Specs from the Spec phase
- Implementation tracking from prior iterations

**Outputs:**
- Source code (`src/`)
- Tests (`tests/`)
- Implementation tracking documents (`impl/impl-{domain}.md`)

**AI Role:**
- Read plans and identify highest-priority unblocked task
- Implement the task
- Generate and run tests on changed files
- Update implementation tracking
- Commit progress frequently

**Human Role:**
- Monitor progress
- Review implementation tracking for anomalies
- Intervene if agent is stuck or going in wrong direction

**Key Principles:**
- Always implement the highest-priority unblocked task
- Run validation gates after each significant change (build -> test -> verify)
- Generate tests for all changed source files
- Update implementation tracking with: files created/modified, issues found, dead ends
- Commit frequently, never push (git as working memory)
- Use sub-agents for discrete subtasks to preserve context window

---

### 3.4 Iterate Phase (Backpropagation)

**Purpose:** When bugs or gaps are found, trace them back to specs/plans and fix at the source.

**Inputs:**
- Failed validations from the Implement phase
- Gaps identified during monitoring
- Manual fixes made by humans

**Outputs:**
- Updated specs with missing requirements/validation
- Updated plans with corrected approaches
- Regression tests
- Systemic prompt improvements

**AI Role:**
- Diagnose the root cause of failures
- Trace issues back to spec or plan gaps
- Update specs (not just code) with missing requirements
- Generate regression tests
- Re-run validation to verify the fix emerges from updated specs alone

**Human Role:**
- Act as **auditor, not implementer**
- Review proposed spec changes
- Make systemic improvements to prompts when issues represent patterns
- Steer direction of backpropagation

**The Backpropagation Process:**
1. **Discover issue** in the running application
2. **Fix it** with the agent (normal debugging session)
3. **Ask the agent:** "What could we have changed about the spec or prompt to have caught this earlier?"
4. **Update the spec** (not the code!) with the missing requirement or validation
5. **Run back-propagation** to trace the manual fix back into context files
6. **Optionally make systemic prompt changes** if the issue represents a pattern
7. **Re-run the iteration loop** -- verify the fix emerges from the updated specs alone
8. **Generate regression tests**

**Key Insight:** If a bug can only be fixed in code, your specs are incomplete. The goal is that specs plus the iteration loop can reproduce any fix autonomously.

---

### 3.5 Monitor Phase

**Purpose:** Continuously observe the running system, agent activity, and development progress.

**Inputs:**
- Running application
- Git history and worktree activity
- Context file changes

**Outputs:**
- Issues and anomalies
- Progress reports
- Convergence metrics
- Backpropagation triggers

**AI Role:**
- Periodically scan worktrees, git history, and context changes
- Report on convergence metrics (test pass rate, change velocity)
- Identify anomalies or stuck agents
- Generate progress summaries

**Human Role:**
- Review monitoring reports
- Steer direction based on progress
- Trigger backpropagation when needed
- Make go/no-go decisions on phase transitions

**Key Metrics:**
- Test pass rate (approaching 100%)
- Change velocity (decreasing = converging)
- Forward progress (% of spec requirements with passing tests)
- Dead end accumulation (increasing = possible spec problem)

---

## 4. The Human is an Auditor, Not an Implementer

This is a critical principle throughout SPIIM:

> The human monitors the process, requests changes as needed, and makes systemic improvements to specs and prompts. The human does NOT write code.

**What the human does:**
- Reviews specs for completeness and accuracy
- Validates architecture decisions in plans
- Monitors implementation progress
- Audits iteration results
- Steers direction when agents are off track
- Makes systemic improvements to prompts
- Triggers backpropagation for discovered issues
- Makes go/no-go decisions at phase gates

**What the human does NOT do:**
- Write code directly
- Fix bugs by editing source files
- Implement features
- Write tests manually

When the human discovers a bug, the correct action is to trace it back to a spec gap and fix the spec, not to fix the code. The iteration loop should then reproduce the fix autonomously from the updated specs.

---

## 5. Phase Gates (Transition Criteria)

Phase gates are mandatory verification checkpoints between phases. No phase transition occurs without passing its gate.

### 5.1 Spec -> Plan Gate

| Criterion | Verification |
|-----------|-------------|
| All domains identified and spec files created | Spec overview lists all domains |
| Every requirement has testable acceptance criteria | Review each `R-` requirement for `[ ]` criteria |
| Cross-references are complete | Each spec links to related specs |
| Scope is defined (in-scope and out-of-scope) | Each spec has explicit exclusions |
| Human has reviewed and approved specs | Human sign-off |

### 5.2 Plan -> Implement Gate

| Criterion | Verification |
|-----------|-------------|
| All spec requirements mapped to plan tasks | Cross-reference check |
| Task dependencies are defined and acyclic | Dependency graph review |
| Test strategies defined for each feature | Each task has test approach |
| Feature frontier established | Tier system documented |
| Framework research complete | Plan references framework docs |
| Human has reviewed architecture decisions | Human sign-off |

### 5.3 Implement -> Iterate Gate

| Criterion | Verification |
|-----------|-------------|
| Build passes | `{BUILD_COMMAND}` exits cleanly |
| Unit tests pass | `{TEST_COMMAND}` exits cleanly |
| Implementation tracking is current | `impl/` files reflect actual state |
| All completed tasks verified | Each DONE task has passing tests |
| No P0 issues outstanding | `plan-known-issues.md` check |

### 5.4 Iterate -> Monitor Gate

| Criterion | Verification |
|-----------|-------------|
| All backpropagation targets addressed | Spec changes committed |
| Regression tests generated and passing | Test suite expanded |
| Iteration loop re-run confirms fixes | Clean iteration pass |
| Implementation tracking updated | Dead ends documented |

### 5.5 Monitor -> Spec Gate (Cycle Back)

| Criterion | Verification |
|-----------|-------------|
| Convergence detected or ceiling diagnosed | Change velocity analysis |
| Gap analysis complete | Built vs intended comparison |
| New requirements or scope changes identified | Spec updates needed |
| Human decision to cycle back | Explicit go/no-go |

---

## 6. The CI Pipeline Analogy

SDD, at its core, is a CI pipeline for AI development:

**Traditional CI/CD:**
```
Code -> Build -> Test -> Deploy
```

**SDD AI Pipeline:**
```
Spec Change
  -> Generate Plans (iteration loop)
    -> Generate Implementation (iteration loop)
      -> Validate (Tests + Review)
        -> Human Audit (Monitor & Steer)
          -> [Gap Found]
            -> Backpropagate
              -> Spec Change (cycle)
```

Each stage feeds the next. Failures at any stage propagate back to the appropriate source (spec, plan, or prompt) rather than being patched at the code level.

---

## 7. When to Use Full SPIIM vs. Lightweight SDD

### Full SPIIM

Use when:
- The codebase is large or complex (50+ source files, multiple domains)
- Requirements will evolve (specs and code need to co-evolve)
- You need multi-agent or multi-prompt pipelines
- Brownfield or production projects where traceability matters
- Multi-team or cross-team environments (shared specs prevent drift)
- Security-sensitive code (validation gates catch vulnerabilities)
- Long-running autonomous work (iteration loops running unattended)

### Lightweight SDD

Use when:
- The task is focused but non-trivial
- You want spec benefits without full pipeline overhead

**Lightweight approach:**
1. Write a focused `context/specs/spec-task.md` capturing requirements
2. Add a `context/plans/plan-task.md` sequencing the implementation
3. Skip full SPIIM; just run the iteration loop against the plan

This is the "SDD floor" -- you get most of the benefit without the overhead of a full multi-phase pipeline.

### Skip SDD

When:
- The task is small and self-contained (~5 files, clear requirements, single session)
- One-off standalone tools with well-defined scope
- Exploratory prototyping where requirements are completely unknown
- Simple bug fixes or small feature additions

**Heuristic:** If you can implement it from memory in a single session, it is probably too small for full SDD.

---

## 8. Prompt Pipeline Patterns

### Greenfield Pattern (3-Prompt)

| Prompt | SPIIM Phase | Input | Output |
|--------|-------------|-------|--------|
| `001-generate-specs-from-refs.md` | **Spec** | `context/refs/` | `context/specs/` |
| `002-generate-plans-from-specs.md` | **Plan** | `context/specs/` | `context/plans/` |
| `003-generate-impl-from-plans.md` | **Implement** | `context/plans/` + `context/specs/` | `src/`, `tests/`, `context/impl/` |

### Rewrite Pattern (6-9 Prompts)

| Prompt | SPIIM Phase | Input | Output |
|--------|-------------|-------|--------|
| `001-generate-refs-from-code.md` | **Spec** (prep) | Old app source | `shared-context/reference/` |
| `002-generate-specs.md` | **Spec** | Feature scope + reference | `shared-context/specs/` |
| `003-validate-specs.md` | **Spec** (verify) | Reference + specs | Validation report |
| `004-create-plans.md` | **Plan** | Specs + framework research | `context/plans/` |
| `005-implement.md` | **Implement** | Plans + specs | `src/` + `tests/` |
| `006-update-specs.md` | **Iterate** | Working prototype | Updated specs |

Prompt 006 back-propagates to 002, creating a feedback loop.

### Shared Principles Across All Pipelines

| Principle | Description |
|-----------|-------------|
| One prompt per SPIIM phase | Clear separation of concerns |
| Explicit input/output directories | Each prompt knows where to read and write |
| Git-based continuity | Agents read git history between iterations |
| Exit criteria with completion signals | `<all-tasks-complete>` signals done |
| Bidirectional spec/plan updates | Plans and impl tracking update each other |
| Test generation on changed files | Validation after every significant change |

---

## 9. Context Directory Structure

Every SDD project follows this standard structure:

```
context/
+-- refs/           # Source of truth (language specs, PRDs, old code docs)
+-- specs/          # Implementation-agnostic specifications
|   +-- CLAUDE.md   # "Specs define WHAT needs implementing"
+-- plans/          # Framework-specific implementation plans
|   +-- CLAUDE.md   # "Plans define HOW to implement something"
+-- impl/           # Living implementation tracking
|   +-- CLAUDE.md   # "Impls record implementation progress"
+-- prompts/        # SPIIM pipeline prompts (001, 002, 003...)
```

Each subdirectory gets a `CLAUDE.md` that describes its conventions. Agents automatically load these when working in that directory. CLAUDE.md is hierarchical -- it loads from the directory AND all parent directories.

---

## 10. Key Insight: Specs as Durable Assets

> "If your specs are good enough and your validation is strong enough, you can regenerate the entire application from specs at any time."

This is the "continuous integration" of AI development. Specs are:
- **Hierarchical** -- organized as a tree, enabling progressive disclosure
- **Reviewable** by humans at a higher level than code
- **Portable** -- framework-agnostic and reusable across different technology stacks
- **Iterable** -- you can improve specs without touching code
- **Testable** -- specs define acceptance criteria that agents validate against

The same specs can drive implementations across different frameworks, enabling apples-to-apples comparison of technology choices.
