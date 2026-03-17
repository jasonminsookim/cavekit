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
	ActionTextInput
	ActionBackspace
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
			if overlayType == OverlayTextInput {
				return ActionTextInput
			}
		case "n":
			if overlayType == OverlayConfirmation {
				return ActionConfirmNo
			}
			if overlayType == OverlayTextInput {
				return ActionTextInput
			}
		case "esc", "ctrl+c":
			return ActionCancel
		case "backspace":
			if overlayType == OverlayTextInput {
				return ActionBackspace
			}
		case "j", "down":
			return ActionNavigateDown
		case "k", "up":
			return ActionNavigateUp
		}
		// In text input mode, treat single printable characters as text input
		if overlayType == OverlayTextInput && len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			return ActionTextInput
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
