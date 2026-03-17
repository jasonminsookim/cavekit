// Package tui implements the bubbletea-based terminal user interface.
package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/julb/blueprint-monitor/internal/exec"
	"github.com/julb/blueprint-monitor/internal/session"
	"github.com/julb/blueprint-monitor/internal/tmux"
	"github.com/julb/blueprint-monitor/internal/worktree"
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

// instanceCreatedMsg is sent when a new instance has been started in the background.
type instanceCreatedMsg struct {
	inst *session.Instance
	err  error
}

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

	// Session management
	sessionMgr  *session.Manager
	store       *session.Store
	projectRoot string
	program     string

	// Set to true when we need to quit
	quitting bool
}

// NewApp creates a new TUI application model.
func NewApp(projectRoot, program string) App {
	executor := exec.NewRealExecutor()
	tmuxMgr := tmux.NewManager(executor)
	wtMgr := worktree.NewManager(executor)
	sessMgr := session.NewManager(tmuxMgr, wtMgr)
	store := session.NewStore("")

	return App{
		activeTab:    TabPreview,
		instanceList: NewInstanceList(),
		tabContent:   NewTabContent(),
		bottomMenu:   NewBottomMenu(),
		overlay:      NewOverlay(),
		sessionMgr:   sessMgr,
		store:        store,
		projectRoot:  projectRoot,
		program:      program,
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
			a.overlay.Show(OverlayTextInput, "New Instance", "Enter instance name:")
		case ActionKill:
			if sel := a.instanceList.Selected(); sel != nil {
				a.overlay.Show(OverlayConfirmation, "Kill Instance", "Kill '"+sel.Title+"'?")
			}
		case ActionTextInput:
			if a.overlay.Active == OverlayTextInput {
				r := msg.Runes
				if len(r) > 0 {
					a.overlay.InputValue += string(r)
				}
			}
		case ActionBackspace:
			if a.overlay.Active == OverlayTextInput && len(a.overlay.InputValue) > 0 {
				a.overlay.InputValue = a.overlay.InputValue[:len(a.overlay.InputValue)-1]
			}
		case ActionConfirmYes:
			switch a.overlay.Active {
			case OverlayTextInput:
				name := a.overlay.InputValue
				if name != "" {
					a.overlay.Hide()
					return a, a.createInstance(name)
				}
			case OverlayConfirmation:
				if sel := a.instanceList.Selected(); sel != nil {
					a.overlay.Hide()
					a.sessionMgr.Kill(context.Background(), sel, a.projectRoot, true)
					a.removeInstance(sel)
				}
			}
		case ActionConfirmNo:
			a.overlay.Hide()
		}

	case tea.MouseMsg:
		if !a.overlay.IsActive() {
			leftWidth := max(int(float64(a.width)*LeftPanelRatio), MinLeftWidth)

			if msg.X < leftWidth {
				// Click in instance list — select the clicked row
				// Account for border (1 row top) and row height (1 per instance)
				row := msg.Y - 1
				if row >= 0 && row < len(a.instances) {
					a.selectedIndex = row
					a.instanceList.SetSelected(a.selectedIndex)
					a.selectedIndex = a.instanceList.SelectedIndex()
				}
			} else {
				// Click in right panel — check if clicking tab bar
				if msg.Y <= 2 {
					// Clicked in tab bar area — cycle tab on click
					a.activeTab = (a.activeTab + 1) % 3
					a.tabContent.SetActiveTab(a.activeTab)
				}
			}
		}

	case instanceCreatedMsg:
		if msg.err == nil && msg.inst != nil {
			a.instances = append(a.instances, msg.inst)
			a.instanceList.SetInstances(a.instances)
			a.selectedIndex = len(a.instances) - 1
			a.instanceList.SetSelected(a.selectedIndex)
			a.store.Save(a.instances)
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

func (a *App) createInstance(name string) tea.Cmd {
	return func() tea.Msg {
		prog := a.program
		if prog == "" {
			prog = "claude"
		}
		inst := a.sessionMgr.Create(name, "", name, prog)
		err := a.sessionMgr.Start(context.Background(), inst, a.projectRoot, name, 3*time.Second)
		return instanceCreatedMsg{inst: inst, err: err}
	}
}

func (a *App) removeInstance(inst *session.Instance) {
	for i, ins := range a.instances {
		if ins == inst {
			a.instances = append(a.instances[:i], a.instances[i+1:]...)
			break
		}
	}
	a.instanceList.SetInstances(a.instances)
	if a.selectedIndex >= len(a.instances) && a.selectedIndex > 0 {
		a.selectedIndex = len(a.instances) - 1
	}
	a.instanceList.SetSelected(a.selectedIndex)
	a.store.Save(a.instances)
}

// Run starts the TUI application.
func Run(projectRoot, program string) error {
	app := NewApp(projectRoot, program)
	// Load persisted instances
	instances, _ := app.store.Load()
	if len(instances) > 0 {
		app.instances = instances
		app.instanceList.SetInstances(instances)
	}

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
