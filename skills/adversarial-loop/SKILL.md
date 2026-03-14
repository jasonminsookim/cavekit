---
name: adversarial-loop
description: >
  Adversarial Ralph Loop — combines SDD specs with a Ralph Loop and true cross-model adversarial
  review using Codex (OpenAI) as an MCP server. Claude builds from specs; Codex reviews adversarially.
  Covers setup, MCP configuration, iteration patterns, convergence detection, and completion criteria.
  Triggers: "adversarial loop", "ralph loop with codex", "sdd ralph", "adversarial build loop",
  "cross-model loop", "codex adversary", "spec to ralph loop"
---

# Adversarial Loop — SDD + Ralph Loop + Codex Adversary

Run an SDD spec through a Ralph Loop where Claude builds and Codex adversarially reviews.
This is the most rigorous automated quality process available: every few iterations, a completely
different model (different training data, different biases, different blind spots) challenges
your implementation.

---

## Why This Works

| Factor | Single-Model Loop | Adversarial Loop |
|--------|-------------------|------------------|
| Blind spots | Same model, same blind spots every iteration | Two models catch different classes of issues |
| Spec drift | Builder may silently deviate from spec | Adversary checks spec compliance explicitly |
| Quality floor | Converges to "good enough for one model" | Converges to "survives cross-examination" |
| Dead ends | May retry failed approaches | Adversary flags repeated patterns |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Ralph Loop                         │
│  (Stop hook feeds same prompt each iteration)        │
│                                                      │
│  ┌──────────┐    ┌──────────────┐    ┌────────────┐ │
│  │  Claude   │───▶│ Build from   │───▶│  Commit    │ │
│  │  (Build)  │    │ spec + fixes │    │  changes   │ │
│  └──────────┘    └──────────────┘    └──────┬─────┘ │
│       ▲                                      │       │
│       │                                      ▼       │
│  ┌──────────┐    ┌──────────────┐    ┌────────────┐ │
│  │  Fix      │◀──│ Parse        │◀──│  Codex MCP │ │
│  │  findings │    │ findings     │    │  (Review)  │ │
│  └──────────┘    └──────────────┘    └────────────┘ │
│                                                      │
│  Completion: all spec requirements met +             │
│              no CRITICAL/HIGH findings               │
└─────────────────────────────────────────────────────┘
```

---

## Quick Start

```bash
# Basic: implement a spec with adversarial review
/sdd:adversarial-loop context/specs/spec-auth.md

# With options
/sdd:adversarial-loop context/specs/spec-api.md --max-iterations 20 --codex-model gpt-5.4-mini

# Review-only mode (review existing code, don't build new)
/sdd:adversarial-loop context/specs/spec-api.md --review-only

# Review every iteration instead of every 2nd
/sdd:adversarial-loop context/specs/spec-auth.md --review-interval 1
```

---

## What the Command Does

1. **Validates** the spec file exists and Codex CLI is installed
2. **Configures** Codex as an MCP server in `.mcp.json` (if not already configured)
3. **Builds** a Ralph Loop prompt that embeds:
   - The spec path and related plan/impl files
   - Instructions to alternate between build and review iterations
   - The adversarial review prompt template for Codex
   - Completion criteria tied to spec acceptance criteria
4. **Starts** the Ralph Loop via the stop hook mechanism

---

## Codex MCP Server

The command configures Codex as an MCP server automatically:

```json
{
  "mcpServers": {
    "codex-adversary": {
      "command": "codex",
      "args": ["mcp-server", "-c", "model=\"gpt-5.4\""]
    }
  }
}
```

Claude calls this MCP server on review iterations to get adversarial feedback. The MCP server
exposes Codex as a tool that accepts prompts and returns responses — Claude sends the spec +
code diff, Codex returns findings.

### Changing the Codex Model

Use `--codex-model` to specify which OpenAI model Codex should use:

```bash
/sdd:adversarial-loop spec.md --codex-model gpt-5.4-mini    # faster, cheaper
/sdd:adversarial-loop spec.md --codex-model gpt-5.4          # default, most capable
```

---

## Iteration Pattern

```
Iteration 1: BUILD  — Read spec, implement first requirement
Iteration 2: REVIEW — Call Codex MCP, get findings, fix CRITICAL/HIGH
Iteration 3: BUILD  — Continue implementing, address remaining findings
Iteration 4: REVIEW — Call Codex MCP again, new findings on new code
...
Iteration N: BUILD  — All requirements met, all findings fixed
             → outputs <promise>SPEC COMPLETE</promise>
```

The review interval is configurable. Default is every 2nd iteration.
Use `--review-interval 1` for maximum rigor (review every iteration).

---

## Adversarial Findings File

Review findings are tracked in `context/adversarial-findings.md`:

```markdown
# Adversarial Review Findings

## Latest Review: Iteration 4 — 2026-03-14T10:30:00Z
### Reviewer: Codex (gpt-5.4)

| # | Severity | File | Issue | Status |
|---|----------|------|-------|--------|
| 1 | CRITICAL | src/auth.ts:L42 | Missing input validation on token | FIXED |
| 2 | HIGH | src/auth.ts:L67 | Race condition in session refresh | FIXED |
| 3 | MEDIUM | src/auth.ts:L15 | Unused import | NEW |
| 4 | LOW | src/auth.ts:L3 | Comment typo | WONTFIX |

## History
### Iteration 2
| # | Severity | File | Issue | Status |
|---|----------|------|-------|--------|
| 1 | CRITICAL | src/auth.ts:L20 | SQL injection in login query | FIXED |
```

---

## Completion Criteria

The loop exits when the completion promise is output. The prompt instructs Claude
to ONLY output it when ALL of these are true:

- All spec requirements (R-numbers) have been implemented
- All acceptance criteria pass
- No CRITICAL or HIGH adversarial findings remain unfixed
- Build passes
- Tests pass
- At least one review iteration completed with no new CRITICAL/HIGH findings

---

## Modes

### Build + Review (default)
Alternates between implementing spec requirements and calling Codex for review.
Use for greenfield implementation from a spec.

### Review Only (`--review-only`)
Skips building. Each iteration calls Codex to review existing code against the spec,
then fixes issues found. Use when code already exists and you want adversarial QA.

---

## Prerequisites

1. **Codex CLI installed**: `npm install -g @openai/codex`
2. **OpenAI API key configured**: Codex needs authentication (via `codex login` or env var)
3. **SDD context directory**: Spec file must exist at the given path
4. **Ralph Loop plugin**: The ralph-loop plugin must be installed (provides the stop hook)

---

## Convergence Signals

The adversarial loop has converged when:
- Codex's findings drop to zero or only LOW/MEDIUM severity
- Code diffs between iterations are minimal
- All spec requirements confirmed as met by both Claude and Codex

If the loop hits max iterations without converging:
- Check `context/adversarial-findings.md` for persistent issues
- Consider whether the spec needs clarification
- Run `/sdd:back-propagate` to trace issues back to specs

---

## Cross-References

- **adversarial-review** — The underlying adversarial review patterns and prompt templates
- **convergence-monitoring** — How to detect convergence vs ceiling
- **validation-first** — Validation gates that run on every build iteration
- **impl-tracking** — How implementation progress is tracked across iterations
- **ralph-loop** — The underlying Ralph Loop mechanism
