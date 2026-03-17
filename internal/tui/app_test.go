package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app.activeTab != TabPreview {
		t.Errorf("default tab should be Preview, got %v", app.activeTab)
	}
	if app.quitting {
		t.Error("should not be quitting on init")
	}
}

func TestApp_Update_Quit(t *testing.T) {
	app := NewApp()
	model, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	updated := model.(App)
	if !updated.quitting {
		t.Error("should be quitting after 'q'")
	}
	if cmd == nil {
		t.Error("should return tea.Quit command")
	}
}

func TestApp_Update_TabSwitch(t *testing.T) {
	app := NewApp()
	if app.activeTab != TabPreview {
		t.Fatal("should start on Preview")
	}

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyTab})
	updated := model.(App)
	if updated.activeTab != TabDiff {
		t.Errorf("after tab, should be Diff, got %v", updated.activeTab)
	}

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyTab})
	updated = model.(App)
	if updated.activeTab != TabTerminal {
		t.Errorf("after tab, should be Terminal, got %v", updated.activeTab)
	}

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyTab})
	updated = model.(App)
	if updated.activeTab != TabPreview {
		t.Errorf("after tab, should wrap to Preview, got %v", updated.activeTab)
	}
}

func TestApp_Update_Navigate(t *testing.T) {
	app := NewApp()
	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := model.(App)
	if updated.selectedIndex != 1 {
		t.Errorf("j should increment index, got %d", updated.selectedIndex)
	}

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = model.(App)
	if updated.selectedIndex != 0 {
		t.Errorf("k should decrement index, got %d", updated.selectedIndex)
	}

	// k at 0 should not go negative
	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = model.(App)
	if updated.selectedIndex != 0 {
		t.Errorf("k at 0 should stay 0, got %d", updated.selectedIndex)
	}
}

func TestApp_Update_Resize(t *testing.T) {
	app := NewApp()
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	updated := model.(App)
	if updated.width != 120 || updated.height != 40 {
		t.Errorf("size should be 120x40, got %dx%d", updated.width, updated.height)
	}
}

func TestApp_View_BeforeResize(t *testing.T) {
	app := NewApp()
	view := app.View()
	if view != "Initializing..." {
		t.Errorf("before resize, should show initializing, got %q", view)
	}
}

func TestApp_View_AfterResize(t *testing.T) {
	app := NewApp()
	model, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	updated := model.(App)
	view := updated.View()

	if !strings.Contains(view, "Preview") {
		t.Error("view should contain Preview tab")
	}
	if !strings.Contains(view, "Diff") {
		t.Error("view should contain Diff tab")
	}
	if !strings.Contains(view, "quit") {
		t.Error("view should contain menu with quit")
	}
}

func TestTab_String(t *testing.T) {
	if TabPreview.String() != "Preview" {
		t.Errorf("TabPreview.String() = %q", TabPreview.String())
	}
	if TabDiff.String() != "Diff" {
		t.Errorf("TabDiff.String() = %q", TabDiff.String())
	}
	if TabTerminal.String() != "Terminal" {
		t.Errorf("TabTerminal.String() = %q", TabTerminal.String())
	}
}
