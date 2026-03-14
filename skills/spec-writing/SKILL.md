---
name: spec-writing
description: |
  How to write SDD-quality specifications that AI agents can consume effectively. Covers
  implementation-agnostic spec design, testable acceptance criteria, hierarchical structure,
  cross-referencing, spec templates, greenfield and rewrite patterns, spec compaction, and gap analysis.
  Trigger phrases: "write specs", "create specifications", "spec this out",
  "define requirements for agents", "how to write specs for AI"
---

# Spec Writing for SDD

## Core Principle: Specs Describe WHAT, Not HOW

Specifications are **implementation-agnostic**. They define what the system must do and how to verify it, but never prescribe a specific framework, language, or architecture.

This is the fundamental distinction in SDD:
- **Specs** = WHAT must be true (framework-agnostic, durable, portable)
- **Plans** = HOW to build it (framework-specific, derived from specs)
- **Code** = the implementation (generated from plans, validated against specs)

### Why Implementation-Agnostic?

When specs avoid prescribing HOW, they become:
- **Portable** — the same specs can drive implementations in different frameworks
- **Durable** — specs survive technology migrations
- **Testable** — acceptance criteria are about behavior, not implementation details
- **Reusable** — the same specs work for greenfield, rewrites, and cross-framework evaluation

**Bad spec requirement:** "Use React useState hook to manage form state"
**Good spec requirement:** "Form state persists across user interactions within a session. Acceptance: entering values, navigating away, and returning preserves all entered values."

---

## Every Requirement Needs Testable Acceptance Criteria

This is the single most important rule in SDD spec writing. If an agent cannot automatically validate a requirement, that requirement will not be met.

### The Validation-First Rule

Every requirement must answer: **"How would an automated test verify this?"**

| Weak Criterion | Strong Criterion |
|----------------|-----------------|
| "UI should look good" | "All interactive elements have minimum 44x44px touch targets" |
| "System should be fast" | "API responses return within 200ms at p95 under 100 concurrent users" |
| "Handle errors gracefully" | "Network failures display a retry prompt with exponential backoff (1s, 2s, 4s)" |
| "Support authentication" | "Valid credentials return a session token; invalid credentials return 401 with error message" |

### Acceptance Criteria Format

Each criterion should be:
- **Observable** — can be checked by reading output, UI state, or logs
- **Deterministic** — same input always produces same pass/fail result
- **Automatable** — an agent can write a test that checks this
- **Independent** — does not depend on subjective judgment

```markdown
**Acceptance Criteria:**
- [ ] {Action} results in {observable outcome}
- [ ] Given {precondition}, when {action}, then {result}
- [ ] {Metric} meets {threshold} under {conditions}
```

---

## Hierarchical Structure with Index

Specs must be organized as a hierarchy — one index file linking to domain-specific sub-specs. This enables progressive disclosure: agents read the index first, then only the sub-specs relevant to their task.

### The Spec Index Pattern

Create a `spec-overview.md` as the entry point:

```markdown
# Spec Overview

## Domains

| Domain | Spec File | Summary |
|--------|-----------|---------|
| Authentication | spec-auth.md | User registration, login, session management, OAuth |
| Data Models | spec-data-models.md | Core entities, relationships, validation rules |
| API | spec-api.md | REST endpoints, request/response formats, error handling |
| UI Components | spec-ui-components.md | Shared components, accessibility, responsive behavior |
| Notifications | spec-notifications.md | Email, push, in-app notification delivery |

## Cross-Cutting Concerns
- Security requirements: see spec-auth.md R3, spec-api.md R7
- Performance budgets: see spec-api.md R12, spec-ui-components.md R5
- Accessibility: see spec-ui-components.md R8-R10
```

### Why Hierarchical?

1. **Context window efficiency** — agents load only the domains they need
2. **Parallel work** — different agents can own different spec domains
3. **Review efficiency** — humans can review domain-by-domain
4. **Cross-referencing** — domains link to each other explicitly

---

## Cross-Referencing Between Specs

Related specs must link to each other. Cross-references prevent requirements from being lost at domain boundaries.

### Cross-Reference Patterns

```markdown
## Cross-References
- **Depends on:** spec-auth.md R1 (session tokens required for API access)
- **Depended on by:** spec-notifications.md R4 (uses user preferences from this spec)
- **Related:** spec-ui-components.md R6 (error display components used by this domain)
```

### When to Cross-Reference

- When one domain's requirement depends on another domain's output
- When shared entities are defined in one spec but used in many
- When validation criteria span multiple domains
- When out-of-scope items are in-scope for another spec

---

## Full Spec Format Template

Use this template for every domain spec:

```markdown
# Spec: {Domain Name}

## Scope
{One paragraph describing what this spec covers and its boundaries.}

## Requirements

### R1: {Requirement Name}
**Description:** {What must be true — stated in terms of behavior, not implementation.}
**Acceptance Criteria:**
- [ ] {Testable criterion 1}
- [ ] {Testable criterion 2}
- [ ] {Testable criterion 3}
**Dependencies:** {Other specs/requirements this depends on, or "None"}

### R2: {Requirement Name}
**Description:** {What must be true}
**Acceptance Criteria:**
- [ ] {Testable criterion 1}
- [ ] {Testable criterion 2}
**Dependencies:** {Dependencies}

### R3: ...

## Out of Scope
{Explicit list of things this spec does NOT cover. This is critical — it prevents
agents from over-building and clarifies domain boundaries.}
- {Thing explicitly excluded and why}
- {Another exclusion}

## Cross-References
- See also: spec-{related-domain}.md — {why it is related}
- Depends on: spec-{dependency}.md R{N} — {what is needed}
- Depended on by: spec-{dependent}.md R{N} — {what depends on this}
```

### Template Rules

1. **Number requirements sequentially** (R1, R2, R3...) — agents reference them by ID
2. **Every requirement gets acceptance criteria** — no exceptions
3. **Out of Scope is mandatory** — explicit exclusions prevent scope creep
4. **Cross-References section is mandatory** — even if it says "None"
5. **Scope section is one paragraph** — concise boundary description

---

## Greenfield Pattern: Reference Material → Specs

When building from scratch, you start with reference materials and derive specs from them.

### Flow

```
context/refs/              context/specs/
├── prd.md          →      ├── spec-overview.md
├── design-doc.md   →      ├── spec-auth.md
├── api-draft.md    →      ├── spec-api.md
└── research/       →      ├── spec-data-models.md
    └── ...         →      └── spec-ui.md
```

### Process

1. **Place all reference materials** in `context/refs/`
2. **Run spec generation** — agent reads all refs, decomposes into domains
3. **Agent produces:**
   - `spec-overview.md` — index with domain summaries
   - One `spec-{domain}.md` per identified domain
   - Cross-references between related domains
4. **Human reviews** specs for completeness and correctness
5. **Iterate** — refine specs based on review feedback

### Greenfield Prompt Pattern

The first prompt in a greenfield pipeline (typically `001-generate-specs-from-refs.md`) should:
- Read all files in `context/refs/`
- Decompose reference material into domains
- Generate specs following the template above
- Create `spec-overview.md` as the index
- Cross-reference related specs

---

## Rewrite Pattern: Old Code → Reference Docs → Specs

When rewriting an existing system, the existing code becomes your reference material. But you never go directly from old code to new code — you always extract specs first.

### Flow

```
Existing codebase          context/refs/              context/specs/
├── src/            →      ├── ref-apis.md      →     ├── spec-overview.md
├── tests/          →      ├── ref-data-models.md →   ├── spec-auth.md
└── docs/           →      ├── ref-ui-components.md →  ├── spec-api.md
                           └── ref-architecture.md →   └── spec-data.md
```

### Process

1. **Agent explores the existing codebase** and generates reference documents
2. **Reference docs capture** the current system's behavior, APIs, data models, and UI patterns
3. **Agent generates specs** from reference docs — implementation-agnostic requirements
4. **Validate specs against existing code** — verify acceptance criteria match current behavior
5. **Proceed with normal SPIIM** — specs drive the new implementation

### Rewrite Prompt Pattern

Rewrites typically use more prompts because of the reverse-engineering step:
- `001`: Generate reference materials from old code
- `002`: Generate specs from references + feature scope
- `003`: Validate specs against existing codebase
- `004+`: Plans and implementation

The key difference from greenfield: step 003 validates that your specs actually describe what the old system does, before you start building the new one.

---

## Spec Compaction

When implementation tracking or spec files grow beyond approximately 500 lines, they become unwieldy for agents to process efficiently. Spec compaction compresses large files while preserving active context.

### When to Compact

- Implementation tracking file exceeds 500 lines
- Spec file has many resolved/completed requirements mixed with active ones
- Agent is spending too much context window on historical information

### How to Compact

1. **Identify resolved content:** completed tasks, resolved issues, archived dead ends
2. **Archive removed content** to a separate file (e.g., `impl/archive/impl-domain-v1.md`)
3. **Preserve in the compacted file:**
   - All active/in-progress tasks
   - All open issues
   - Recent dead ends (last 2-3 sessions)
   - Current test health status
   - Active cross-references
4. **Target:** under 500 lines in the active file

### Compaction Rule

Never delete information — move it to an archive. Agents can still find archived context if needed, but it will not consume context window during normal operations.

---

## Gap Analysis

Gap analysis compares what was built against what was intended, identifying where specs, plans, or validation fell short.

### How to Perform Gap Analysis

1. **Read specs** (intended behavior) and **implementation tracking** (what was built)
2. **For each spec requirement,** check if acceptance criteria are satisfied
3. **Classify each requirement:**

| Status | Meaning |
|--------|---------|
| **Complete** | All acceptance criteria pass |
| **Partial** | Some criteria pass, others do not |
| **Missing** | Requirement not implemented at all |
| **Over-built** | Implementation exceeds spec (may indicate spec gap) |

4. **Report gaps** with: which spec, which criterion, what is missing
5. **Feed gaps into backpropagation** — update specs if needed, then re-implement

### Gap Analysis as Feedback

Gap analysis is not a one-time activity. Run it:
- After each implementation iteration
- Before starting a new session (to prioritize work)
- When convergence stalls (to identify what is blocking progress)

---

## Integration with Other Skills

### With `superpowers:brainstorming`

Use brainstorming during the Spec phase to explore requirements you may not have considered. Brainstorm first, then formalize the output into spec format with acceptance criteria.

**Pattern:**
1. Brainstorm the domain with the agent
2. Capture insights as draft requirements
3. Formalize each requirement with acceptance criteria
4. Add to the spec using the template

### With `sdd:validation-first`

Every acceptance criterion in a spec must map to at least one validation gate. When writing specs, think about which gate will verify each requirement:

| Acceptance Criterion Type | Likely Gate |
|--------------------------|-------------|
| "Code compiles without errors" | Gate 1: Build |
| "Function returns correct output for input X" | Gate 2: Unit Tests |
| "User can complete workflow end-to-end" | Gate 3: E2E/Integration |
| "Response time under N ms" | Gate 4: Performance |
| "Application starts and displays main screen" | Gate 5: Launch Verification |
| "UI matches design intent" | Gate 6: Human Review |

### With `sdd:context-architecture`

Specs live in the `context/specs/` directory. See `sdd:context-architecture` for the full context directory structure, CLAUDE.md conventions, and multi-repo strategies.

### With `sdd:impl-tracking`

As specs are implemented, progress is tracked in `context/impl/` documents. Dead ends discovered during implementation should be recorded to prevent future agents from retrying failed approaches.

---

## Common Mistakes

### 1. Writing Implementation-Specific Specs

**Wrong:** "Use PostgreSQL with a users table containing columns: id (UUID), email (VARCHAR), ..."
**Right:** "User accounts have a unique identifier and email. Email must be unique across all accounts. Acceptance: creating two accounts with the same email fails with a duplicate error."

### 2. Vague Acceptance Criteria

**Wrong:** "System handles errors properly"
**Right:** "When a network request fails, the UI displays an error message within 2 seconds and offers a retry action. Acceptance: simulating network failure shows error banner with retry button."

### 3. Missing Out of Scope

Every spec needs explicit exclusions. Without them, agents will over-build or make assumptions.

### 4. No Cross-References

Domains do not exist in isolation. If spec-auth defines session tokens that spec-api uses, both specs must cross-reference each other.

### 5. Monolithic Specs

A single 1000-line spec file defeats progressive disclosure. Decompose into domains with a clear index.

---

## Summary

Writing specs for AI agents follows these rules:

1. **WHAT, not HOW** — describe behavior, not implementation
2. **Every requirement gets testable acceptance criteria** — if agents cannot validate it, it will not be met
3. **Hierarchical with an index** — progressive disclosure for context efficiency
4. **Cross-referenced** — related domains link to each other
5. **Explicitly scoped** — out-of-scope section prevents over-building
6. **Compact when large** — archive resolved content, keep active files under 500 lines
7. **Living documents** — specs evolve through backpropagation as gaps are discovered
