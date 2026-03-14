---
name: sdd-methodology
description: |
  Core Spec-Driven Development (SDD) methodology — the master skill that teaches the SPIIM lifecycle
  and routes to all sub-skills. Covers the Idea → Specs → Code principle, the UDP/TCP reliability analogy,
  5-phase SPIIM lifecycle, decision matrix for when to use SDD, and CI pipeline analogy.
  Trigger phrases: "use SDD", "spec-driven", "start SDD project", "sdd methodology",
  "how should I structure this project for AI agents"
---

# Spec-Driven Development (SDD) Methodology

## Core Principle: Idea → Specs → Code

**Never go directly from source material to code. Always go through specs first.**

SDD is a methodology for building software with AI coding agents that treats **specifications as the primary artifact**, not code. Whether starting from scratch or modernizing an existing system, the principle is the same:

- **Greenfield projects:** reference material → specs → code
- **Rewrites:** old code → specs → new code

In both cases, the specs become a durable, iterable asset that agents can consume to continuously build, test, and improve the application.

### Why Specs Are the First-Class Citizen

| Property | Benefit |
|----------|---------|
| **Hierarchical** | Organized as a tree, enabling progressive disclosure |
| **Reviewable** | Humans can audit at a higher level than code |
| **Portable** | Framework-agnostic, reusable across technology stacks |
| **Iterable** | Improve specs without touching code |
| **Testable** | Define acceptance criteria that agents validate against |

> **Key Insight:** If your specs are good enough and your validation is strong enough, you can regenerate the entire application from specs at any time. This is the "continuous integration" of AI development.

---

## The UDP/TCP Analogy

LLMs are inherently non-deterministic — like UDP, each individual call is unreliable. But through the right guardrails — rigorous specs, validation requirements, and iterative convergence loops — we build a reliable, deterministic process on top of an unreliable substrate.

**SDD is the TCP layer over LLM-powered agents.**

| Layer | Analogy | What It Does |
|-------|---------|-------------|
| **LLM calls** | UDP packets | Each call may produce different output; no guarantee of correctness |
| **Specifications** | Protocol definition | Define what "correct" means — the contract |
| **Validation gates** | Checksums / ACKs | Verify each output meets the contract |
| **Convergence loops** | Retransmission | Re-run until output stabilizes and passes all gates |
| **Implementation tracking** | Sequence numbers | Track progress, prevent duplicate/lost work |
| **Backpropagation** | Error correction | Trace failures back to source, fix the protocol |

The result: a reliable, reproducible development process built from unreliable components.

---

## The 5 SPIIM Phases

SPIIM stands for **Spec, Plan, Implement, Iterate, Monitor**. Each phase has dedicated prompts that drive it.

| Phase | Input | Output | AI Role | Human Role |
|-------|-------|--------|---------|------------|
| **Spec** | Old code, reference docs, research | Implementation-agnostic specs | Analyze, document, organize | Review specs for completeness |
| **Plan** | Specs + framework research | Framework-specific implementation plans | Architect, decompose, sequence | Validate architecture decisions |
| **Implement** | Plans + specs | Working code + tests + tracking docs | Build, test, validate | Monitor progress |
| **Iterate** | Failed validations, gaps, manual fixes | Updated specs/plans via backpropagation | Diagnose, backpropagate, fix | Audit results, steer direction |
| **Monitor** | Running application, git history | Issues, anomalies, progress reports | Observe, scan, report | Review reports, steer direction |

### Phase Transitions

Each phase has **gate conditions** that must be met before moving to the next:

1. **Spec → Plan:** All domains have specs with testable acceptance criteria. Human has reviewed for completeness.
2. **Plan → Implement:** Plans reference specs, define implementation sequence, and include test strategies. Architecture decisions validated.
3. **Implement → Iterate:** Code builds, tests pass at current coverage level, implementation tracking is up to date.
4. **Iterate → Monitor:** Convergence detected (changes decreasing iteration-over-iteration). Remaining changes are trivial.
5. **Monitor → Spec (cycle):** Gap found or new requirement identified. Backpropagate to specs and restart the cycle.

The **Iterate** phase is where the human acts as an **auditor, not an implementer**. You monitor the process, request changes as needed, and make systemic improvements to specs and prompts.

> For the full SPIIM phase reference, see `references/spiim-phases.md`.

---

## Decision Matrix: When to Use SDD

### Full SDD

Use when the project has significant scope, evolving requirements, or needs autonomous agent execution.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | 50+ source files |
| Requirements | Evolving, multi-domain |
| Agent coordination | Multi-agent or multi-prompt pipelines |
| Environment | Production, security-sensitive, brownfield |
| Team structure | Multi-team or cross-team |
| Execution mode | Long-running autonomous work (overnight, unattended) |

**What you get:** Full SPIIM lifecycle, context directory with specs/plans/impl tracking, prompt pipeline, convergence loops, backpropagation, validation gates.

### Lightweight SDD

Use when scope is moderate — too complex for ad-hoc but not worth a full pipeline.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | 5-50 files |
| Requirements | Mostly clear, focused |
| Agent coordination | Single agent, possibly with sub-agents |
| Execution mode | Interactive with occasional iteration loops |

**What you do:**
1. Write a focused `context/specs/spec-task.md` capturing requirements
2. Add a `context/plans/plan-task.md` sequencing the implementation
3. Skip full SPIIM — just run an iteration loop against the plan

This is the "SDD floor" — most of the benefit without the overhead of a full multi-phase pipeline.

### Skip SDD

Use when the task is trivially small.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | Less than 5 files |
| Task type | One-off tools, simple bug fixes, exploratory prototypes |
| Implementation | Can be done from memory in a single session |

**Heuristic:** If you can implement it from memory in a single session, it is probably too small for full SDD.

### Growth Path

Start with lightweight SDD even if the project is small. If the scope expands, you already have the structure in place to scale up. It is much harder to retrofit specs onto a large codebase than to grow a spec directory from the beginning.

---

## The CI Pipeline Analogy

SDD is, at its core, a **CI pipeline for AI development:**

```
Traditional CI/CD:
  Code → Build → Test → Deploy

SDD AI Pipeline:
  Spec Change
    → Generate Plans (iteration loop)
    → Generate Implementation (iteration loop)
    → Validate (Tests + Review)
    → Human Audit (Monitor & Steer)
    → [Gap Found]
    → Backpropagate
    → Spec Change (cycle repeats)
```

Each stage can run as an iteration loop — the same prompt executed repeatedly until output converges. The convergence loop is the engine that turns unreliable LLM calls into reliable software.

### The Iteration Loop

The iteration loop is the fundamental execution primitive in SDD. Run the same prompt over the same codebase multiple times until changes converge to zero.

**How it works:**
1. Run a prompt against the codebase
2. Agent reads git history and tracking files to understand prior state
3. Agent makes changes, commits progress
4. Repeat from step 1

**Convergence signal:** An exponentially decreasing number of changes from run to run (200 lines changed, then 100, then 50, then approximately 20 lines of minor adjustments). The signal is not zero changes — it is that remaining changes are trivial.

**Non-convergence is a signal to improve your specs, not to run more iterations.**

If changes are not decreasing:
- Your specs are too fuzzy (requirements ambiguous)
- Your validation is too weak (agent cannot verify correctness)
- Sub-agents are fighting each other (conflicting ownership)

---

## Cross-References to Sub-Skills

SDD is composed of techniques that work together. This methodology skill is the index — each sub-skill below is self-contained but cross-references others.

### Foundation Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `sdd:spec-writing` | Write implementation-agnostic specs with testable acceptance criteria | Spec phase — always the first step |
| `sdd:context-architecture` | Organize context for progressive disclosure | Project setup and ongoing maintenance |
| `sdd:impl-tracking` | Track implementation progress, dead ends, test health | Implement and Iterate phases |
| `sdd:validation-first` | Design validation gates agents can execute | All phases — validation is continuous |

### Pipeline Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `sdd:prompt-pipeline` | Design numbered prompt pipelines for SPIIM | Setting up automation |
| `sdd:backpropagation` | Trace bugs back to specs and fix at the source | Iterate phase — after finding gaps |
| `sdd:brownfield-adoption` | Adopt SDD on existing codebases | Starting SDD on legacy projects |

### Advanced Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `sdd:adversarial-review` | Use a second agent to challenge the first | Quality gates, architecture review |
| `sdd:leader-follower` | Stagger pipeline stages for parallelism | Optimizing long pipelines |
| `sdd:convergence-monitoring` | Detect convergence vs ceiling | Monitoring iteration loops |
| `sdd:documentation-inversion` | Turn documentation into agent-consumable skills | Library/module documentation |

### Integration with Existing Skills

SDD works **with** existing skills, not as a replacement:

| Existing Skill | SDD Integration |
|----------------|-----------------|
| `superpowers:brainstorming` | Use during spec generation to explore requirements |
| `superpowers:writing-plans` | Use during plan generation for structured planning |
| `superpowers:test-driven-development` | TDD-within-SDD: spec acceptance criteria become failing tests |
| `superpowers:verification-before-completion` | Use for gate validation in every phase |
| `superpowers:executing-plans` | Use during implementation phase |
| `superpowers:dispatching-parallel-agents` | Use for agent team coordination |

---

## Quick Start

### For a New Project (Greenfield)

1. **Set up context directory:**
   ```
   context/
   ├── refs/           # Source materials (PRDs, language specs, research)
   ├── specs/          # Implementation-agnostic specifications
   ├── plans/          # Framework-specific implementation plans
   ├── impl/           # Living implementation tracking
   └── prompts/        # SPIIM pipeline prompts
   ```

2. **Write specs** from your reference materials (see `sdd:spec-writing`)
3. **Generate plans** from specs (see `sdd:prompt-pipeline`)
4. **Implement** with validation gates (see `sdd:validation-first`)
5. **Track progress** in implementation documents (see `sdd:impl-tracking`)
6. **Iterate** — when gaps are found, backpropagate to specs (see `sdd:backpropagation`)

### For an Existing Project (Brownfield)

1. **Set up context directory** (same structure as above)
2. **Designate existing codebase as reference material**
3. **Generate specs from code** (see `sdd:brownfield-adoption`)
4. **Validate specs match behavior** — run tests against generated specs
5. **Proceed with normal SPIIM** — future changes flow through specs first

---

## Summary

SDD is not a tool — it is a methodology. The core loop is simple:

1. **Describe what you want** (specs with testable criteria)
2. **Let agents build it** (plans → implementation → validation)
3. **Fix the specs, not the code** (backpropagation)
4. **Repeat until converged** (iteration loops)

The more structure you give agents — through well-crafted specs, validation gates, and iterative convergence — the more they can do without you hovering over every line. This does NOT replace software engineers. Your expertise in architecture, spec design, system thinking, and quality judgment is exactly what makes SDD work. The methodology amplifies what you can do — it turns one engineer's vision into a full implementation pipeline.
