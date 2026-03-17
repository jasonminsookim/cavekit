package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TabContent renders the right panel with tabbed content.
type TabContent struct {
	activeTab Tab
	width     int
	height    int

	// Content for each tab
	previewContent  string
	diffContent     string
	terminalContent string
}

// NewTabContent creates a new tabbed content component.
func NewTabContent() *TabContent {
	return &TabContent{
		activeTab: TabPreview,
	}
}

// SetActiveTab sets the currently visible tab.
func (t *TabContent) SetActiveTab(tab Tab) {
	t.activeTab = tab
}

// SetSize updates available dimensions.
func (t *TabContent) SetSize(w, h int) {
	t.width = w
	t.height = h
}

// SetPreview sets the preview tab content.
func (t *TabContent) SetPreview(content string) {
	t.previewContent = content
}

// SetDiff sets the diff tab content.
func (t *TabContent) SetDiff(content string) {
	t.diffContent = content
}

// SetTerminal sets the terminal tab content.
func (t *TabContent) SetTerminal(content string) {
	t.terminalContent = content
}

// View renders the tabbed content panel.
func (t *TabContent) View() string {
	tabBar := t.renderTabBar()

	var content string
	switch t.activeTab {
	case TabPreview:
		content = t.previewContent
		if content == "" {
			content = "Select an instance to preview."
		}
	case TabDiff:
		content = t.diffContent
		if content == "" {
			content = "No diff available."
		}
	case TabTerminal:
		content = t.terminalContent
		if content == "" {
			content = "Press Enter to open terminal."
		}
	}

	// Truncate content to fit
	lines := strings.Split(content, "\n")
	maxLines := t.height - 2 // account for tab bar
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	return tabBar + "\n" + strings.Join(lines, "\n")
}

func (t *TabContent) renderTabBar() string {
	var tabs []string
	for _, tab := range []Tab{TabPreview, TabDiff, TabTerminal} {
		if tab == t.activeTab {
			tabs = append(tabs, ActiveTabStyle.Render(tab.String()))
		} else {
			tabs = append(tabs, InactiveTabStyle.Render(tab.String()))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}
