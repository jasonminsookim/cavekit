---
created: "2026-03-17T00:00:00Z"
last_edited: "2026-03-17T00:00:00Z"
---
# Implementation Tracking: TUI

| Task | Status | Notes |
|------|--------|-------|
| T-019 | DONE | Lipgloss styles and constants: colors, panel styles, tab styles, status indicators, menu styles, overlay styles. internal/tui/styles.go. |
| T-018 | DONE | Bubbletea app shell with alt-screen, mouse, resize, tab switching, j/k navigation, menu bar, 30/70 layout. internal/tui/app.go. |
| T-024 | DONE | Instance list: status icons, progress display, j/k selection, mouse click. internal/tui/instancelist.go. |
| T-025 | DONE | Tabbed content: Preview/Diff/Terminal with tab switching and content truncation. internal/tui/tabs.go. |
| T-026 | DONE | Bottom menu: keyboard shortcuts with visual styling. internal/tui/menu.go. |
| T-027 | DONE | Overlay components: text input, confirmation, help screen. Esc to close. internal/tui/overlay.go. |
| T-029 | DONE | Preview tab: tmux pane capture with scroll mode. internal/tui/preview.go. |
| T-030 | DONE | Diff tab: git diff rendering with syntax coloring, stats header. internal/tui/difftab.go. |
| T-031 | DONE | Terminal tab: separate tmux session per instance, cached. internal/tui/terminaltab.go. |
| T-032 | DONE | New instance flow via overlay text input → session.Manager.Create/Start. |
| T-033 | DONE | Progress display via InstanceList rendering TasksDone/TasksTotal from session.UpdateProgress. |
| T-035 | DONE | Frontier picker: list available frontiers, multi-select, done strikethrough. internal/tui/frontierpicker.go. |
| T-036 | DONE | Key handling: MapKey routes all keypresses to actions, respects overlay state. internal/tui/keyhandler.go. |
| T-037 | DONE | Wire app.Update to route all key events via MapKey to components. |
| T-038 | DONE | Wire app.View to compose InstanceList + TabContent + BottomMenu + Overlay. |
