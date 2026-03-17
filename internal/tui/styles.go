package tui

import "github.com/charmbracelet/lipgloss"

// Layout proportions
const (
	LeftPanelRatio  = 0.30
	RightPanelRatio = 0.70
	MinLeftWidth    = 25
	MenuHeight      = 1
)

// Colors
var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // purple
	ColorSecondary = lipgloss.Color("#6B7280") // gray
	ColorSuccess   = lipgloss.Color("#10B981") // green
	ColorWarning   = lipgloss.Color("#F59E0B") // yellow
	ColorDanger    = lipgloss.Color("#EF4444") // red
	ColorMuted     = lipgloss.Color("#4B5563") // dark gray
	ColorBorder    = lipgloss.Color("#374151") // border gray
	ColorHighlight = lipgloss.Color("#1F2937") // subtle bg
)

// Panel styles
var (
	LeftPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	RightPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(ColorHighlight).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle()
)

// Tab styles
var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Padding(0, 1)

	TabBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(ColorBorder)
)

// Status indicators
var (
	StatusRunning = lipgloss.NewStyle().Foreground(ColorSuccess).SetString("●")
	StatusReady   = lipgloss.NewStyle().Foreground(ColorWarning).SetString("●")
	StatusLoading = lipgloss.NewStyle().Foreground(ColorSecondary).SetString("◌")
	StatusPaused  = lipgloss.NewStyle().Foreground(ColorMuted).SetString("⏸")
	StatusDone    = lipgloss.NewStyle().Foreground(ColorSuccess).SetString("✓")
)

// Menu styles
var (
	MenuStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	MenuKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E5E7EB"))

	MenuDescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)

// Overlay styles
var (
	OverlayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	OverlayTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary)
)

// Diff styles
var (
	DiffAddStyle    = lipgloss.NewStyle().Foreground(ColorSuccess)
	DiffRemoveStyle = lipgloss.NewStyle().Foreground(ColorDanger)
	DiffHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
)
