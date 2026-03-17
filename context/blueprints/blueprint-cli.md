---
created: "2026-03-17T00:00:00Z"
last_edited: "2026-03-17T00:00:00Z"
---

# Spec: CLI Interface

## Scope
The command-line interface that replaces the current `blueprint` bash script. Provides the `blueprint` binary with subcommands for monitor, status, kill, and analytics.

## Requirements

### R1: Binary and Installation
**Description:** A single Go binary named `blueprint` that can be installed via `go install` or downloaded.
**Acceptance Criteria:**
- [ ] Compiles to a single static binary
- [ ] Module path: `github.com/julb/blueprint-monitor` (or similar)
- [ ] Supports `go install github.com/julb/blueprint-monitor@latest`
- [ ] Binary name is `blueprint` (or `blueprint-monitor` to avoid conflict with the existing `blueprint` script)
**Dependencies:** none

### R2: Monitor Command (Default)
**Description:** The primary command that launches the TUI.
**Acceptance Criteria:**
- [ ] `blueprint` or `blueprint monitor` launches the TUI
- [ ] `--program <cmd>` overrides the default program (default: `claude`)
- [ ] `--autoyes` / `-y` enables auto-approval of permission prompts
- [ ] Preflight checks: tmux installed, program (claude) installed, git repo detected
- [ ] Loads persisted instances from previous session
**Dependencies:** R1, blueprint-tui R1

### R3: Status Command
**Description:** Shows frontier progress without launching the TUI.
**Acceptance Criteria:**
- [ ] `blueprint status` prints per-worktree progress to stdout
- [ ] Format: `{name}: {icon} {done}/{total} tasks done`
- [ ] Works from any terminal (doesn't require the TUI to be running)
- [ ] Exits after printing
**Dependencies:** R1, blueprint-site R3, blueprint-worktree R3

### R4: Kill Command
**Description:** Stops all Blueprint sessions and cleans up.
**Acceptance Criteria:**
- [ ] `blueprint kill` kills all `blueprint_*` tmux sessions
- [ ] Removes all `{project}-blueprint-*` worktrees
- [ ] Deletes all `blueprint/*` branches
- [ ] Cleans up `.claude/ralph-loop.local.md` from project root and worktrees
- [ ] Reports count of killed sessions, cleaned worktrees, deleted branches
**Dependencies:** R1, blueprint-tmux R1, blueprint-worktree R1

### R5: Configuration
**Description:** Persistent configuration for the monitor.
**Acceptance Criteria:**
- [ ] Config file at `~/.blueprint-monitor/config.json`
- [ ] Configurable: default_program, stagger_delay, max_instances (default 10)
- [ ] `blueprint debug` prints config paths for troubleshooting
- [ ] `blueprint reset` clears all stored instances
- [ ] `blueprint version` prints the version
**Dependencies:** R1

## Out of Scope
- Analytics command (keep as separate bash script for now)
- Merge command (keep as Claude Code slash command)
- Plugin system

## Cross-References
- See also: blueprint-tui.md (monitor launches TUI)
- See also: blueprint-session.md (persistence paths)
