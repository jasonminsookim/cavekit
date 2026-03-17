package tui

// KeyAction represents a user action triggered by a keypress.
type KeyAction int

const (
	ActionNone KeyAction = iota
	ActionNew
	ActionKill
	ActionOpen
	ActionPush
	ActionCheckout
	ActionResume
	ActionSwitchTab
	ActionHelp
	ActionQuit
	ActionNavigateUp
	ActionNavigateDown
	ActionConfirmYes
	ActionConfirmNo
	ActionCancel
	ActionScrollUp
	ActionScrollDown
)

// MapKey maps a key string to an action, respecting overlay state.
func MapKey(key string, overlayActive bool, overlayType OverlayType) KeyAction {
	// When overlay is active, only overlay-specific keys work
	if overlayActive {
		switch key {
		case "enter":
			return ActionConfirmYes
		case "y":
			if overlayType == OverlayConfirmation {
				return ActionConfirmYes
			}
		case "n":
			if overlayType == OverlayConfirmation {
				return ActionConfirmNo
			}
		case "esc", "ctrl+c":
			return ActionCancel
		case "j", "down":
			return ActionNavigateDown
		case "k", "up":
			return ActionNavigateUp
		}
		return ActionNone
	}

	// Normal mode
	switch key {
	case "n":
		return ActionNew
	case "D":
		return ActionKill
	case "enter", "o":
		return ActionOpen
	case "p":
		return ActionPush
	case "c":
		return ActionCheckout
	case "r":
		return ActionResume
	case "tab":
		return ActionSwitchTab
	case "?":
		return ActionHelp
	case "q", "ctrl+c":
		return ActionQuit
	case "j", "down":
		return ActionNavigateDown
	case "k", "up":
		return ActionNavigateUp
	case "shift+up", "K":
		return ActionScrollUp
	case "shift+down", "J":
		return ActionScrollDown
	}
	return ActionNone
}
