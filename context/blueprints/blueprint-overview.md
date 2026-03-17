---
created: "2026-03-17T00:00:00Z"
last_edited: "2026-03-17T00:00:00Z"
---

# Spec Overview

## Project
**blueprint-monitor** — A Go TUI application for managing multiple parallel Claude Code agents executing Blueprint frontiers. Replaces the current bash/tmux launcher with a flicker-free, claude-squad-style interface that adds Blueprint-specific progress tracking.

## Domain Index
| Domain | Spec File | Requirements | Status | Description |
|--------|-----------|-------------|--------|-------------|
| tmux | blueprint-tmux.md | 5 | DRAFT | Detached tmux session lifecycle, capture, attach/detach |
| worktree | blueprint-worktree.md | 4 | DRAFT | Git worktree creation, diff stats, discovery |
| frontier | blueprint-site.md | 5 | DRAFT | Frontier file discovery, parsing, task status tracking |
| session | blueprint-session.md | 6 | DRAFT | Instance model, lifecycle, persistence, progress |
| tui | blueprint-tui.md | 11 | DRAFT | Bubbletea TUI with list, tabs, overlays, Blueprint progress |
| cli | blueprint-cli.md | 5 | DRAFT | Binary, subcommands, config |

## Cross-Reference Map
| Domain A | Interacts With | Interaction Type |
|----------|---------------|-----------------|
| session | tmux | session creates and controls tmux sessions |
| session | worktree | session creates worktrees for isolation |
| session | frontier | session reads frontier for progress data |
| tui | session | TUI displays and controls sessions |
| tui | tmux | TUI captures pane content for preview |
| tui | frontier | TUI displays frontier progress |
| tui | worktree | TUI displays diff stats |
| cli | tui | CLI launches TUI |
| cli | session | CLI loads/saves session state |
| cli | tmux | kill command cleans up tmux sessions |
| cli | worktree | kill/status commands interact with worktrees |

## Dependency Graph
```
Tier 0 (no dependencies):    tmux, worktree, frontier
Tier 1 (depends on Tier 0):  session (depends on tmux, worktree, frontier)
Tier 2 (depends on Tier 1):  tui (depends on session, tmux, frontier, worktree)
Tier 3 (depends on Tier 2):  cli (depends on tui, session)
```

## Technology Stack
- **Language:** Go 1.22+
- **TUI framework:** charmbracelet/bubbletea
- **Styling:** charmbracelet/lipgloss
- **Components:** charmbracelet/bubbles (spinner, viewport, textinput)
- **PTY:** creack/pty
- **Tmux:** exec.Command wrapping tmux CLI
- **Git:** exec.Command wrapping git CLI
