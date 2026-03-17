package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/julb/blueprint-monitor/internal/session"
)

// InstanceList renders the left panel with all active instances.
type InstanceList struct {
	instances     []*session.Instance
	selectedIndex int
	width         int
	height        int
}

// NewInstanceList creates a new instance list component.
func NewInstanceList() *InstanceList {
	return &InstanceList{}
}

// SetInstances updates the instance list.
func (l *InstanceList) SetInstances(instances []*session.Instance) {
	l.instances = instances
}

// SetSelected sets the selected index.
func (l *InstanceList) SetSelected(index int) {
	if index < 0 {
		index = 0
	}
	if len(l.instances) > 0 && index >= len(l.instances) {
		index = len(l.instances) - 1
	}
	l.selectedIndex = index
}

// Selected returns the currently selected instance, or nil.
func (l *InstanceList) Selected() *session.Instance {
	if l.selectedIndex >= 0 && l.selectedIndex < len(l.instances) {
		return l.instances[l.selectedIndex]
	}
	return nil
}

// SelectedIndex returns the current selection index.
func (l *InstanceList) SelectedIndex() int {
	return l.selectedIndex
}

// SetSize updates the available dimensions.
func (l *InstanceList) SetSize(w, h int) {
	l.width = w
	l.height = h
}

// View renders the instance list.
func (l *InstanceList) View() string {
	if len(l.instances) == 0 {
		return lipgloss.NewStyle().
			Foreground(ColorMuted).
			Render("No instances.\n\nPress 'n' to create one.")
	}

	var rows []string
	for i, inst := range l.instances {
		row := l.renderRow(i, inst, i == l.selectedIndex)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func (l *InstanceList) renderRow(index int, inst *session.Instance, selected bool) string {
	// Status indicator
	var statusIcon string
	switch inst.Status {
	case session.StatusRunning:
		statusIcon = StatusRunning.String()
	case session.StatusReady:
		statusIcon = StatusReady.String()
	case session.StatusLoading:
		statusIcon = StatusLoading.String()
	case session.StatusPaused:
		statusIcon = StatusPaused.String()
	case session.StatusDone:
		statusIcon = StatusDone.String()
	}

	// Progress
	progress := ""
	if inst.TasksTotal > 0 {
		progress = fmt.Sprintf(" %d/%d", inst.TasksDone, inst.TasksTotal)
	}

	// Build row text
	text := fmt.Sprintf(" %d %s %s%s", index+1, statusIcon, inst.Title, progress)

	if selected {
		return SelectedItemStyle.Width(l.width).Render(text)
	}
	return NormalItemStyle.Width(l.width).Render(text)
}
