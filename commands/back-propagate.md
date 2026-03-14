---
name: sdd-back-propagate
description: Trace recent manual code fixes back into specs and context files
---

# SDD Back-Propagate — Trace Fixes to Specs

You are performing backpropagation: tracing recent manual code changes back into specs, plans, and context files so that the convergence loop can reproduce them autonomously. The core principle: if a bug can only be fixed in code, your specs are incomplete.

## Step 1: Analyze Recent Commits

Run `git log --oneline -20` (or since the last convergence loop run) to gather recent commits. Read the diffs for each commit.

Classify each commit into one of three categories:

| Category | Description | Action |
|----------|-------------|--------|
| **Manual fix** | Human-authored code change fixing a bug or adding behavior | Back-propagate (proceed to Step 2) |
| **Iteration loop** | Changes made by an automated convergence loop session | Skip — already spec-driven |
| **Infrastructure** | Build config, CI, tooling, dependency updates | Skip — not spec-relevant |

Report the classification table to the user.

## Step 2: Analyze Each Manual Fix

For each commit classified as a manual fix, determine:

1. **WHAT** changed — which files, which functions, what behavior was added/modified/removed
2. **WHY** it changed — from the commit message, PR description, diff context, and surrounding code comments
3. **RULE** — what invariant or requirement was violated that necessitated this fix
4. **LAYER** — which spec, plan, or prompt should have caught this:
   - **Spec gap**: the requirement was never specified
   - **Plan gap**: the spec existed but the plan didn't implement it
   - **Prompt gap**: the plan existed but the prompt didn't guide the agent to it
   - **Validation gap**: everything existed but no test caught the regression

## Step 3: Discover Governing Plan Files

For each changed source file, determine which plan file governs it:
- Search `context/plans/` for plan files that reference the changed paths
- If no plan covers the file, flag it as an **untracked file** (potential spec gap)

## Step 4: Update Context Files

For each manual fix, update the appropriate context files:

### Spec Updates (context/specs/)
- If the fix reveals a missing requirement, add it with testable acceptance criteria
- If the fix reveals an ambiguous requirement, clarify it
- Add a `## Back-propagated` section or annotation noting the source commit

### Plan Updates (context/plans/)
- If the plan missed a task, add it with proper T- prefix and dependencies
- If the plan had incorrect sequencing, fix the dependency graph
- Update plan-known-issues.md if the fix reveals a systemic issue

### Impl Tracking Updates (context/impl/)
- Record the manual fix in the relevant impl tracking file
- Add to the "Dead Ends & Failed Approaches" section if the fix replaced a failed approach
- Update test health if new tests were added

## Step 5: Run Tests

Run the project's test suite to verify that:
- The manual fixes still pass
- No regressions were introduced by context file updates
- Any new acceptance criteria added in Step 4 have corresponding tests

If the build command is not obvious, ask the user for `{BUILD_COMMAND}` and `{TEST_COMMAND}`.

## Step 6: Report

Generate a summary report:

```markdown
## Back-Propagation Report

### Commits Analyzed
| Commit | Category | Action |
|--------|----------|--------|

### Manual Fixes Traced
| Commit | WHAT | WHY | RULE | LAYER | Files Updated |
|--------|------|-----|------|-------|---------------|

### Context Files Updated
- specs: {list of updated spec files}
- plans: {list of updated plan files}
- impl: {list of updated impl files}

### Test Results
- Pass: {count}
- Fail: {count}
- New tests added: {count}

### Recommendations
- {Any systemic prompt changes suggested}
- {Any specs that need deeper review}
```

Present this report to the user.
