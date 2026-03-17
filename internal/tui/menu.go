package tui

import "github.com/charmbracelet/lipgloss"

// MenuItem represents a menu entry.
type MenuItem struct {
	Key  string
	Desc string
}

// DefaultMenu returns the standard menu items.
func DefaultMenu() []MenuItem {
	return []MenuItem{
		{"n", "new"},
		{"D", "kill"},
		{"Enter", "open"},
		{"p", "push"},
		{"c", "checkout"},
		{"tab", "switch tab"},
		{"?", "help"},
		{"q", "quit"},
	}
}

// RenderMenu renders a menu bar from items.
func RenderMenu(items []MenuItem, width int) string {
	var parts []string
	for _, item := range items {
		parts = append(parts,
			MenuKeyStyle.Render(item.Key)+" "+MenuDescStyle.Render(item.Desc))
	}

	content := joinMenuParts(parts, " │ ")
	return MenuStyle.Width(width).Render(content)
}

// BottomMenu is a stateful menu component.
type BottomMenu struct {
	items []MenuItem
	width int
}

// NewBottomMenu creates the bottom menu with default items.
func NewBottomMenu() *BottomMenu {
	return &BottomMenu{
		items: DefaultMenu(),
	}
}

// SetWidth updates the menu width.
func (m *BottomMenu) SetWidth(w int) {
	m.width = w
}

// SetItems replaces the menu items (for context-dependent menus).
func (m *BottomMenu) SetItems(items []MenuItem) {
	m.items = items
}

// View renders the menu.
func (m *BottomMenu) View() string {
	return RenderMenu(m.items, m.width)
}

// joinWithSep is re-declared here to avoid depending on app.go ordering.
func joinMenuParts(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

// RenderMenuCompact renders a compact version for narrow terminals.
func RenderMenuCompact(items []MenuItem) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, lipgloss.NewStyle().Bold(true).Render(item.Key))
	}
	return MenuStyle.Render(joinMenuParts(parts, " "))
}
