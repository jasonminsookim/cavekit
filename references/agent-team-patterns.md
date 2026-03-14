# Agent Team Patterns Reference

Coordination patterns for multi-agent development. Covers team structure, batching, delegation, worktree isolation, file ownership, merge protocol, and shutdown procedures.

---

## 1. Overview

Agent teams decompose work into parallel sub-tasks, each handled by a separate AI agent instance with its own context window. This enables parallelism, context isolation, and coordinated development of complex systems that exceed what a single context window can hold.

**When to Use Agent Teams:**
- Work spans multiple domains that can be developed in parallel
- Total context (specs + plans + code) exceeds a single context window
- Tasks have limited interdependencies
- You need to finish faster through parallelism

**When NOT to Use Agent Teams:**
- Work is small enough for a single agent
- Tasks are highly interdependent (serial is better)
- The codebase has extensive shared state that cannot be cleanly partitioned

---

## 2. Team Structure

### Roles

| Role | Responsibility | Code Writing |
|------|---------------|-------------|
| **Team Lead** | Orchestrate, coordinate, summarize | Never -- delegate mode |
| **Teammate** | Implement assigned domain | Yes, within owned files |
| **Sub-agent** | Discrete subtask for teammate | Limited, reports to teammate |

### The Lead Never Writes Code

This is a critical rule. The team lead operates in delegate mode -- it:
- Creates and assigns tasks
- Spawns teammates with complete context
- Monitors progress
- Coordinates merge order
- Resolves conflicts between teammates
- Summarizes results

The lead does NOT:
- Write code directly
- Modify source files
- Run build/test commands on production code
- Make implementation decisions without delegating

**Why:** Delegate mode forces proper decomposition. If the lead wrote code, it would accumulate implementation details in its context window, reducing its ability to coordinate. The lead's context window is reserved for orchestration.

---

## 3. Max 3 Concurrent Teammates

### The Rule

Never spawn more than 3 teammates simultaneously. This limit exists because:

1. **Resource exhaustion:** Each teammate consumes CPU, memory, and API rate limits
2. **Terminal multiplexer race conditions:** More than 3 parallel sessions causes coordination issues
3. **Diminishing returns:** Coordination overhead increases faster than parallelism benefit
4. **Merge complexity:** More concurrent teammates means more potential conflicts

### The Pattern

```
Batch 1: Spawn A, B, C  ->  Wait for completion  ->  Shutdown  ->  Merge
Batch 2: Spawn D, E, F  ->  Wait for completion  ->  Shutdown  ->  Merge
Batch 3: Spawn G         ->  Wait for completion  ->  Shutdown  ->  Merge
```

### Sub-Agent Limit

Each teammate can spawn up to 3 sub-agents. This gives a maximum hierarchy of:

```
Lead (1)
+-- Teammate A (1) + sub-agents (up to 3)
+-- Teammate B (1) + sub-agents (up to 3)
+-- Teammate C (1) + sub-agents (up to 3)
```

Maximum simultaneous agents: 1 lead + 3 teammates + 9 sub-agents = 13

In practice, sub-agents are short-lived and rarely all active simultaneously.

---

## 4. Batch Phases

Work is divided into batches. Each batch goes through a complete lifecycle.

### Batch Lifecycle

```
1. PLAN
   - Lead identifies tasks for this batch
   - Lead assigns tasks to teammates
   - Lead creates file ownership table
   - Lead prepares spawn prompts

2. SPAWN
   - Lead creates worktrees for each teammate
   - Lead spawns teammates with complete context
   - Maximum 3 teammates per batch

3. EXECUTE
   - Teammates work in parallel
   - Each teammate works only in their worktree
   - Each teammate modifies only files they own
   - Teammates commit frequently
   - Teammates update implementation tracking

4. COMPLETE
   - Teammates report completion to lead
   - Lead verifies all teammates are done

5. SHUTDOWN
   - All teammates in the batch are shut down
   - Their context windows are released

6. MERGE
   - Lead merges one teammate at a time
   - After each merge: build -> test -> verify
   - Fix any issues before merging next teammate
   - Clean up worktrees and branches

7. TRANSITION
   - Lead assesses remaining work
   - Lead plans next batch (if needed)
   - Repeat from step 1
```

### Why Shutdown Between Batches

Teammates must be shut down between batches because:
- Their context windows are stale after merge (code has changed)
- Resource cleanup prevents accumulation
- Clean slate for next batch prevents confusion
- Merge results may change the task plan for the next batch

---

## 5. Worktree Isolation

Every teammate works in an isolated git worktree. This prevents file conflicts entirely at the filesystem level.

### Worktree Creation

```bash
# Lead creates worktree for each teammate
git worktree add ./worktrees/{teammate-name} -b feat/impl/{teammate-name}
```

### Worktree Layout

```
project-root/
+-- src/                    # Main branch (lead's view)
+-- context/
+-- worktrees/
    +-- auth/               # Teammate A's isolated copy
    |   +-- src/
    |   +-- context/
    +-- api/                # Teammate B's isolated copy
    |   +-- src/
    |   +-- context/
    +-- ui/                 # Teammate C's isolated copy
        +-- src/
        +-- context/
```

### Rules for Teammates in Worktrees

1. Work ONLY within your assigned worktree directory
2. Do NOT access other teammates' worktrees
3. Commit to your branch frequently
4. Do NOT push to remote
5. Do NOT merge your branch into main
6. Read specs/plans from the worktree copy (they are shared via git)

### Worktree Cleanup

```bash
# After merging teammate's work
git worktree remove ./worktrees/{teammate-name}
git branch -d feat/impl/{teammate-name}
```

---

## 6. Explicit File Ownership

File ownership eliminates merge conflicts by ensuring exactly one teammate can modify any given file.

### Ownership Table Format

```markdown
| File/Directory | Owner | Notes |
|---------------|-------|-------|
| `src/auth/*` | auth-teammate | Authentication module |
| `src/api/*` | api-teammate | API endpoints |
| `src/ui/*` | ui-teammate | UI components |
| `src/shared/types.ts` | auth-teammate | Shared type definitions |
| `src/shared/config.ts` | api-teammate | Configuration |
| `src/shared/constants.ts` | ui-teammate | UI constants |
| `package.json` | LEAD ONLY | Modified during merge phase |
| `tsconfig.json` | LEAD ONLY | Modified during merge phase |
```

### Ownership Rules

1. **Every modifiable file has exactly one owner**
2. **Non-owners may READ any file** but must NOT modify
3. **Shared configuration files** (package.json, build configs) are owned by the lead and modified only during merge phases
4. **If a teammate needs a change in a file they do not own**, they must:
   - Document the needed change in their implementation tracking
   - The lead coordinates the change during merge phase
5. **Test files follow source ownership** -- if you own `src/auth/*`, you own `tests/auth/*`

### Cross-Plan Ownership

When work spans multiple plans, the ownership table must map files to plans:

```markdown
| File | Plan Owner | Teammate |
|------|-----------|----------|
| `src/rpc/types.ts` | plan-auth | auth-teammate |
| `src/rpc/client.ts` | plan-api | api-teammate |
| `src/rpc/server.ts` | plan-api | api-teammate |
```

---

## 7. Merge Protocol

Merging is the most critical phase of agent team coordination. Errors here can undo all parallel work.

### The Protocol

```
For each teammate (one at a time, in dependency order):

1. MERGE
   git checkout main
   git merge feat/impl/{teammate-name} --no-ff

2. BUILD
   {BUILD_COMMAND}
   If fail -> fix -> rebuild -> re-verify

3. TEST
   {TEST_COMMAND}
   If fail -> fix -> retest -> re-verify

4. VERIFY
   Verify the merged feature works as expected
   Check implementation tracking for expected behavior

5. CLEAN
   git worktree remove ./worktrees/{teammate-name}
   git branch -d feat/impl/{teammate-name}

6. NEXT
   Only proceed to next teammate after ALL gates pass
```

### Merge Order

Merge in dependency order:
1. Foundation/infrastructure teammates first (shared types, config)
2. Core domain teammates next (business logic, data models)
3. Dependent teammates last (UI, integration layers)

### Conflict Resolution

If a merge conflict occurs:
1. The lead resolves the conflict (it owns the merge process)
2. After resolution: build -> test -> verify before proceeding
3. Document the conflict and resolution in implementation tracking
4. Consider whether the file ownership table needs updating

### Post-Merge Validation Checklist

```markdown
- [ ] `{BUILD_COMMAND}` passes with zero errors
- [ ] `{TEST_COMMAND}` passes with zero failures
- [ ] No regressions in previously passing tests
- [ ] Implementation tracking updated with merge results
- [ ] Worktree removed
- [ ] Branch deleted
```

---

## 8. Shutdown Protocol

Shutdown happens between batch phases and at session end.

### Between Batches

1. Wait for all teammates in the batch to report completion
2. Collect status reports from each teammate
3. Shut down all teammate sessions
4. Proceed to merge phase

### Graceful Shutdown

When shutting down a teammate:
1. Teammate commits all work in progress
2. Teammate updates implementation tracking with final status
3. Teammate reports completion to lead
4. Lead terminates the teammate session

### Emergency Shutdown

If a teammate is stuck or producing bad output:
1. Lead sends stop signal
2. Teammate commits whatever work is salvageable
3. Lead terminates the session
4. Lead documents the failure in implementation tracking
5. Work may be re-assigned to a new teammate in the next batch

---

## 9. Communication Patterns

### Lead -> Teammate (Spawn)

The spawn prompt is the only communication from lead to teammate at creation time. It must be self-contained.

### Teammate -> Lead (Report)

Teammates report back to the lead:
- When they complete all assigned tasks
- When they hit a blocker they cannot resolve
- When they need a change in a file they do not own

### Teammate -> Teammate

Teammates do NOT communicate directly. All coordination goes through the lead. This prevents:
- Race conditions in communication
- Context pollution between teammates
- Uncoordinated changes

### Lead -> Sub-agent (via Teammate)

Teammates can spawn sub-agents for discrete tasks. Sub-agents report to their spawning teammate, not to the lead.

---

## 10. Team Sizing Guidelines

### Small Project (3-5 domains)

```
Lead
+-- Teammate A: domain-1 + domain-2
+-- Teammate B: domain-3
+-- Teammate C: domain-4 + domain-5
```

One batch, 3 teammates, each handling 1-2 domains.

### Medium Project (6-12 domains)

```
Batch 1:
  Lead
  +-- Teammate A: domain-1, domain-2
  +-- Teammate B: domain-3, domain-4
  +-- Teammate C: domain-5, domain-6

Batch 2:
  Lead
  +-- Teammate D: domain-7, domain-8
  +-- Teammate E: domain-9, domain-10
  +-- Teammate F: domain-11, domain-12
```

Two batches, 3 teammates each.

### Large Project (12+ domains)

Multiple batches with dependency-ordered execution:

```
Batch 1 (Foundation): Shared types, config, core models
Batch 2 (Core): Business logic, data access, API
Batch 3 (Features): UI, integrations, advanced features
Batch 4 (Polish): Performance, edge cases, documentation
```

---

## 11. Common Failure Modes

### Teammate Modifies Unowned File

**Symptom:** Merge conflict during merge phase.
**Prevention:** Strict file ownership table in spawn prompt.
**Recovery:** Resolve conflict, update ownership table, re-spawn if needed.

### Lead Writes Code Directly

**Symptom:** Lead's context window fills with implementation details.
**Prevention:** Delegate mode enforcement in the prompt.
**Recovery:** Move code writing to a teammate in the next batch.

### Too Many Concurrent Teammates

**Symptom:** System slowdown, API rate limits, race conditions.
**Prevention:** Max 3 concurrent teammates, batch phases.
**Recovery:** Shut down excess teammates, re-batch.

### Teammate Ignores Time Guards

**Symptom:** Teammate spends 2 hours on one task.
**Prevention:** Explicit time guards in spawn prompt.
**Recovery:** Lead sends stop signal, re-assigns task.

### Stale Worktree After Merge

**Symptom:** Teammate continues working with outdated code.
**Prevention:** Shutdown between batches, clean worktree creation.
**Recovery:** Shut down stale teammate, create fresh worktree.

---

## 12. Worktree Commands Reference

### Create Worktree

```bash
git worktree add ./worktrees/{name} -b feat/impl/{name}
```

### List Worktrees

```bash
git worktree list
```

### Remove Worktree

```bash
git worktree remove ./worktrees/{name}
```

### Force Remove (if locked)

```bash
git worktree remove --force ./worktrees/{name}
```

### Clean Up Branches

```bash
# Delete merged branch
git branch -d feat/impl/{name}

# Force delete unmerged branch (use with caution)
git branch -D feat/impl/{name}
```

### Verify Clean State

```bash
# After all merges and cleanup
git worktree list          # Should show only main worktree
git branch --list 'feat/*' # Should show no remaining feature branches
```

---

## 13. Integration with Iteration Loops

Agent teams can be used within iteration loops. Each iteration of the loop:

1. Lead reads git state and implementation tracking
2. Lead plans the next batch of work
3. Lead spawns teammates
4. Teammates execute
5. Lead merges
6. Lead updates tracking
7. Lead checks exit criteria
8. If not complete, emit status and wait for next iteration

### Convergence in Team Context

Convergence metrics for teams:
- **Per-teammate:** Lines changed per iteration should decrease
- **Per-merge:** Number of conflicts should decrease
- **Overall:** Total remaining tasks should decrease monotonically
- **Test health:** Test count should increase, failure count should decrease

### When Teams Are Overkill in Later Iterations

As the project converges, the remaining work may be small enough for a single agent. The lead should recognize this and switch from team mode to single-agent mode when:
- Remaining tasks are in a single domain
- Changes are minor fixes or polishing
- Team overhead exceeds the parallelism benefit
