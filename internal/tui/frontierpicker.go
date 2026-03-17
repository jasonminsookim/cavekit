package tui

import (
	"fmt"
	"strings"

	"github.com/julb/blueprint-monitor/internal/frontier"
)

// FrontierPicker shows available frontiers for selection.
type FrontierPicker struct {
	items         []FrontierPickerItem
	selectedIndex int
	multiSelect   map[int]bool
	visible       bool
}

// FrontierPickerItem represents a selectable frontier.
type FrontierPickerItem struct {
	Name      string
	Path      string
	Status    frontier.FrontierStatus
	TasksDone int
	TasksTotal int
}

// NewFrontierPicker creates a frontier picker.
func NewFrontierPicker() *FrontierPicker {
	return &FrontierPicker{
		multiSelect: make(map[int]bool),
	}
}

// SetItems populates the picker.
func (p *FrontierPicker) SetItems(items []FrontierPickerItem) {
	p.items = items
	p.selectedIndex = 0
	p.multiSelect = make(map[int]bool)
}

// Show makes the picker visible.
func (p *FrontierPicker) Show() {
	p.visible = true
}

// Hide hides the picker.
func (p *FrontierPicker) Hide() {
	p.visible = false
}

// IsVisible returns whether the picker is showing.
func (p *FrontierPicker) IsVisible() bool {
	return p.visible
}

// MoveDown moves selection down.
func (p *FrontierPicker) MoveDown() {
	if p.selectedIndex < len(p.items)-1 {
		p.selectedIndex++
	}
}

// MoveUp moves selection up.
func (p *FrontierPicker) MoveUp() {
	if p.selectedIndex > 0 {
		p.selectedIndex--
	}
}

// ToggleSelect toggles multi-select on current item.
func (p *FrontierPicker) ToggleSelect() {
	if p.selectedIndex < len(p.items) {
		item := p.items[p.selectedIndex]
		if item.Status == frontier.FrontierDone {
			return // Can't select done frontiers
		}
		p.multiSelect[p.selectedIndex] = !p.multiSelect[p.selectedIndex]
	}
}

// SelectedItems returns the selected frontiers.
func (p *FrontierPicker) SelectedItems() []FrontierPickerItem {
	var result []FrontierPickerItem
	if len(p.multiSelect) == 0 {
		// Single select mode: return current
		if p.selectedIndex < len(p.items) {
			result = append(result, p.items[p.selectedIndex])
		}
		return result
	}
	for i, selected := range p.multiSelect {
		if selected && i < len(p.items) {
			result = append(result, p.items[i])
		}
	}
	return result
}

// View renders the picker.
func (p *FrontierPicker) View() string {
	if !p.visible || len(p.items) == 0 {
		return ""
	}

	var rows []string
	rows = append(rows, OverlayTitleStyle.Render("Select Frontier")+"\n")

	for i, item := range p.items {
		marker := "  "
		if p.multiSelect[i] {
			marker = "● "
		}
		if i == p.selectedIndex {
			marker = "→ "
		}

		status := item.Status.Icon()
		progress := fmt.Sprintf("%d/%d", item.TasksDone, item.TasksTotal)

		style := NormalItemStyle
		if item.Status == frontier.FrontierDone {
			style = style.Strikethrough(true).Foreground(ColorMuted)
		}

		row := style.Render(fmt.Sprintf("%s%s %s %s", marker, status, item.Name, progress))
		rows = append(rows, row)
	}

	rows = append(rows, "\n"+MenuDescStyle.Render("Space to select · Enter to confirm · Esc to cancel"))
	return strings.Join(rows, "\n")
}
