---
name: brownfield-adoption
description: >
  Step-by-step process for adopting Spec-Driven Development on an existing codebase.
  Covers the 6-step brownfield process, bootstrap prompt design, spec validation against
  existing behavior, and the decision between brownfield adoption vs deliberate rewrite.
  Trigger phrases: "brownfield", "existing codebase", "add SDD to existing project",
  "adopt SDD", "layer specs on code", "retrofit specs"
---

# Brownfield Adoption: Adding SDD to Existing Codebases

Brownfield adoption layers specs on top of existing code without rewriting it. The existing codebase becomes reference material, and specs are reverse-engineered from what the code actually does. Once specs exist, all future changes flow through the SDD lifecycle.

**Core principle:** The existing code is not the enemy -- it is the source of truth for spec generation. Respect what works; spec what matters.

---

## 1. When to Use Brownfield Adoption

Brownfield adoption is the right choice when:

- You have a **working codebase** that you want to improve incrementally
- You want to adopt SDD **without stopping development**
- The codebase is too large or critical for a full rewrite
- You want **traceability** between specs and code for future changes
- You need to **onboard AI agents** to an existing project safely
- The team wants to start with SDD on a subset of the codebase

**Brownfield is NOT the right choice when:**
- You are migrating to a completely different framework (use a deliberate rewrite instead)
- The existing code is so broken that specs would just document bugs
- The codebase is being sunset or replaced

---

## 2. Brownfield vs Deliberate Rewrite

Before starting, decide which approach fits your situation:

| Aspect | Brownfield Adoption | Deliberate Rewrite |
|--------|--------------------|--------------------|
| **Goal** | Layer specs on existing code | Replace code with new implementation |
| **Code fate** | Kept and evolved | Discarded after spec extraction |
| **Risk** | Lower -- existing code continues working | Higher -- new code must reach feature parity |
| **Speed** | Faster to start, incremental progress | Slower to start, bigger payoff if successful |
| **Best for** | Legacy codebases, incremental SDD adoption, production systems | Framework migrations, clean-slate rebuilds, tech debt bankruptcy |
| **Spec source** | Reverse-engineered from existing code | Forward-designed from requirements |
| **When code is broken** | Specs document current (broken) behavior, then fix via SDD | Specs document intended behavior, build fresh |
| **Team disruption** | Minimal -- development continues | Significant -- parallel development required |

### Decision flowchart

```
Is the existing code fundamentally sound?
  YES -> Are you changing frameworks?
           YES -> Deliberate Rewrite (extract specs, build new)
           NO  -> Brownfield Adoption (layer specs, evolve)
  NO  -> Is a rewrite feasible (time, budget, risk)?
           YES -> Deliberate Rewrite
           NO  -> Brownfield Adoption (spec the broken parts, fix incrementally)
```

---

## 3. The 6-Step Brownfield Process

### Step 1: Set Up the Context Directory

Create the standard SDD context directory structure alongside your existing codebase:

```bash
mkdir -p context/{refs,specs,plans,impl,prompts}
```

Resulting structure:

```
your-project/
+-- src/                    # Existing source code (untouched)
+-- tests/                  # Existing tests (untouched)
+-- package.json            # Existing config (untouched)
+-- context/
    +-- refs/
    |   +-- architecture-overview.md   # High-level description of existing system
    +-- specs/
    |   +-- CLAUDE.md                  # "Specs define WHAT needs implementing"
    +-- plans/
    |   +-- CLAUDE.md                  # "Plans define HOW to implement something"
    +-- impl/
    |   +-- CLAUDE.md                  # "Impls record implementation progress"
    +-- prompts/
        +-- 000-generate-specs-from-code.md   # Bootstrap prompt (this step)
```

**Create `context/refs/architecture-overview.md`** with a high-level description of the existing system:

```markdown
# Architecture Overview

## System Description
{Brief description of what the application does}

## Technology Stack
- Language: {LANGUAGE}
- Framework: {FRAMEWORK}
- Build: {BUILD_COMMAND}
- Test: {TEST_COMMAND}

## Directory Structure
{Key directories and their purposes}

## Key Domains
{List the major functional areas of the application}

## External Dependencies
{APIs, databases, services the application depends on}

## Known Issues / Tech Debt
{Major known issues that specs should account for}
```

### Step 2: Designate the Codebase as Reference Material

The existing codebase itself becomes the reference material. Unlike greenfield projects (where refs are PRDs or language specs), brownfield refs are the living code.

**In `context/refs/`, add a pointer:**

```markdown
# Reference: Existing Codebase

The existing source code at `src/` is the primary reference material for spec generation.

## How to Use This Reference
1. Explore the codebase structure to identify domains
2. Read source files to understand current behavior
3. Run existing tests to understand expected behavior
4. Check git history for context on design decisions

## What the Codebase Tells Us
- Current behavior (what the code DOES)
- Implicit requirements (what the code assumes)
- Test coverage (what is validated)
- Architecture decisions (how domains interact)

## What the Codebase Does NOT Tell Us
- Why decisions were made (check git history, docs)
- What behavior is intentional vs accidental
- What requirements are missing
- What the system SHOULD do vs what it DOES
```

### Step 3: Create the Bootstrap Prompt (000)

The bootstrap prompt is numbered `000` because it runs first and only once. It reverse-engineers specs from the existing code.

```markdown
# 000: Generate Specs from Existing Code (Brownfield Bootstrap)

## Runtime Inputs
- Framework: {FRAMEWORK}
- Build command: {BUILD_COMMAND}
- Test command: {TEST_COMMAND}
- Source directory: {SRC_DIR}

## Context
This is a brownfield adoption. The existing codebase at `{SRC_DIR}` is the reference material.
Read `context/refs/architecture-overview.md` for system context.

## Task

### Phase 1: Explore and Discover
1. Read the architecture overview
2. Explore the source directory structure
3. Identify distinct functional domains (auth, data, UI, API, etc.)
4. Read key source files in each domain
5. Run existing tests to understand expected behavior: `{TEST_COMMAND}`

### Phase 2: Generate Specs
For each identified domain:
1. Create `context/specs/spec-{domain}.md`
2. Each spec must include:
   - **Scope:** What this domain covers
   - **Requirements:** What the code currently does, expressed as requirements
   - **Acceptance Criteria:** Testable criteria derived from existing behavior
   - **Dependencies:** What other domains this depends on
   - **Out of Scope:** What this spec explicitly excludes
   - **Cross-References:** Links to related specs

3. Create `context/specs/spec-overview.md` as the index:
   - One-line summary per domain spec
   - Dependency graph between domains
   - Overall system architecture summary

### Phase 3: Validate
For each acceptance criterion in the generated specs:
1. Verify the existing code satisfies it
2. If a test exists that validates it, reference the test
3. If no test exists, note it as a coverage gap

## Exit Criteria
- [ ] All major domains have corresponding spec files
- [ ] Every requirement has testable acceptance criteria
- [ ] spec-overview.md indexes all specs
- [ ] Validation report shows which criteria are covered by existing tests
- [ ] Coverage gaps are documented

## Completion Signal
<all-tasks-complete>
```

### Step 4: Run the Iteration Loop

Run the bootstrap prompt through the iteration loop:

```bash
# Run 3-5 iterations to stabilize specs
iteration-loop context/prompts/000-generate-specs-from-code.md -n 5 -t 1h
```

**What happens during iteration:**
- **Iteration 1:** Agent explores codebase, generates initial specs (broad but shallow)
- **Iteration 2:** Agent refines specs based on git history from iteration 1, adds detail
- **Iteration 3:** Agent validates specs against code, fills coverage gaps
- **Iterations 4-5:** Convergence -- minor refinements, polishing cross-references

**Watch for convergence:** Specs should stabilize after 3-5 iterations. If they do not, the codebase may be too large for a single prompt. Split into domain-specific bootstrap prompts.

### Step 5: Validate Specs Match Behavior

After the bootstrap prompt converges, validate that the generated specs accurately describe the existing code:

#### 5a. Run tests against specs

```bash
# Use TDD to verify specs match behavior
# For each domain spec, generate tests from acceptance criteria
# then verify existing code passes them
{TEST_COMMAND}
```

#### 5b. Manual review checklist

```markdown
## Spec Validation Checklist
- [ ] Each domain in the codebase has a corresponding spec
- [ ] Acceptance criteria match actual code behavior (not aspirational)
- [ ] Dependencies between specs match actual code dependencies
- [ ] No orphan code -- every significant module is covered by a spec
- [ ] No phantom requirements -- specs do not describe behavior that does not exist
- [ ] Cross-references are accurate
```

#### 5c. Handle mismatches

| Mismatch Type | Action |
|--------------|--------|
| **Spec describes behavior that does not exist** | Remove the requirement (phantom requirement) |
| **Code has behavior not in any spec** | Add a requirement (coverage gap) |
| **Spec and code disagree on behavior** | Determine which is correct; update the other |
| **Code has bugs that specs documented as-is** | Mark as known issue in spec; fix via normal SDD |

### Step 6: Proceed with Normal SPIIM

Once specs are validated, the project is ready for full SDD. All future changes flow through specs first:

```
Future change workflow:
  1. Update spec with new/changed requirement
  2. Generate/update plans from specs (prompt 002)
  3. Implement from plans (prompt 003)
  4. Validate: build + test + acceptance criteria
  5. If issues found: backpropagate to specs
```

Create the standard pipeline prompts:

```bash
# Create greenfield-style prompts for ongoing development
# (000 was the bootstrap; 001-003 are the ongoing pipeline)
context/prompts/001-generate-specs-from-refs.md    # For new features
context/prompts/002-generate-plans-from-specs.md   # Plan generation
context/prompts/003-generate-impl-from-plans.md    # Implementation
```

---

## 4. Incremental Adoption Strategy

You do not have to spec the entire codebase at once. Start with the most active or highest-risk areas:

### Priority matrix for spec coverage

| Priority | Criteria | Example |
|----------|----------|---------|
| **P0: Spec immediately** | Code changes frequently, high risk, many bugs | Auth system, payment processing |
| **P1: Spec soon** | Active development area, moderate complexity | Feature modules, API endpoints |
| **P2: Spec when touched** | Stable code, rarely changes | Utility libraries, config modules |
| **P3: Skip until needed** | Dead code, deprecated features | Legacy compatibility layers |

### Incremental process

```
Week 1: Bootstrap specs for P0 domains
  -> Run 000 prompt scoped to P0 directories only
  -> Validate and refine

Week 2-3: Extend to P1 domains
  -> Add P1 directories to the bootstrap prompt
  -> Cross-reference with existing P0 specs

Week 4+: Spec-on-touch
  -> When any P2 file is modified, generate its spec first
  -> Gradually expand coverage
```

### Scoping the bootstrap prompt

For incremental adoption, modify prompt 000 to target specific directories:

```markdown
## Scope
This bootstrap targets the following domains only:
- `src/auth/` -> spec-auth.md
- `src/payments/` -> spec-payments.md

Do NOT generate specs for other directories at this time.
```

---

## 5. Common Challenges and Solutions

### Challenge: Codebase is too large for one context window

**Solution:** Split the bootstrap into domain-specific prompts:

```
context/prompts/
+-- 000a-generate-specs-auth.md
+-- 000b-generate-specs-data.md
+-- 000c-generate-specs-ui.md
```

Run each independently, then create a manual `spec-overview.md` that ties them together.

### Challenge: No existing tests

**Solution:** The bootstrap prompt generates specs from code behavior, not tests. After specs exist, use the implementation prompt to generate tests:

```bash
# After bootstrap, generate tests from specs
iteration-loop context/prompts/003-generate-impl-from-plans.md -n 5 -t 1h
# Focus on test generation, not code changes
```

### Challenge: Code has undocumented behavior

**Solution:** Use git history to understand intent:

```markdown
# In the bootstrap prompt, add:

## Discovery Strategy
1. Read source code for current behavior
2. Read `git log --oneline -50` for recent changes
3. Read `git log --follow {file}` for individual file history
4. Infer requirements from both code AND history
```

### Challenge: Code has known bugs

**Solution:** Spec the intended behavior, not the buggy behavior. Mark known bugs as issues:

```markdown
### R3: Search Results Pagination
**Description:** Search results are paginated with 20 items per page
**Acceptance Criteria:**
- [ ] Results are paginated
- [ ] Page size is configurable (default 20)
**Known Issues:**
- BUG: Off-by-one error on last page (see issue #142)
```

### Challenge: Team resistance to SDD

**Solution:** Start small, show results:
1. Pick ONE upcoming feature
2. Write a spec before implementing it
3. Show how the spec caught issues the team would have missed
4. Gradually expand SDD coverage based on demonstrated value

---

## 6. Lightweight SDD for Small Projects

Even small projects benefit from minimal SDD. The "SDD floor" is:

```
your-small-project/
+-- src/
+-- context/
    +-- specs/
    |   +-- spec-task.md          # One spec for the current task
    +-- plans/
        +-- plan-task.md          # One plan for the current task
```

**No prompts directory needed.** Just write a focused spec and plan, then use the iteration loop against the plan.

**Why bother for small projects?**
- The spec catches requirements you would have missed
- The plan sequences work so the agent does not thrash
- If the project grows, you already have the structure in place
- It is much easier to scale up from lightweight SDD than to retrofit full SDD later

### Lightweight SDD process

1. Write `context/specs/spec-task.md` (15-30 minutes)
2. Write `context/plans/plan-task.md` (10-20 minutes)
3. Run the iteration loop against the plan
4. If the project grows, add the full context directory structure

---

## 7. Transition Milestones

Track your brownfield adoption progress with these milestones:

```markdown
## Brownfield Adoption Progress

### Milestone 1: Foundation
- [ ] Context directory created
- [ ] Architecture overview written
- [ ] Bootstrap prompt created

### Milestone 2: Initial Specs
- [ ] P0 domains have specs
- [ ] Specs validated against existing code
- [ ] Coverage gaps documented

### Milestone 3: Pipeline Active
- [ ] Standard prompts (001-003) created
- [ ] First feature developed through SDD pipeline
- [ ] Backpropagation process tested

### Milestone 4: Steady State
- [ ] All active domains have specs
- [ ] All new features go through specs first
- [ ] Backpropagation is routine
- [ ] Iteration loop runs are predictable (convergence in 3-5 iterations)

### Milestone 5: Full SDD
- [ ] All domains have specs
- [ ] All changes flow through SPIIM
- [ ] Convergence monitoring active
- [ ] Team comfortable with the process
```

---

## Cross-References

- **Context architecture:** See `sdd:context-architecture` skill for the full context directory structure and progressive disclosure patterns.
- **Prompt pipeline:** See `sdd:prompt-pipeline` skill for designing the 001-003 prompts after bootstrap.
- **Spec writing:** See `sdd:spec-writing` skill for how to write high-quality specs with testable acceptance criteria.
- **Backpropagation:** See `sdd:backpropagation` skill for tracing bugs back to specs after brownfield adoption.
- **Convergence monitoring:** See `sdd:convergence-monitoring` skill for detecting when the bootstrap prompt has converged.
