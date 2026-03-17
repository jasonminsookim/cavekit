package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// OverlayType identifies the kind of overlay.
type OverlayType int

const (
	OverlayNone OverlayType = iota
	OverlayTextInput
	OverlayConfirmation
	OverlayHelp
)

// Overlay renders a centered modal on top of the main content.
type Overlay struct {
	Active      OverlayType
	Title       string
	Message     string
	InputValue  string
	Width       int
	Height      int
}

// NewOverlay creates an inactive overlay.
func NewOverlay() *Overlay {
	return &Overlay{Active: OverlayNone}
}

// Show activates the overlay.
func (o *Overlay) Show(typ OverlayType, title, message string) {
	o.Active = typ
	o.Title = title
	o.Message = message
	o.InputValue = ""
}

// Hide deactivates the overlay.
func (o *Overlay) Hide() {
	o.Active = OverlayNone
	o.Title = ""
	o.Message = ""
	o.InputValue = ""
}

// IsActive returns true if an overlay is showing.
func (o *Overlay) IsActive() bool {
	return o.Active != OverlayNone
}

// SetSize updates the available screen dimensions for centering.
func (o *Overlay) SetSize(w, h int) {
	o.Width = w
	o.Height = h
}

// View renders the overlay content (without background).
func (o *Overlay) View() string {
	if !o.IsActive() {
		return ""
	}

	var content string
	switch o.Active {
	case OverlayTextInput:
		content = OverlayTitleStyle.Render(o.Title) + "\n\n" +
			o.Message + "\n\n" +
			"> " + o.InputValue + "█" + "\n\n" +
			MenuDescStyle.Render("Enter to confirm · Esc to cancel")

	case OverlayConfirmation:
		content = OverlayTitleStyle.Render(o.Title) + "\n\n" +
			o.Message + "\n\n" +
			MenuKeyStyle.Render("y") + " yes  " +
			MenuKeyStyle.Render("n") + " no"

	case OverlayHelp:
		content = OverlayTitleStyle.Render("Keyboard Shortcuts") + "\n\n" +
			helpText() + "\n\n" +
			MenuDescStyle.Render("Press Esc or ? to close")
	}

	// Calculate overlay dimensions
	overlayWidth := min(60, o.Width-4)
	rendered := OverlayStyle.Width(overlayWidth).Render(content)

	// Center vertically and horizontally
	return lipgloss.Place(o.Width, o.Height, lipgloss.Center, lipgloss.Center, rendered)
}

func helpText() string {
	return `Navigation:
  j/k, ↑/↓    Navigate instances
  Tab          Switch tab (Preview/Diff/Terminal)
  Enter/o      Attach to selected instance
  Ctrl+Q       Detach from instance

Instance Management:
  n            New instance
  D            Kill selected instance
  p            Push branch
  c            Checkout worktree

Other:
  ?            Toggle help
  q, Ctrl+C   Quit`
}
