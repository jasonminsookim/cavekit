---
name: adversarial-review
description: >
  Patterns for using a second AI agent or model to challenge the primary builder agent's work.
  Covers six integration strategies (Code Review, Architecture Review, Shared Thread, Agent-to-Agent,
  Tie-Breaking, Test Coverage Review), how to set up adversarial review with any model via MCP server,
  adversarial iteration loops that alternate builder and reviewer prompts, and prompt templates for
  each strategy. The adversary's job is to find what the builder missed, not to agree.
  Triggers: "adversarial review", "adversarial agent", "use another model to review",
  "second opinion on code", "cross-model review".
---

# Adversarial Review

Use a second AI agent to review and challenge the first agent's work. The adversary exists to find
what the builder missed -- not to agree, not to be polite, and not to rubber-stamp. This is the
single most effective quality gate you can add beyond automated tests.

## Core Principle

> **The adversary's job is to find what the builder missed, not to agree.**

A review that says "looks good" is a wasted review. The adversarial model should be given explicit
instructions to be critical, to challenge assumptions, and to look for what is *not* there rather
than what is.

---

## Why Adversarial Review Works

LLMs have blind spots. Every model has patterns it over-relies on, edge cases it misses, and
architectural assumptions it makes implicitly. A second model -- or the same model with a different
prompt and role -- catches a different set of issues.

**The analogy:** In traditional engineering, code review exists because the author has cognitive
blind spots about their own work. The same principle applies to AI agents, but the blind spots are
different: they are systematic patterns in training data, context window limitations, and prompt
interpretation biases.

**What adversarial review catches that automated tests miss:**
- Architectural over-engineering or under-engineering
- Missing error handling patterns
- Security vulnerabilities the builder didn't consider
- Spec requirements that were technically met but poorly implemented
- Dead code, unused imports, and unnecessary complexity
- Performance pitfalls that only manifest at scale
- Missing edge cases not covered by the spec

---

## Integration Strategies

| Strategy | When | How |
|----------|------|-----|
| **Code Review** | After implementation | Second model reviews changed files with a targeted review prompt; first model applies fixes |
| **Architecture Review** | During planning | Second model designs alternative architecture; first model validates both against specs and picks the stronger approach |
| **Shared Thread** | Complex discussions | Multiple replies on the same conversation thread, maintaining context across turns |
| **Agent-to-Agent** | Large tasks | A teammate agent owns the adversary interaction, manages the back-and-forth, and reports findings to the team lead |
| **Tie-Breaking** | Conflicting opinions | Team lead presents both perspectives to the adversarial model and uses its analysis to decide |
| **Test Coverage Review** | During validation | Test analysis output (coverage reports, gap analysis) feeds into adversarial review for deeper scrutiny |

### When to Use Each Strategy

```
Need adversarial review
├─ Reviewing completed code?
│   ├─ Small changeset (< 500 lines) → Code Review
│   └─ Large changeset or full feature → Agent-to-Agent
├─ Designing architecture?
│   ├─ Single decision point → Tie-Breaking
│   └─ Full system design → Architecture Review
├─ Debating trade-offs?
│   ├─ Need extended back-and-forth → Shared Thread
│   └─ Need a decisive answer → Tie-Breaking
└─ Validating test quality?
    └─ Test Coverage Review
```

---

## Setting Up Adversarial Review via MCP Server

Any AI model that exposes an MCP server interface can serve as an adversarial reviewer. The setup
is model-agnostic -- the pattern works with any model that supports the MCP protocol.

### Generic MCP Configuration

Add the adversarial model as an MCP server in your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "adversary": {
      "command": "{ADVERSARY_CLI}",
      "args": ["mcp-server"],
      "env": {
        "API_KEY": "{ADVERSARY_API_KEY}"
      }
    }
  }
}
```

Replace `{ADVERSARY_CLI}` with the CLI command for your chosen model (e.g., any model's CLI tool
that supports MCP server mode) and `{ADVERSARY_API_KEY}` with the appropriate credentials.

### Two Core MCP Tools

Most adversarial model MCP servers expose two tools:

1. **Start session** -- Begin a new conversation with the adversarial model
   - Parameters: prompt, approval policy, sandbox mode, model selection
   - Returns: a thread/session identifier

2. **Reply to session** -- Continue an existing conversation
   - Parameters: thread/session ID, follow-up message
   - Returns: the model's response

The thread/session identifier is critical -- it allows multi-turn conversations where the adversary
builds on previous context.

### Example: Starting an Adversarial Review Session

```
Tool: adversary.start_session
Parameters:
  prompt: "Review the following code changes for bugs, security issues,
           missing edge cases, and spec compliance. Be critical -- your
           job is to find problems, not to agree. Here are the changes:
           {DIFF_CONTENT}"
  model: "{ADVERSARY_MODEL}"
```

### Example: Multi-Turn Follow-Up

```
Tool: adversary.reply_to_session
Parameters:
  thread_id: "{THREAD_ID_FROM_PREVIOUS}"
  message: "Good findings. Now focus specifically on error handling paths.
            For each function that can fail, verify there is explicit
            error handling and that errors propagate correctly."
```

---

## Strategy Details

### 1. Code Review

**When:** After a builder agent completes implementation of a feature or fix.

**Process:**
1. Builder agent implements the feature and commits
2. Generate a diff of all changes: `git diff {BASE_BRANCH}...HEAD`
3. Send the diff to the adversarial model with a code review prompt
4. Parse the adversary's findings into actionable items
5. Builder agent applies fixes for valid findings
6. Optionally: send fixes back to adversary for re-review

**Review Prompt Template:**
```markdown
You are a senior code reviewer. Review the following code changes critically.

## What to look for:
- Bugs, logic errors, off-by-one errors
- Security vulnerabilities (injection, auth bypass, data exposure)
- Missing error handling and edge cases
- Performance issues (N+1 queries, unnecessary allocations, blocking calls)
- Spec compliance: does this implementation match the requirements?
- Code quality: naming, structure, unnecessary complexity

## What NOT to do:
- Do not say "looks good" unless you genuinely found zero issues
- Do not suggest stylistic changes unless they affect readability significantly
- Do not rewrite the code -- describe the problem and where it is

## Spec requirements for this feature:
{SPEC_REQUIREMENTS}

## Code changes:
{DIFF_CONTENT}

## Output format:
For each finding:
- **Severity:** CRITICAL / HIGH / MEDIUM / LOW
- **File:** path and line range
- **Issue:** what is wrong
- **Why:** why this matters
- **Suggestion:** how to fix it
```

### 2. Architecture Review

**When:** During the planning phase, before implementation begins.

**Process:**
1. Builder agent drafts an architecture or plan
2. Send the plan + specs to the adversarial model
3. Adversary proposes alternative approaches or critiques the plan
4. Builder validates both approaches against specs
5. Human makes the final decision if there is a genuine trade-off

**Architecture Review Prompt Template:**
```markdown
You are a systems architect reviewing a proposed design. Your goal is to
find weaknesses, over-engineering, missing considerations, and better
alternatives.

## Specs (what must be built):
{SPEC_CONTENT}

## Proposed architecture:
{PLAN_CONTENT}

## Evaluate:
1. Does this architecture satisfy all spec requirements?
2. Is it over-engineered for the scope?
3. Are there simpler alternatives that meet the same requirements?
4. What failure modes exist? How does the system recover?
5. What are the scaling bottlenecks?
6. What dependencies introduce risk?
```

### 3. Shared Thread

**When:** Complex design discussions that require extended back-and-forth.

**Process:**
1. Start a session with the adversarial model presenting the problem
2. Use reply-to-session to continue the conversation across multiple turns
3. Maintain the thread ID throughout the discussion
4. Summarize conclusions when the discussion converges

**Key consideration:** Thread-based conversations accumulate context. Keep the
conversation focused on a single topic to avoid context dilution.

### 4. Agent-to-Agent

**When:** Large tasks where the adversarial review itself is substantial.

**Process:**
1. Team lead spawns a teammate specifically for adversarial coordination
2. The teammate owns the adversary MCP interaction
3. Teammate manages multi-turn review sessions
4. Teammate summarizes findings and reports to the team lead
5. Team lead assigns fixes to the appropriate builder teammates

**Why delegate:** The adversarial back-and-forth can consume significant context
window. Delegating it to a dedicated teammate preserves the team lead's context
for coordination.

### 5. Tie-Breaking

**When:** The builder agent and human (or two agents) disagree on an approach.

**Process:**
1. Present both perspectives to the adversarial model
2. Ask it to evaluate the trade-offs of each approach
3. Ask it to recommend one, with explicit reasoning
4. Use the recommendation to inform the decision (human has final say)

**Tie-Breaking Prompt Template:**
```markdown
Two approaches have been proposed for the same problem. Evaluate both
critically and recommend one.

## Context:
{PROBLEM_DESCRIPTION}

## Approach A:
{APPROACH_A}

## Approach B:
{APPROACH_B}

## Evaluation criteria:
- Correctness: which approach is more likely to be correct?
- Simplicity: which is easier to understand and maintain?
- Performance: which performs better for the expected use case?
- Risk: which has fewer failure modes?

## Your recommendation:
Pick one and explain why. If neither is clearly better, say so and
explain what additional information would break the tie.
```

### 6. Test Coverage Review

**When:** During validation, after tests have been generated and run.

**Process:**
1. Run test coverage analysis on the codebase
2. Generate a coverage report (which files/functions are covered)
3. Send the coverage report + specs to the adversarial model
4. Adversary identifies: untested edge cases, missing integration tests,
   spec requirements without corresponding tests
5. Builder adds missing tests

---

## Adversarial Iteration (Convergence Loop with Review)

Instead of a simple build-then-review, run alternating convergence loops where each
iteration alternates between building and reviewing.

### The Pattern

```
Iteration 1: Builder runs against spec → produces code
Iteration 2: Reviewer runs against code + spec → produces findings
Iteration 3: Builder runs against spec + findings → fixes code
Iteration 4: Reviewer runs against updated code + spec → produces new findings
...repeat until findings converge to zero (or trivial)
```

### Implementation with Separate Prompts

Create two prompt files:

**`prompts/build.md`** -- The builder prompt:
```markdown
Implement the requirements in the spec. Read implementation tracking for
context on what has been done. Read any review findings and address them.

Input: specs/, plans/, impl/, review-findings.md (if exists)
Output: source code, updated impl tracking
Exit: all spec requirements implemented, all review findings addressed
```

**`prompts/review.md`** -- The reviewer prompt:
```markdown
Review the current implementation against the spec. Be critical. Find
bugs, missing requirements, security issues, and quality problems.

Input: specs/, plans/, source code, impl/
Output: review-findings.md
Exit: all source files reviewed against all spec requirements
```

### Running Adversarial Iteration

```bash
# Terminal 1: Builder convergence loop
{LOOP_TOOL} prompts/build.md -n 5 -t 2h

# Terminal 2: Reviewer convergence loop (staggered by 30 min)
{LOOP_TOOL} prompts/review.md -n 5 -t 2h -d 30m
```

The builder and reviewer share the same git repository. The reviewer reads the
builder's latest committed code; the builder reads the reviewer's latest
`review-findings.md`. They converge naturally through git.

### Convergence Signal

The adversarial loop has converged when:
- The reviewer's findings drop to zero or only LOW severity items remain
- The builder's diffs between iterations are minimal
- All spec requirements have been reviewed and confirmed as met

---

## Anti-Patterns

### 1. Adversary as Yes-Man
**Problem:** The adversarial model says "looks good" without finding real issues.
**Fix:** Explicitly instruct the adversary to find problems. Add to the prompt:
"If you find zero issues, explain what areas you checked and why you believe
they are correct. An empty review is suspicious."

### 2. Adversary Rewrites Everything
**Problem:** The adversary provides complete rewrites instead of identifying issues.
**Fix:** Instruct the adversary to describe problems and locations, not to write
code. "Your output is a list of findings, not a pull request."

### 3. Builder Ignores Findings
**Problem:** The builder agent dismisses adversary findings without addressing them.
**Fix:** Require the builder to explicitly respond to each finding: "For each
review finding, either fix it and explain the fix, or explain why the finding
is not valid. You may not skip any finding."

### 4. Infinite Disagreement Loop
**Problem:** Builder and reviewer keep going back and forth without converging.
**Fix:** Set a maximum iteration count. After N iterations, escalate to human.
If the disagreement persists, it likely indicates an ambiguous spec that needs
human clarification.

### 5. Same Model Reviewing Itself
**Problem:** Using the same model with the same prompt for both building and reviewing.
**Fix:** At minimum, use different prompts with different roles. Ideally, use a
different model or a different model version. The value of adversarial review comes
from diverse perspectives.

---

## Prompt Templates Quick Reference

| Strategy | Key Prompt Instruction |
|----------|----------------------|
| Code Review | "Find bugs, security issues, missing edge cases. Do not say 'looks good'." |
| Architecture Review | "Find weaknesses and simpler alternatives. Evaluate failure modes." |
| Shared Thread | "Continue the discussion. Build on previous context." |
| Agent-to-Agent | "Own the adversary interaction. Summarize findings for the lead." |
| Tie-Breaking | "Evaluate both approaches. Recommend one with explicit reasoning." |
| Test Coverage Review | "Identify untested edge cases and spec requirements without tests." |

---

## Integration with SDD Lifecycle

Adversarial review fits into the SPIIM lifecycle at multiple points:

| SPIIM Phase | Adversarial Role |
|-------------|-----------------|
| **Spec** | Review specs for completeness, ambiguity, missing edge cases |
| **Plan** | Architecture Review: challenge the plan before implementation begins |
| **Implement** | Code Review: review implementation against specs after each feature |
| **Iterate** | Adversarial iteration loop: alternate build/review convergence |
| **Monitor** | Test Coverage Review: validate that monitoring covers all failure modes |

The most impactful point is during **Iterate** -- adversarial iteration catches issues
that neither automated tests nor single-agent convergence loops find.

---

## Cross-References

- **convergence-monitoring** -- How to detect when adversarial iterations have converged
- **validation-first** -- Adversarial review is Gate 6 (human/agent review) in the validation pipeline
- **prompt-pipeline** -- How to structure builder and reviewer prompts in the SPIIM pipeline
- **backpropagation** -- When the adversary finds a spec gap, backpropagate the fix to specs
- **impl-tracking** -- Record adversarial findings in implementation tracking documents
