---
created: "2026-03-17T00:00:00Z"
last_edited: "2026-03-17T00:00:00Z"
---

# Spec: Terminal User Interface

## Scope
The bubbletea-based TUI that replaces the current tmux split-pane layout. Renders a session list, preview/diff/terminal tabs, and a bottom menu — all in a single process that controls rendering (no flicker).

## Requirements

### R1: Application Shell
**Description:** The main bubbletea application with alt-screen, mouse support, and component layout.
**Acceptance Criteria:**
- [ ] Uses `tea.WithAltScreen()` for full-screen takeover
- [ ] Uses `tea.WithMouseCellMotion()` for mouse scroll support
- [ ] Layout: left panel (30% width) = instance list, right panel (70%) = tabbed content, bottom = menu bar
- [ ] Responds to terminal resize events
- [ ] Clean exit saves state and restores terminal
**Dependencies:** none

### R2: Instance List
**Description:** Left panel showing all active instances with status indicators.
**Acceptance Criteria:**
- [ ] Each row shows: index number, title, branch name, diff stats (+N, -N), status indicator
- [ ] Status indicators: green dot = Running, yellow dot = Ready, spinner = Loading, grey = Paused
- [ ] Blueprint-specific: shows `{done}/{total} tasks` next to each instance
- [ ] Navigate with j/k or arrow keys
- [ ] Selected instance highlighted with border/background
- [ ] Mouse click selects instance
**Dependencies:** R1, blueprint-session R1

### R3: Tabbed Content Window
**Description:** Right panel with Preview, Diff, and Terminal tabs.
**Acceptance Criteria:**
- [ ] Tab bar at top with Preview | Diff | Terminal labels
- [ ] Switch tabs with Tab key
- [ ] Active tab is visually highlighted
- [ ] Each tab renders independently based on the selected instance
**Dependencies:** R1

### R4: Preview Tab
**Description:** Shows a snapshot of the selected instance's tmux pane content.
**Acceptance Criteria:**
- [ ] Captures tmux pane content every 100ms for the selected instance only
- [ ] Renders content with ANSI color preservation
- [ ] Supports scroll mode (Shift+Up/Down to enter, Esc to exit)
- [ ] In scroll mode, captures full scrollback history
- [ ] Shows fallback message when no instance is selected or instance is paused
**Dependencies:** R3, blueprint-tmux R2

### R5: Diff Tab
**Description:** Shows git diff between the instance's branch and main.
**Acceptance Criteria:**
- [ ] Renders `git diff main...HEAD` output with syntax highlighting
- [ ] Scrollable with Shift+Up/Down
- [ ] Shows diff stats summary at top (files changed, +insertions, -deletions)
- [ ] Updates when instance changes or on metadata tick
**Dependencies:** R3, blueprint-worktree R2

### R6: Terminal Tab
**Description:** An independent terminal session in the instance's worktree directory.
**Acceptance Criteria:**
- [ ] Creates a separate tmux session per instance for shell access in the worktree
- [ ] Displays captured pane content (same as preview but for the terminal session)
- [ ] Press Enter to attach full-screen (interact directly)
- [ ] Sessions are cached per instance and preserved when switching between instances
**Dependencies:** R3, blueprint-tmux R1

### R7: Bottom Menu
**Description:** Shows available keyboard shortcuts.
**Acceptance Criteria:**
- [ ] Displays: `n new | D kill | Enter/o open | p push branch | c checkout | tab switch tab | ? help | q quit`
- [ ] Highlights the active shortcut briefly when pressed (visual feedback)
- [ ] Menu adapts to current state (different options when creating new instance vs. default)
**Dependencies:** R1

### R8: New Instance Flow
**Description:** Interactive flow for creating a new agent instance.
**Acceptance Criteria:**
- [ ] Press `n` to start: shows name input field at the bottom of the instance list
- [ ] Press `N` (Shift+N) to start with prompt: shows name input, then prompt+branch picker overlay
- [ ] Name input validates: non-empty, max 32 chars
- [ ] On Enter: creates instance, starts worktree+tmux, sends `/bp:build --filter {name}`
- [ ] Esc cancels and removes the pending instance
**Dependencies:** R2, R7, blueprint-session R2

### R9: Blueprint Progress Display
**Description:** Blueprint-specific progress information integrated into the UI.
**Acceptance Criteria:**
- [ ] Instance list shows task progress: `3/12` or `✓` if complete
- [ ] Progress bar or fraction visible per-instance
- [ ] Current working task ID shown when instance is Running
- [ ] Tier completion markers (Tier 0 ✓, Tier 1 >, Tier 2 -)
**Dependencies:** R2, blueprint-site R3, blueprint-site R5

### R10: Overlays
**Description:** Modal overlays for text input, confirmations, and help screens.
**Acceptance Criteria:**
- [ ] Text input overlay: centered modal with text field (for prompts)
- [ ] Confirmation overlay: yes/no modal for destructive actions (kill, push)
- [ ] Help overlay: shows keyboard shortcuts and usage info
- [ ] Overlays capture all input while active
- [ ] Esc or Ctrl+C closes overlays
**Dependencies:** R1

### R11: Frontier Picker Integration
**Description:** On startup (or when pressing `n`), show available frontiers for selection.
**Acceptance Criteria:**
- [ ] Lists frontiers from `context/frontiers/` with status (available/in-progress/done)
- [ ] Shows task count per frontier
- [ ] Done frontiers shown as struck-through/disabled
- [ ] In-progress frontiers shown with resume indicator
- [ ] Multi-select support for launching multiple agents at once
**Dependencies:** R8, blueprint-site R1, blueprint-site R4

## Out of Scope
- Web-based UI
- Multiple project support (single project root per TUI instance)
- Theming/color customization (uses sensible defaults)

## Cross-References
- See also: blueprint-session.md (TUI controls sessions)
- See also: blueprint-tmux.md (preview captures)
- See also: blueprint-site.md (progress data)
- See also: blueprint-cli.md (TUI is launched by CLI)
