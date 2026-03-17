// Package tui implements the bubbletea-based terminal user interface.
package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/julb/blueprint-monitor/internal/session"
)

// Tab represents the active content tab.
type Tab int

const (
	TabPreview Tab = iota
	TabDiff
	TabTerminal
)

func (t Tab) String() string {
	switch t {
	case TabPreview:
		return "Preview"
	case TabDiff:
		return "Diff"
	case TabTerminal:
		return "Terminal"
	default:
		return "Unknown"
	}
}

// tickMsg triggers periodic updates (metadata, capture).
type tickMsg time.Time

// App is the main bubbletea model.
type App struct {
	width  int
	height int

	activeTab     Tab
	selectedIndex int

	// Components
	instanceList *InstanceList
	tabContent   *TabContent
	bottomMenu   *BottomMenu
	overlay      *Overlay

	// Data
	instances []*session.Instance

	// Set to true when we need to quit
	quitting bool
}

// NewApp creates a new TUI application model.
func NewApp() App {
	return App{
		activeTab:    TabPreview,
		instanceList: NewInstanceList(),
		tabContent:   NewTabContent(),
		bottomMenu:   NewBottomMenu(),
		overlay:      NewOverlay(),
	}
}

// Init implements tea.Model.
func (a App) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update implements tea.Model.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		action := MapKey(key, a.overlay.IsActive(), a.overlay.Active)

		switch action {
		case ActionQuit:
			a.quitting = true
			return a, tea.Quit
		case ActionSwitchTab:
			a.activeTab = (a.activeTab + 1) % 3
			a.tabContent.SetActiveTab(a.activeTab)
		case ActionNavigateDown:
			a.selectedIndex++
			a.instanceList.SetSelected(a.selectedIndex)
			a.selectedIndex = a.instanceList.SelectedIndex()
		case ActionNavigateUp:
			if a.selectedIndex > 0 {
				a.selectedIndex--
			}
			a.instanceList.SetSelected(a.selectedIndex)
			a.selectedIndex = a.instanceList.SelectedIndex()
		case ActionHelp:
			if a.overlay.Active == OverlayHelp {
				a.overlay.Hide()
			} else {
				a.overlay.Show(OverlayHelp, "", "")
			}
		case ActionCancel:
			a.overlay.Hide()
		case ActionNew:
			a.overlay.Show(OverlayTextInput, "New Instance", "Enter frontier name:")
		case ActionKill:
			if sel := a.instanceList.Selected(); sel != nil {
				a.overlay.Show(OverlayConfirmation, "Kill Instance", "Kill '"+sel.Title+"'?")
			}
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateLayout()

	case tickMsg:
		return a, tickCmd()
	}

	return a, nil
}

func (a *App) updateLayout() {
	leftWidth := max(int(float64(a.width)*LeftPanelRatio), MinLeftWidth)
	rightWidth := a.width - leftWidth - 2
	contentHeight := a.height - MenuHeight - 2

	a.instanceList.SetSize(leftWidth-2, contentHeight-2)
	a.tabContent.SetSize(rightWidth-2, contentHeight-2)
	a.bottomMenu.SetWidth(a.width)
	a.overlay.SetSize(a.width, a.height)
}

// View implements tea.Model.
func (a App) View() string {
	if a.quitting {
		return ""
	}
	if a.width == 0 || a.height == 0 {
		return "Initializing..."
	}

	leftWidth := max(int(float64(a.width)*LeftPanelRatio), MinLeftWidth)
	rightWidth := a.width - leftWidth - 2
	contentHeight := a.height - MenuHeight - 2

	// Render panels
	left := LeftPanelStyle.
		Width(leftWidth).
		Height(contentHeight).
		Render(a.instanceList.View())

	right := RightPanelStyle.
		Width(rightWidth).
		Height(contentHeight).
		Render(a.tabContent.View())

	menu := a.bottomMenu.View()

	// Compose layout
	panels := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	base := lipgloss.JoinVertical(lipgloss.Left, panels, menu)

	// Overlay on top
	if a.overlay.IsActive() {
		return a.overlay.View()
	}

	return base
}

// SetInstances updates the displayed instances.
func (a *App) SetInstances(instances []*session.Instance) {
	a.instances = instances
	a.instanceList.SetInstances(instances)
}

// Run starts the TUI application.
func Run() error {
	p := tea.NewProgram(
		NewApp(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
