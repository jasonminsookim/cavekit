---
name: sdd-execute
description: "Start a Ralph Loop that implements the feature frontier — builds, validates, commits, repeats"
argument-hint: "[--filter PATTERN] [--adversarial] [--max-iterations N] [--completion-promise TEXT]"
allowed-tools: ["Bash(${CLAUDE_PLUGIN_ROOT}/scripts/setup-execute.sh:*)"]
---

# SDD Execute — Run the Implementation Loop

This is the third phase of SDD. Execute the setup script:

```!
"${CLAUDE_PLUGIN_ROOT}/scripts/setup-execute.sh" $ARGUMENTS
```

You are now in a Ralph Loop implementing the feature frontier. Follow the prompt instructions exactly.

## How This Works

1. Archives any previous loop cycle automatically
2. Reads the feature frontier to find unblocked tasks
3. Each iteration: pick task → read spec → implement → validate → commit
4. With `--adversarial`: alternates build iterations with Codex (GPT-5.4) review
5. Exits when all tasks are done

## Critical Rules

- NEVER output the completion promise unless ALL tasks are genuinely DONE
- ONE task per iteration — stay focused
- If stuck 2+ iterations, document as dead end and move on
- Always run validation gates before committing
