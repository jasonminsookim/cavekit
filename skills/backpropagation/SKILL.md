---
name: backpropagation
description: >
  The technique of tracing bugs and manual fixes back to specs and prompts, then fixing at the source
  so the iteration loop can reproduce the fix autonomously. Covers the 8-step backpropagation process,
  commit classification, spec-level root cause analysis, and regression test generation.
  Trigger phrases: "backpropagate", "back-propagate", "trace bug to spec", "fix the spec not the code",
  "why did this bug happen", "update specs from bug"
---

# Backpropagation: Tracing Bugs Back to Specs

Backpropagation is borrowed from machine learning. In ML, when the output is wrong, you trace the error backward through the network to adjust weights. In SDD, when the built software has bugs or gaps, you trace the issue back to the specs and prompts and fix at the source -- not just in code.

**Key insight:** If a bug can only be fixed in code, your specs are incomplete. The goal is that specs plus the iteration loop can reproduce any fix autonomously.

---

## 1. Why Backpropagation Matters

Without backpropagation, every bug fix is a one-off patch. The next time the iteration loop runs, it may reintroduce the bug because nothing in the specs or plans prevents it.

With backpropagation:
- Bug fixes become **spec improvements** that persist across all future iterations
- The iteration loop becomes **self-correcting** -- it learns from every manual intervention
- Specs become **progressively more complete** over time
- The gap between "what specs describe" and "what works" **shrinks monotonically**

```
Without backpropagation:
  Bug found -> Fix code -> Bug may return next iteration

With backpropagation:
  Bug found -> Fix code -> Update spec -> Re-run iteration loop -> Fix emerges from specs alone
```

---

## 2. The 8-Step Backpropagation Process

This is the complete process for tracing a bug back to its spec-level root cause.

### Step 1: Discover the Issue

Find a bug, gap, or unexpected behavior in the running application. This can come from:
- Manual testing
- Automated test failures
- User reports
- Monitoring alerts
- Gap analysis (comparing built vs intended)

### Step 2: Fix It (Normal Debugging)

Fix the issue through a normal debugging session. This produces a working code fix. **Do not stop here** -- this fix is temporary until it is encoded in specs.

```bash
# The fix produces commits that we will analyze
git log --oneline -5
# abc1234 Fix: auth token not refreshing on 401 response
# def5678 Fix: race condition in concurrent data fetch
```

### Step 3: Ask What the Spec Missed

This is the critical step. Ask: **"What could we have changed about the spec or prompt to have caught this earlier?"**

Analyze the fix and determine:
- **WHAT** changed (files, functions, behavior)
- **WHY** it was wrong (what assumption was violated)
- **The RULE** (what invariant should have been specified)
- **The LAYER** (which spec, plan, or prompt should have caught this)

Example analysis:

```markdown
## Backpropagation Analysis: Auth Token Refresh

**WHAT changed:** Added 401 response handler in `src/auth/client.ts`
**WHY:** The auth spec did not mention token refresh behavior on HTTP errors
**RULE:** "When an API call receives a 401 response, the auth module MUST
          attempt token refresh before failing"
**LAYER:** spec-auth.md (missing requirement), plan-auth.md (missing task)
**Spec implications:** Add requirement R7 to spec-auth.md with acceptance criteria
```

### Step 4: Update the Spec (Not the Code)

Add the missing requirement or validation to the appropriate spec file:

```markdown
# In context/specs/spec-auth.md, add:

### R7: Token Refresh on Authentication Failure
**Description:** When any API call receives a 401 Unauthorized response,
the auth module must attempt to refresh the authentication token before
propagating the error to the caller.
**Acceptance Criteria:**
- [ ] 401 response triggers automatic token refresh
- [ ] Refresh failure propagates as authentication error
- [ ] Concurrent 401s are deduplicated (only one refresh attempt)
- [ ] Successful refresh retries the original request
**Dependencies:** R1 (token storage), R3 (API client)
```

### Step 5: Trace the Fix into Context Files

Run back-propagation analysis to trace the manual fix into all relevant context files:

1. **Identify affected plan files:** Which plans govern the changed source paths?
2. **Update plans:** Add tasks or notes reflecting the new requirement.
3. **Update impl tracking:** Record the fix and its spec-level root cause.
4. **Add traced annotations:** Mark updated sections with backpropagation metadata.

```markdown
# In context/plans/plan-auth.md, add:

### T-AUTH-007: Implement token refresh on 401
- **Status:** DONE (backpropagated from manual fix abc1234)
- **Spec:** R7 in spec-auth.md
- **Files:** src/auth/client.ts
- **Acceptance criteria:**
  - [ ] 401 triggers refresh
  - [ ] Concurrent dedup
  - [ ] Retry on success
```

### Step 6: Systemic Prompt Changes (If Applicable)

If the issue represents a **pattern** (not a one-off), make systemic changes to the prompts:

**Pattern indicators:**
- Same class of bug appears in multiple domains
- The spec gap is structural (e.g., no specs ever mention error handling)
- The issue reflects a missing validation gate

**Example systemic fix:**

```markdown
# In prompt 003, add to the validation section:

## Error Handling Validation
For every API integration, verify:
- [ ] Error responses (4xx, 5xx) are handled explicitly
- [ ] Authentication failures trigger token refresh
- [ ] Network errors have retry logic with backoff
- [ ] All error paths have corresponding test coverage
```

### Step 7: Re-Run the Iteration Loop

This is the **proof step**. Run the iteration loop against the updated specs and verify that the fix emerges from specs alone, without the manual code fix:

```bash
# Option A: Reset the code fix and re-run
git stash  # temporarily remove the manual fix
iteration-loop context/prompts/003-generate-impl-from-plans.md -n 5 -t 1h
# Verify the fix appears in the implementation

# Option B: Run on a clean branch
git checkout -b test-backprop
git reset --hard <commit-before-fix>
iteration-loop context/prompts/003-generate-impl-from-plans.md -n 5 -t 1h
```

If the fix does NOT emerge from specs alone, the spec update is insufficient. Go back to Step 4 and strengthen the spec.

### Step 8: Generate Regression Tests

Create regression tests that will catch this specific issue if it ever recurs:

```bash
# Generate tests targeting the updated spec
{TEST_COMMAND} --spec context/specs/spec-auth.md

# Or manually create a regression test
# tests/auth/token-refresh-on-401.test.ts
```

The regression test should:
- Directly test the acceptance criteria from Step 4
- Fail if the fix is reverted
- Be included in the standard test suite

---

## 3. Back-Propagate Analysis (Automated)

The back-propagate analysis automates Steps 3-5 by examining recent git history.

### 3.1 Classify Commits

Analyze recent commits and classify each as:

| Classification | Meaning | Action |
|---------------|---------|--------|
| **Manual fix** | Human or interactive agent fixed a bug | Trace back to spec -- this is a backpropagation target |
| **Iteration loop** | Automated iteration loop made the change | No action -- this is the system working as intended |
| **Infrastructure** | Build config, CI, tooling changes | No action -- not spec-related |

**How to classify:**
- Commits from iteration loop sessions have predictable patterns (automated commit messages, batch changes)
- Manual fixes are typically single-issue, focused commits with descriptive messages
- Infrastructure changes touch config files, build scripts, CI pipelines

### 3.2 Analyze Each Manual Fix

For each commit classified as a manual fix, determine:

```markdown
## Commit: abc1234 "Fix: auth token not refreshing on 401"

### WHAT changed
- File: src/auth/client.ts
- Function: handleApiResponse()
- Behavior: Added 401 detection and token refresh logic

### WHY it was wrong
- The auth module did not handle 401 responses
- Tokens would expire and never refresh, causing cascading auth failures

### RULE (invariant that should have been specified)
- "Authentication tokens must be refreshed automatically on 401 responses"

### LAYER (which context file should have caught this)
- spec-auth.md: Missing requirement for error-based token refresh
- plan-auth.md: No task for 401 handling

### Spec Implications
- Add R7 to spec-auth.md: Token Refresh on Authentication Failure
- Add T-AUTH-007 to plan-auth.md: Implement token refresh on 401
```

### 3.3 Discover Affected Plan Files

Dynamically discover which plan files govern the changed source paths:

```
Changed file: src/auth/client.ts
  -> Matches pattern: src/auth/*
  -> Governed by: plan-auth.md
  -> Spec: spec-auth.md

Changed file: src/data/api.ts
  -> Matches pattern: src/data/*
  -> Governed by: plan-data.md
  -> Spec: spec-data.md
```

Use file ownership tables (from prompts) or directory conventions to map source files to plan/spec files.

### 3.4 Update Context Files

For each backpropagation target, update:

1. **Spec file:** Add missing requirement with acceptance criteria
2. **Plan file:** Add task referencing the new requirement
3. **Impl tracking:** Record the backpropagation event

```markdown
# In context/impl/impl-auth.md, add:

## Backpropagation Log
| Date | Commit | Issue | Spec Update | Plan Update |
|------|--------|-------|-------------|-------------|
| 2026-03-14 | abc1234 | 401 not handled | R7 added to spec-auth.md | T-AUTH-007 added |
```

### 3.5 Run Tests

After updating context files, run the test suite to verify nothing broke:

```bash
{BUILD_COMMAND}
{TEST_COMMAND}
```

### 3.6 Generate Regression Tests

For each backpropagation target, generate a regression test that:
- Tests the specific acceptance criteria from the new spec requirement
- Would fail if the fix were reverted
- Is included in the standard test suite going forward

---

## 4. Patterns and Anti-Patterns

### Healthy backpropagation patterns

| Pattern | What It Looks Like |
|---------|--------------------|
| **Decreasing frequency** | Fewer manual fixes needed over time as specs become more complete |
| **Spec coverage expanding** | Each backpropagation adds requirements that prevent entire classes of bugs |
| **Systemic improvements** | Prompt-level changes prevent bugs across all domains, not just the one that triggered it |
| **Iteration loop reproducing fixes** | After spec update, the iteration loop generates the same fix autonomously |

### Anti-patterns to watch for

| Anti-Pattern | Symptom | Fix |
|-------------|---------|-----|
| **Fixing code without updating specs** | Same class of bug keeps recurring | Enforce the 8-step process; never skip Step 4 |
| **Spec updates that are too narrow** | Each fix only prevents the exact same bug, not related ones | Think about the RULE, not just the fix -- generalize |
| **Skipping the proof step** | Updated specs but did not verify the iteration loop reproduces the fix | Always run Step 7; if it fails, the spec is insufficient |
| **Over-specifying** | Specs become so detailed they are brittle and break on minor changes | Specs should describe WHAT and WHY, not exact HOW |
| **Backpropagation debt** | Many manual fixes accumulate without being traced back to specs | Schedule regular backpropagation sessions; do not let debt grow |

---

## 5. When NOT to Backpropagate

Not every code fix needs backpropagation:

- **One-off environment issues** (wrong config, missing dependency) -- these are infrastructure, not spec gaps
- **Typos and formatting** -- trivial fixes that do not reflect missing requirements
- **Exploratory changes** during prototyping -- specs are still being formed
- **Performance optimizations** that do not change behavior -- unless performance is a spec requirement

**Rule of thumb:** If the iteration loop could plausibly reintroduce the bug, backpropagate. If not, skip it.

---

## 6. Backpropagation and Convergence

Backpropagation directly improves convergence:

```
Iteration 1: 200 lines changed, 5 manual fixes needed
  -> Backpropagate all 5 fixes into specs
Iteration 2: 100 lines changed, 2 manual fixes needed
  -> Backpropagate 2 fixes
Iteration 3: 50 lines changed, 0 manual fixes needed
  -> Convergence achieved
```

Each backpropagation cycle makes the specs more complete, which means the iteration loop needs fewer iterations to reach a stable solution. If you are NOT seeing convergence, check whether you are consistently backpropagating manual fixes.

**Non-convergence + frequent manual fixes = backpropagation debt.** The specs are not capturing what the code needs, so the iteration loop keeps producing broken output that requires manual intervention.

---

## 7. Integration with Other SDD Skills

- **Convergence monitoring:** Use `sdd:convergence-monitoring` to detect when manual fixes are decreasing (good) or increasing (backpropagation debt).
- **Prompt pipeline:** Backpropagation may trigger changes to prompts (Step 6), which affects the `sdd:prompt-pipeline` design.
- **Validation-first design:** Stronger validation gates catch issues earlier, reducing the need for backpropagation.
- **Gap analysis:** Systematic gap analysis (`/sdd:gap-analysis`) identifies backpropagation targets proactively, rather than waiting for bugs.

---

## Cross-References

- **Convergence patterns:** See `references/convergence-patterns.md` for how backpropagation drives convergence.
- **Prompt pipeline:** See `sdd:prompt-pipeline` skill for how prompt 006 (rewrite pattern) implements automated backpropagation.
- **Impl tracking:** See `sdd:impl-tracking` skill for the backpropagation log format in implementation tracking documents.
- **Validation gates:** See `sdd:validation-first` skill for validation layers that catch issues before they require backpropagation.
