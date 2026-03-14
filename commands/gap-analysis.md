---
name: sdd-gap-analysis
description: Compare what was built against what was intended
---

# SDD Gap Analysis — Built vs. Intended

You are performing a gap analysis: comparing what was actually built (implementation tracking, code, test results) against what was intended (specs, acceptance criteria). This identifies where specs, plans, or validation fell short and feeds into backpropagation.

Dispatch a `gap-analyzer` agent (via the Task tool) with the following instructions. If the Task tool is unavailable, execute the instructions directly.

## Agent Instructions for gap-analyzer

### Phase 1: Read Specs (Intended)

1. Read `context/specs/spec-overview.md` to get the full domain map
2. For each domain, read `context/specs/spec-{domain}.md`
3. Catalog every requirement and its acceptance criteria into a checklist:
   - Requirement ID (R1, R2, etc.)
   - Requirement description
   - Each acceptance criterion
   - Source spec file

### Phase 2: Read Implementation (Built)

1. Read all impl tracking files in `context/impl/`
2. Read `context/plans/plan-feature-frontier.md` for tier/progress state
3. Read `context/plans/plan-known-issues.md` for documented issues
4. Optionally run the test suite and capture results

### Phase 3: Compare

For each spec requirement and acceptance criterion, classify it:

| Status | Definition |
|--------|-----------|
| **COMPLETE** | Acceptance criteria met, tests pass, code exists |
| **PARTIAL** | Some criteria met, others missing or failing |
| **MISSING** | Not implemented at all |
| **OVER-BUILT** | Implementation goes beyond what the spec requires (scope creep) |
| **UNTESTABLE** | Acceptance criteria cannot be automatically validated |

### Phase 4: Identify Backpropagation Targets

For each gap (PARTIAL, MISSING, OVER-BUILT, UNTESTABLE), determine the root cause:

| Root Cause | Description | Fix Target |
|-----------|-------------|------------|
| **Spec gap** | Requirement was ambiguous or missing detail | Update spec |
| **Plan gap** | Spec was clear but plan didn't cover it | Update plan |
| **Implementation gap** | Plan covered it but implementation is incomplete | Continue implementing |
| **Validation gap** | Built but no test verifies it | Add tests |
| **Scope creep** | Built without a spec requirement | Either add spec or remove code |

### Phase 5: Generate Report

```markdown
## Gap Analysis Report

### Summary
| Status | Count | Percentage |
|--------|-------|-----------|
| COMPLETE | {n} | {%} |
| PARTIAL | {n} | {%} |
| MISSING | {n} | {%} |
| OVER-BUILT | {n} | {%} |
| UNTESTABLE | {n} | {%} |

### Overall Coverage: {complete + partial} / {total requirements}

### Detailed Findings

#### Complete
| Requirement | Spec | Evidence |
|------------|------|----------|
| R1: {name} | spec-{domain}.md | {test file or impl reference} |

#### Partial
| Requirement | Spec | Met | Missing | Root Cause |
|------------|------|-----|---------|-----------|
| R2: {name} | spec-{domain}.md | {criteria met} | {criteria missing} | {spec/plan/impl/validation gap} |

#### Missing
| Requirement | Spec | Root Cause | Suggested Action |
|------------|------|-----------|-----------------|
| R3: {name} | spec-{domain}.md | {root cause} | {action} |

#### Over-Built
| Feature | Files | Spec Coverage | Suggested Action |
|---------|-------|--------------|-----------------|
| {feature} | {files} | No spec | Add spec or remove |

#### Untestable
| Requirement | Spec | Why Untestable | Suggested Action |
|------------|------|---------------|-----------------|
| R4: {name} | spec-{domain}.md | {reason} | {rewrite criteria / add human review gate} |

### Backpropagation Targets
| Priority | Target File | Change Needed | Affected Requirements |
|----------|------------|--------------|----------------------|
| P0 | spec-{domain}.md | {what to update} | R{n}, R{n} |
| P1 | plan-{domain}.md | {what to update} | R{n} |

### Recommended Next Steps
1. Run `/sdd:back-propagate` to trace gaps into context files
2. {Specific spec updates needed}
3. {Specific plan updates needed}
4. {Implementation work remaining}
```

Present this report to the user when complete.
