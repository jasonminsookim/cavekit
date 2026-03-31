---
name: bp-build
description: "Implement a build site — automatically parallelizes independent tasks and progresses through tiers autonomously"
argument-hint: "[--filter PATTERN] [--peer-review] [--max-iterations N] [--completion-promise TEXT]"
allowed-tools: ["Bash(${CLAUDE_PLUGIN_ROOT}/scripts/setup-build.sh:*)", "Bash(git *)"]
---

# Blueprint Build — Autonomous Implementation

This is the third phase of Blueprint. Execute the setup script:

```!
"${CLAUDE_PLUGIN_ROOT}/scripts/setup-build.sh" $ARGUMENTS
```

## If site selection is required

If the output contains `BLUEPRINT_SITE_SELECTION_REQUIRED=true`, multiple build sites were found. **Ask the user which one to implement.** Then re-run with `--filter <their-choice>`.

## Execution Loop

Once the setup script completes (outputs the ralph prompt), you run the execution loop autonomously. Progress through all tiers without stopping.

### Each Wave

1. **Read state**: Read the build site + all `context/impl/impl-*.md` files + `context/impl/dead-ends.md`. If this is the first wave of a new tier, capture the tier start ref: `TIER_START_REF=$(git rev-parse HEAD)`
2. **Compute frontier**: Find all tasks that are NOT done AND whose `blockedBy` dependencies are ALL done
3. **Report**:
   ```
   ═══ Wave {N} ═══
   {count} task(s) ready:
     {task_id}: {title} (tier {N}, deps: {deps})
   ```

4. **Execute based on frontier size**:

   **0 ready tasks** → Check if ALL tasks are done. If yes → completion. If not → report blockage and stop.

   **1 ready task** → Implement it directly yourself. Follow the ralph prompt instructions (read blueprint, implement, validate, commit, track).

   **2+ ready tasks** → Partition the frontier into a small set of coherent work packets, then parallelize those. Optimize for throughput, not raw agent count.

   **Build work packets using these rules:**
   - Group tasks when they touch the same subsystem, blueprint/domain, or expected file set
   - Group small related tasks when the combined scope is still coherent for one agent
   - Split large, risky, or file-disjoint tasks into separate packets
   - Keep ownership clean: each packet should have a clear primary file/module surface
   - Prefer a few meaningful packets, usually 2-4 concurrent subagents, rather than one agent per task
   - If a tiny task is faster to do inline than to delegate, keep it yourself and only delegate the heavier or cleaner-isolated packets

   ```
   Agent(
     subagent_type: "bp:task-builder",
     isolation: "worktree",
     prompt: "TASKS:
   - {task_id}: {title}
   - {task_id}: {title}

   SHARED CONTEXT:
   - Domain/spec: {spec_name}
   - Requirement IDs: {requirement_ids}
   BUILD SITE: {path to build site}
   BLUEPRINTS: {paths to relevant blueprint files}
   EXPECTED FILE OWNERSHIP: {files or modules this packet should own}

   ACCEPTANCE CRITERIA (from blueprints):
   {paste the acceptance criteria for each task in this packet}

   DEAD ENDS TO AVOID:
   {paste relevant dead ends, or 'None'}

   INSTRUCTIONS:
   1. Read each listed blueprint requirement for full context
   2. Implement the packet as one coherent slice of work
   3. Keep changes inside the owned files/modules unless a requirement forces expansion
   4. Write tests as needed
   5. Run validation: build must pass, tests must pass
   6. Commit with a message that names the packet's primary task BEFORE reporting
   7. Report result:
      TASK RESULT:
      - Tasks: {ids and titles}
      - Status: COMPLETE | PARTIAL | BLOCKED
      - Files changed: {list}
      - Issues: {any}"
   )
   ```

   Dispatch all packets for the wave in a single message with multiple Agent tool calls.

5. **After wave completes**:
   - For parallel waves: merge each subagent's branch into current branch (`git merge <branch> --no-edit`), then delete it (`git branch -D <branch>`). Skip if subagent reported no changes.
   - Update `context/impl/impl-*.md` with status for each completed task
   - Record any dead ends in `context/impl/dead-ends.md`

6. **Tier boundary check** — after updating impl tracking, check whether all tasks in the current tier are now done. If the current tier still has undone tasks, skip this step. If the tier is complete, run the Codex tier gate review (the `TIER_START_REF` was captured in step 1 at the start of this tier):

   a. Source `codex-config.sh` and check `tier_gate_mode` via `bp_config_get tier_gate_mode`. If the value is `"off"`, skip the review and log:
      ```
      [bp:tier-gate] Tier gate review disabled (tier_gate_mode=off). Skipping.
      ```

   b. Source `codex-detect.sh` and check `codex_available`. If `false`, log a note and continue:
      ```
      [bp:tier-gate] Codex unavailable — skipping tier boundary review. Continuing to next tier.
      ```

   c. Otherwise, run the review inline (wait for it to complete before advancing):
      ```
      scripts/codex-review.sh --base $TIER_START_REF
      ```

   d. **Severity-based gating** — after the review, source `scripts/codex-gate.sh` and run `bp_tier_gate`:
      - If `GATE_RESULT=proceed`: log the tier review summary and advance.
      - If `GATE_RESULT=blocked`: the tier has P0/P1 findings (or all findings in `strict` mode) that must be fixed before advancing.

   e. **When blocked** — generate fix tasks and execute them as an additional wave within the current tier:
      1. Run `bp_generate_fix_tasks` to get the list of blocking findings
      2. For each blocking finding, create a fix task: read the finding's file and description, implement the fix
      3. After all fix tasks complete, commit the fixes
      4. Mark fixed findings with `bp_findings_update_status <F-ID> FIXED`
      5. **Re-run the Codex review** on the updated diff: `scripts/codex-review.sh --base $TIER_START_REF`
      6. Re-evaluate with `bp_tier_gate`
      7. **Maximum 2 review-fix cycles per tier** — after 2 cycles, log any remaining P0/P1 findings as warnings and advance:
         ```
         [bp:tier-gate] WARNING: Advancing after 2 review-fix cycles with {N} unresolved P0/P1 findings
         ```

   ```
   ═══ Tier {N} Complete — Codex Review ═══
   Review: {CLEAN | N findings (M blocking, K deferred)}
   Gate: {PROCEED | BLOCKED → fix cycle {1|2}}
   ```

7. **Immediately proceed to next wave** — do NOT wait for user input between waves.

### Completion

When all tasks in the build site are done:

```
═══ BUILD COMPLETE ═══
Waves executed: {N}
Tasks completed: {done}/{total}
```

Then output the completion promise from the ralph prompt.

## Circuit Breakers

- **3 consecutive test failures on same task** → mark BLOCKED, document in dead-ends.md, skip
- **Merge conflict unresolvable** → stop the wave, report which branches conflict
- **All remaining tasks blocked** → report the dependency chain and stop

## Critical Rules

- Parallelize by work packet, not blindly by task count
- Group related small/medium tasks when they share files or context
- Split large or file-disjoint work so concurrent agents have clean ownership
- Merge after EVERY wave — do not accumulate unmerged branches
- Update impl tracking after EVERY wave — next wave reads it for frontier computation
- Progress through tiers autonomously — never pause between waves
- NEVER output completion promise unless ALL tasks are genuinely DONE
- NEVER mark a task DONE because existing code "looks related" — verify each acceptance criterion
