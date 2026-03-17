---
created: "2026-03-17T00:00:00Z"
last_edited: "2026-03-17T12:00:00Z"
---

# Feature Frontier

42 tasks across 6 tiers from 6 specs.

---

## Tier 0 — No Dependencies (Start Here)

| Task | Title | Spec | Requirement | Effort |
|------|-------|------|------------|--------|
| T-001 | Go module init, go.mod with dependencies | spec-cli.md | R1 | S |
| T-002 | Tmux session create/kill/exists | spec-tmux.md | R1 | M |
| T-003 | Tmux pane content capture | spec-tmux.md | R2 | S |
| T-004 | Git worktree create/detect/remove | spec-worktree.md | R1 | M |
| T-005 | Frontier file discovery and name derivation | spec-frontier.md | R1 | M |
| T-006 | Frontier markdown parsing (tasks, tiers, table rows) | spec-frontier.md | R2 | M |
| T-007 | Session instance model and status enum | spec-session.md | R1 | S |
| T-008 | Command executor abstraction (for testability) | spec-tmux.md | R1 | S |

---

## Tier 1 — Depends on Tier 0

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-009 | PTY-based tmux attach/detach with Ctrl+Q | spec-tmux.md | R3 | T-002 | L |
| T-010 | Tmux status detection (active/prompt via content hashing) | spec-tmux.md | R4 | T-002, T-003 | M |
| T-011 | Tmux input injection (keystrokes, prompts) | spec-tmux.md | R5 | T-002 | S |
| T-012 | Git diff stats (files changed, insertions, deletions) | spec-worktree.md | R2 | T-004 | S |
| T-013 | Worktree discovery (scan sibling dirs) | spec-worktree.md | R3 | T-004 | S |
| T-014 | Task status tracking from impl files | spec-frontier.md | R3 | T-005, T-006 | M |
| T-015 | Frontier status classification (done/in-progress/available) | spec-frontier.md | R4 | T-005, T-014 | S |
| T-016 | Frontier multi-candidate ranking and selection | spec-frontier.md | R6 | T-005, T-014 | M |
| T-017 | Progress summary string generation | spec-frontier.md | R5 | T-014, T-015 | S |
| T-018 | Bubbletea app shell (alt-screen, mouse, resize) | spec-tui.md | R1 | T-001 | M |
| T-019 | Lipgloss styles and constants | spec-tui.md | R1 | T-001 | S |

---

## Tier 2 — Depends on Tier 1

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-020 | Session lifecycle (create, start with worktree+tmux) | spec-session.md | R2 | T-002, T-004, T-007, T-011 | L |
| T-021 | Session persistence (save/load JSON) | spec-session.md | R3 | T-007, T-013 | M |
| T-022 | Auto-yes mode (permission prompt auto-approve) | spec-session.md | R5 | T-010, T-011 | S |
| T-023 | SDD progress integration on instance | spec-session.md | R6 | T-007, T-014 | S |
| T-024 | Instance list component | spec-tui.md | R2 | T-018, T-019, T-007 | M |
| T-025 | Tabbed content window component | spec-tui.md | R3 | T-018, T-019 | M |
| T-026 | Bottom menu component | spec-tui.md | R7 | T-018, T-019 | S |
| T-027 | Overlay components (text input, confirmation, help) | spec-tui.md | R10 | T-018, T-019 | M |
| T-028 | Branch push from worktree | spec-worktree.md | R4 | T-004, T-012 | S |

---

## Tier 3 — Depends on Tier 2

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-029 | Preview tab (capture + render tmux snapshots) | spec-tui.md | R4 | T-025, T-003, T-020 | M |
| T-030 | Diff tab (git diff rendering with scroll) | spec-tui.md | R5 | T-025, T-012 | M |
| T-031 | Terminal tab (separate tmux session per instance) | spec-tui.md | R6 | T-025, T-002, T-009 | L |
| T-032 | New instance flow (name input → start) | spec-tui.md | R8 | T-024, T-027, T-020 | M |
| T-033 | SDD progress display in instance list | spec-tui.md | R9 | T-024, T-023, T-017 | M |
| T-034 | Staggered launch for multiple instances | spec-session.md | R4 | T-020 | S |
| T-035 | Frontier picker integration | spec-tui.md | R11 | T-027, T-005, T-015 | M |
| T-036 | Key handling: kill, push, checkout, resume | spec-tui.md | R8 | T-024, T-027, T-020, T-028 | M |

---

## Tier 4 — Depends on Tier 3

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-037 | Wire app.Update to route all key events | spec-tui.md | R1 | T-029, T-030, T-031, T-032, T-033, T-036 | L |
| T-038 | Wire app.View to compose all components | spec-tui.md | R1 | T-024, T-025, T-026, T-027, T-029, T-030, T-031, T-033 | M |
| T-039 | Monitor command (launch TUI with preflight) | spec-cli.md | R2 | T-018, T-021, T-037 | M |

---

## Tier 5 — Depends on Tier 4

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-040 | Status command (print progress, exit) | spec-cli.md | R3 | T-013, T-014, T-017 | S |
| T-041 | Kill command (cleanup sessions, worktrees, branches) | spec-cli.md | R4 | T-002, T-004, T-039 | S |
| T-042 | Config file and debug/reset/version commands | spec-cli.md | R5 | T-039 | S |

---

## Summary

| Tier | Tasks | Effort |
|------|-------|--------|
| 0 | 8 | 2S, 5M, 1S |
| 1 | 11 | 4S, 5M, 1L, 1S |
| 2 | 9 | 3S, 5M, 1L |
| 3 | 8 | 1S, 5M, 1L, 1M |
| 4 | 3 | 1L, 2M |
| 5 | 3 | 3S |

**Total: 42 tasks, 6 tiers**

---

## Tier 6 — Inspection Fixes (depends on Tier 5)

| Task | Title | Spec | Requirement | blockedBy | Effort | Finding |
|------|-------|------|------------|-----------|--------|---------|
| T-043 | Wire tick loop: capture, diff, progress, auto-yes, status | spec-tui.md | R12 | T-037 | L | F-001 |
| T-044 | Instantiate PreviewTab/DiffTab/TerminalTab in App, pipe to TabContent | spec-tui.md | R12 | T-029, T-030, T-031, T-037 | M | F-002 |
| T-045 | Handle ActionOpen: tmux attach/detach with TUI suspend | spec-tui.md | R13 | T-009, T-037 | M | F-003 |
| T-046 | Handle ActionPush: branch push with confirmation overlay | spec-tui.md | R13 | T-028, T-027, T-037 | S | F-006 |
| T-047 | Handle ActionCheckout/ActionResume | spec-tui.md | R13 | T-037 | S | F-006 |
| T-048 | Handle ActionScrollUp/Down for preview and diff tabs | spec-tui.md | R13 | T-029, T-030, T-037 | S | F-008 |
| T-049 | Integrate FrontierPicker into new-instance flow | spec-tui.md | R11 | T-035, T-032 | M | F-007 |
| T-050 | Validate persistence on load (tmux/worktree existence) | spec-session.md | R3 | T-021 | S | F-009 |
| T-051 | Parse --autoyes/-y flag, pass to TUI and AutoYes | spec-cli.md | R2 | T-039 | S | F-010 |
| T-052 | Status command: show progress counts from frontier tracking | spec-cli.md | R3 | T-040, T-014, T-017 | S | F-012 |
| T-053 | Instance list: add branch name and diff stats to rows | spec-tui.md | R2 | T-024, T-012 | S | F-011 |
| T-054 | DiffTab: apply scrollPos to Content() output | spec-tui.md | R5 | T-030 | S | F-014 |
| T-055 | Context-adaptive menu items | spec-tui.md | R7 | T-026 | S | F-013 |
| T-056 | Preflight check for program binary | spec-cli.md | R2 | T-039 | S | F-015 |

**Inspection additions: 14 tasks, 1 tier**
**New total: 56 tasks, 7 tiers**
