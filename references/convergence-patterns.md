# Convergence Patterns Reference

Detecting whether agent iterations are converging toward a solution or hitting a ceiling. Covers the exponential decay curve, convergence vs ceiling distinction, non-convergence diagnosis, test pass rate as signal, and forward progress metrics.

---

## 1. Overview

Convergence is the central feedback signal in SDD. It tells you whether the iteration loop is working -- whether the agent is approaching a stable solution that satisfies all specifications and validation gates.

> **Convergence means the agent's output is stabilizing. The signal is not zero changes -- it is that remaining changes are trivial.**

Understanding convergence is critical because it answers the most important question in any iterative AI workflow: **when are we done?**

---

## 2. The Exponential Decay Curve

Healthy convergence follows an exponential decay pattern in the volume of changes per iteration:

```
Changes (lines)
    ^
500 |*
    |
400 |
    |  *
300 |
    |
200 |     *
    |
100 |        *
    |
 50 |           *
 20 |              * * *    <-- convergence zone
  0 +--|--|--|--|--|--|--|-->  Iteration
    1  2  3  4  5  6  7  8
```

### Typical Progression

| Iteration | Lines Changed | Character |
|-----------|--------------|-----------|
| 1 | 400-500 | Major scaffolding, initial implementation |
| 2 | 200-300 | Filling gaps, fixing broken tests |
| 3 | 100-150 | Edge cases, error handling |
| 4 | 50-80 | Polish, minor fixes |
| 5 | 20-40 | Trivial adjustments |
| 6 | 10-20 | Convergence zone -- nearly stable |
| 7 | 5-15 | Convergence zone -- essentially done |
| 8 | 0-10 | Converged |

### What "Convergence Zone" Means

The convergence zone is reached when:
- Changes are small (< ~20 lines per iteration)
- Changes are trivial (formatting, comments, minor adjustments)
- No new functionality is being added
- No new tests are being added
- Existing tests remain passing

**You do NOT need zero changes.** Perfect zero is rare and unnecessary. The signal is that changes are trivial and decreasing.

---

## 3. Convergence vs. Ceiling

Both convergence and ceiling produce small diffs, but they have fundamentally different implications.

### Convergence

**Definition:** The agent has successfully implemented all requirements and validation passes. Remaining changes are trivial refinements.

**Indicators:**
- All tests pass
- Build is clean
- Implementation tracking shows all tasks DONE
- Changes are cosmetic (formatting, comments, minor naming)
- Agent emits completion signal
- Exit criteria checklist is fully checked

**Action:** You are done. Move to the next phase or stop iterating.

### Ceiling

**Definition:** The agent cannot make further progress due to an external constraint. It is stuck, not done.

**Indicators:**
- Some tests still fail (same ones each iteration)
- Build may have persistent warnings or errors
- Implementation tracking shows tasks stuck at IN_PROGRESS or BLOCKED
- Changes are repetitive (agent tries the same thing each iteration)
- Agent does NOT emit completion signal
- Exit criteria checklist has unchecked items
- Agent may explicitly report being blocked

**Action:** Investigate the blocker. The agent needs help -- a spec change, a dependency fix, a tooling update, or human guidance.

### Comparison Table

| Characteristic | Convergence | Ceiling |
|---------------|-------------|---------|
| Diff size | Small, decreasing | Small, flat |
| Test pass rate | 100% or approaching | Stuck below 100% |
| Task status | All DONE | Some BLOCKED |
| Change character | Trivial refinements | Repetitive attempts |
| Completion signal | Emitted | Not emitted |
| Root cause | Success | External constraint |
| Action | Stop iterating | Investigate blocker |

---

## 4. Detecting Convergence

### Method 1: Change Volume Tracking

Track the total lines changed per iteration:

```bash
# After each iteration, measure changes
git diff HEAD~1 --stat | tail -1
# Output: "15 files changed, 450 insertions(+), 120 deletions(-)"
```

Plot or log these numbers. Convergence = monotonically decreasing volume.

### Method 2: Test Pass Rate

Track the test pass rate across iterations:

| Iteration | Tests | Pass | Fail | Skip | Pass Rate |
|-----------|-------|------|------|------|-----------|
| 1 | 50 | 30 | 15 | 5 | 60% |
| 2 | 75 | 55 | 12 | 8 | 73% |
| 3 | 90 | 78 | 8 | 4 | 87% |
| 4 | 95 | 90 | 3 | 2 | 95% |
| 5 | 98 | 96 | 1 | 1 | 98% |
| 6 | 100 | 99 | 0 | 1 | 99% |

Healthy convergence shows:
- Test count increasing (more coverage)
- Pass rate approaching 100%
- Fewer failures per iteration
- Skip count decreasing

### Method 3: Task Completion Rate

Track completion of spec requirements:

```markdown
Iteration 1: 3/20 requirements satisfied (15%)
Iteration 2: 8/20 requirements satisfied (40%)
Iteration 3: 14/20 requirements satisfied (70%)
Iteration 4: 18/20 requirements satisfied (90%)
Iteration 5: 19/20 requirements satisfied (95%)
Iteration 6: 20/20 requirements satisfied (100%)
```

### Method 4: Implementation Tracking

Read the implementation tracking document for task status progression:

```markdown
After iteration 1: 3 DONE, 12 IN_PROGRESS, 5 TODO
After iteration 2: 8 DONE, 9 IN_PROGRESS, 3 TODO
After iteration 3: 14 DONE, 5 IN_PROGRESS, 1 TODO
After iteration 4: 18 DONE, 2 IN_PROGRESS, 0 TODO
After iteration 5: 19 DONE, 1 IN_PROGRESS, 0 TODO
After iteration 6: 20 DONE, 0 IN_PROGRESS, 0 TODO
```

---

## 5. Non-Convergence Diagnosis

When changes are NOT decreasing iteration-over-iteration, something is wrong. The iteration loop is not the problem -- the inputs to the loop are.

### Cause 1: Fuzzy Specs

**Symptom:** Agent makes different changes each iteration. Code oscillates between approaches.

**Diagnosis:** The specification is ambiguous. The agent interprets requirements differently each pass because there is no single clear interpretation.

**Examples:**
- "The UI should be responsive" -- How responsive? What breakpoints? What elements?
- "Handle errors gracefully" -- What errors? What does graceful mean? Log? Retry? Show UI?
- "Good performance" -- What metric? What threshold? Under what load?

**Fix:** Make specs specific and testable. Replace vague requirements with measurable acceptance criteria:

```markdown
# Bad
R1: The API should be fast

# Good
R1: API Response Time
**Acceptance Criteria:**
- [ ] GET /api/items responds in < 200ms for up to 1000 items
- [ ] POST /api/items responds in < 100ms
- [ ] No endpoint exceeds 500ms under normal load
```

### Cause 2: Weak Validation

**Symptom:** Agent makes changes but does not know if they are correct. Changes do not systematically address failures.

**Diagnosis:** The validation gates are too weak. The agent has no signal about whether its changes improve or regress the system.

**Examples:**
- No tests for a requirement (agent cannot verify it works)
- Tests are too broad (pass even when behavior is wrong)
- Build succeeds but application does not actually work
- No launch verification (agent never checks if the app starts)

**Fix:** Strengthen validation gates. Every spec requirement must have at least one automated validation:

```markdown
# Each requirement maps to a specific test
R1: User login
  -> Unit test: login function returns session token
  -> Integration test: POST /login with valid credentials returns 200
  -> E2E test: User can log in via the UI
  -> Launch test: Login page renders after app startup
```

### Cause 3: Fighting Sub-Agents

**Symptom:** Changes oscillate. One iteration adds code, the next removes it. Lines changed stays flat or increases.

**Diagnosis:** Multiple agents (or agent team teammates) are making conflicting changes. One agent's output is being undone by another.

**Examples:**
- Teammate A refactors a shared file, Teammate B reverts the refactor
- Iteration N adds a feature, iteration N+1 removes it because it breaks something else
- Two tasks have contradictory requirements in the specs

**Fix:**
1. Check file ownership tables -- ensure no shared file has multiple writers
2. Check specs for contradictions -- resolve any conflicting requirements
3. Check task dependencies -- ensure tasks are ordered correctly
4. If using agent teams, ensure teammates work in isolated worktrees

### Cause 4: Missing Dependencies

**Symptom:** Agent repeatedly tries the same approach and fails each time.

**Diagnosis:** The agent needs something that does not exist yet -- a library, a service, a configuration, a tool.

**Examples:**
- Tests require a database that is not configured
- Implementation needs an API that is not deployed
- Build requires a dependency that is not installed

**Fix:** Identify the missing dependency and either:
- Install/configure it
- Mock it for the agent
- Update specs to exclude it from current scope
- Document it as a blocker in implementation tracking

---

## 6. Non-Convergence is a Signal to Improve Specs

This is a critical principle:

> **Non-convergence is a signal to improve your specs, not to run more iterations.**

Running more iterations when the system is not converging is waste. The agent will keep making the same mistakes because the underlying problem is in the specs, validation, or coordination -- not in the number of passes.

### Decision Tree

```
Are changes decreasing iteration-over-iteration?
  |
  YES -> Continue iterating. Convergence is working.
  |
  NO  -> Stop. Diagnose the cause:
         |
         +-> Are specs ambiguous? -> Fix specs, restart
         |
         +-> Is validation weak? -> Add tests/gates, restart
         |
         +-> Are agents fighting? -> Fix ownership/deps, restart
         |
         +-> Is there a missing dep? -> Fix or mock it, restart
```

---

## 7. Test Pass Rate as First-Class Convergence Signal

Test pass rate is the most reliable convergence signal because it is:
- **Objective:** Tests either pass or fail
- **Quantitative:** Easily tracked as a number
- **Incremental:** Shows progress between 0% and 100%
- **Granular:** Individual test failures point to specific issues

### Tracking Test Health

```markdown
## Test Health Dashboard

| Iteration | Total | Pass | Fail | Skip | New Tests | New Passes |
|-----------|-------|------|------|------|-----------|------------|
| 1 | 50 | 30 | 15 | 5 | 50 | 30 |
| 2 | 75 | 55 | 12 | 8 | 25 | 25 |
| 3 | 90 | 78 | 8 | 4 | 15 | 23 |
| 4 | 95 | 90 | 3 | 2 | 5 | 12 |
| 5 | 98 | 96 | 1 | 1 | 3 | 6 |
```

### Healthy Signals

- Total test count increasing (agent is adding coverage)
- Pass count increasing faster than total count (fixing existing failures AND adding new passing tests)
- Fail count decreasing monotonically
- Skip count decreasing (fewer known issues being deferred)

### Unhealthy Signals

- Pass rate oscillating (tests that passed now fail, suggesting regression)
- Total test count flat (agent is not adding coverage)
- Same tests failing every iteration (ceiling -- agent cannot fix them)
- Test count decreasing (agent is deleting tests instead of fixing them)

---

## 8. Forward Progress Metrics

For large projects with many spec requirements, track the percentage of requirements with passing tests.

### Forward Progress Calculation

```
Forward Progress = (Requirements with ALL acceptance criteria passing)
                   / (Total requirements)
                   * 100%
```

### Example Tracking

```markdown
## Forward Progress

| Domain | Total Reqs | Passing | Progress |
|--------|-----------|---------|----------|
| Auth | 8 | 7 | 87.5% |
| API | 12 | 10 | 83.3% |
| UI | 15 | 9 | 60.0% |
| Data | 6 | 6 | 100% |
| Perf | 4 | 2 | 50.0% |
| **Total** | **45** | **34** | **75.6%** |
```

### Using Forward Progress

- Track per domain to identify which areas need more iteration
- 100% in a domain means that domain has converged
- Domains below 50% may have spec or validation issues
- Overall forward progress should increase monotonically
- Plateaus indicate ceiling -- investigate the blocked domains

---

## 9. Convergence Monitoring in Practice

### Manual Monitoring

After each iteration:
1. Check `git diff --stat` for change volume
2. Run `{TEST_COMMAND}` and record pass/fail counts
3. Read implementation tracking for task status
4. Log the results in a convergence tracking file

### Automated Monitoring

A convergence monitor agent can automate this:

```markdown
## Convergence Monitor

Every N minutes (or after each iteration):
1. Read git log for recent commits
2. Calculate change volume (insertions + deletions)
3. Run test suite and record results
4. Read implementation tracking
5. Report:
   - Change trend (decreasing/flat/increasing)
   - Test trend (pass rate direction)
   - Task trend (completion direction)
   - Diagnosis (converging/ceiling/non-converging)
6. If non-converging: alert human with diagnosis
```

---

## 10. Iteration Count Guidelines

While convergence is the true "done" signal, practical iteration limits prevent runaway loops:

| Project Size | Typical Iterations to Convergence | Max Iterations |
|-------------|----------------------------------|----------------|
| Small (5-10 files) | 3-5 | 10 |
| Medium (10-50 files) | 5-8 | 15 |
| Large (50-200 files) | 8-15 | 25 |
| Very large (200+ files) | 10-20+ | 30 |

### When Max Is Reached Without Convergence

If the iteration limit is reached and the system has not converged:
1. Stop iterating
2. Analyze the convergence metrics
3. Diagnose the cause (fuzzy specs, weak validation, fighting agents, missing deps)
4. Fix the root cause
5. Reset the iteration counter
6. Resume iterating

Never increase the iteration limit to compensate for non-convergence.

---

## 11. Convergence in Multi-Agent Context

When using agent teams, convergence has additional dimensions:

### Per-Teammate Convergence

Each teammate should converge independently within their domain:
- Changes to their files should decrease per iteration
- Tests in their domain should approach 100% pass rate
- Their tasks should progress toward DONE

### Cross-Team Convergence

After merging all teammates' work:
- The combined build should pass
- The combined test suite should pass
- Integration tests should show increasing pass rate
- Merge conflicts should decrease or disappear

### Team-Level Non-Convergence

If the team as a whole is not converging despite individual teammates converging:
- Check for integration issues at merge boundaries
- Verify file ownership is correct (no conflicts)
- Check for specification contradictions across domains
- Verify task dependencies are correct

---

## 12. Convergence Patterns Summary

| Pattern | Change Volume | Test Pass Rate | Action |
|---------|--------------|----------------|--------|
| Healthy convergence | Decreasing exponentially | Approaching 100% | Continue, almost done |
| Slow convergence | Decreasing linearly | Increasing slowly | More iterations needed, but working |
| Ceiling | Flat (small) | Stuck below 100% | Investigate blocker |
| Non-convergence | Flat or increasing | Oscillating or flat | Stop, fix specs/validation |
| Divergence | Increasing | Decreasing | Stop immediately, major issue |

### The Cardinal Rules

1. **Convergence tells you when you are done.** Not iteration count, not time -- convergence.
2. **Non-convergence is a signal to improve specs, not to run more iterations.**
3. **Test pass rate is the most reliable convergence signal.**
4. **Both convergence and ceiling produce small diffs -- learn to distinguish them.**
5. **For large projects, track forward progress as percentage of spec requirements with passing tests.**
