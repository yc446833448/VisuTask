//go:build darwin

package action

import (
	"fmt"
	"os/exec"
	"strings"
)

// ─── MacKeyboard ───

// MacKeyboard implements KeyboardController using AppleScript System Events.
// Requires Accessibility permissions to be granted (System Preferences > Security & Privacy).
type MacKeyboard struct{}

func NewMacKeyboard() *MacKeyboard { return &MacKeyboard{} }

// Type simulates typing a text string.
// ASCII text is typed via System Events keystroke.
// CJK text is pasted via clipboard to ensure correct input.
func (k *MacKeyboard) Type(text string) error {
	if hasCJK(text) {
		return k.pasteText(text)
	}
	return k.keystroke(text)
}

// HotKey presses a key combination (e.g., "command", "c" for Cmd+C).
// Modifier keys should come first, followed by the main key.
func (k *MacKeyboard) HotKey(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// Map key names to macOS modifier names for AppleScript
	modifiers := []string{}
	var mainKey string

	for _, key := range keys {
		lower := strings.ToLower(key)
		switch lower {
		case "command", "cmd", "super", "win":
			modifiers = append(modifiers, "command down")
		case "shift":
			modifiers = append(modifiers, "shift down")
		case "option", "alt":
			modifiers = append(modifiers, "option down")
		case "control", "ctrl":
			modifiers = append(modifiers, "control down")
		default:
			mainKey = key
		}
	}

	if mainKey == "" {
		// All keys were modifiers — just press the last one
		mainKey = keys[len(keys)-1]
	}

	modPart := ""
	if len(modifiers) > 0 {
		modPart = " using {" + strings.Join(modifiers, ", ") + "}"
	}

	// Escape special characters in AppleScript strings
	escapedKey := escapeAppleScript(mainKey)
	script := fmt.Sprintf(`tell application "System Events" to keystroke "%s"%s`, escapedKey, modPart)
	return runAppleScript(script)
}

// KeyPress presses and releases a single key
func (k *MacKeyboard) KeyPress(key string) error {
	lower := strings.ToLower(key)
	// Map to AppleScript key code or keystroke
	switch lower {
	case "enter", "return":
		return runAppleScript(`tell application "System Events" to keystroke return`)
	case "tab":
		return runAppleScript(`tell application "System Events" to keystroke tab`)
	case "space":
		return runAppleScript(`tell application "System Events" to keystroke space`)
	case "escape", "esc":
		return runAppleScript(`tell application "System Events" to key code 53`)
	case "backspace":
		return runAppleScript(`tell application "System Events" to keystroke (ASCII character 8)`)
	case "forwarddelete":
		return runAppleScript(`tell application "System Events" to key code 117`)
	case "home":
		return runAppleScript(`tell application "System Events" to key code 115`)
	case "end":
		return runAppleScript(`tell application "System Events" to key code 119`)
	case "pageup":
		return runAppleScript(`tell application "System Events" to key code 116`)
	case "pagedown":
		return runAppleScript(`tell application "System Events" to key code 121`)
	case "up":
		return runAppleScript(`tell application "System Events" to key code 126`)
	case "down":
		return runAppleScript(`tell application "System Events" to key code 125`)
	case "left":
		return runAppleScript(`tell application "System Events" to key code 123`)
	case "right":
		return runAppleScript(`tell application "System Events" to key code 124`)
	case "f1", "f2", "f3", "f4", "f5", "f6",
		"f7", "f8", "f9", "f10", "f11", "f12":
		// F1-F12: key codes 122-133
		fNum := int(lower[1] - '0')
		if len(lower) == 3 {
			fNum = fNum*10 + int(lower[2]-'0')
		}
		keyCode := 121 + fNum
		script := fmt.Sprintf(`tell application "System Events" to key code %d`, keyCode)
		return runAppleScript(script)
	default:
		if len(key) == 1 {
			escaped := escapeAppleScript(key)
			script := fmt.Sprintf(`tell application "System Events" to keystroke "%s"`, escaped)
			return runAppleScript(script)
		}
		return fmt.Errorf("unsupported key: %s", key)
	}
}

// keystroke sends a string via System Events keystroke
func (k *MacKeyboard) keystroke(text string) error {
	escaped := escapeAppleScript(text)
	script := fmt.Sprintf(`tell application "System Events" to keystroke "%s"`, escaped)
	return runAppleScript(script)
}

// pasteText copies text to clipboard and pastes it (for CJK input)
func (k *MacKeyboard) pasteText(text string) error {
	// Step 1: Copy text to clipboard via pbcopy
	copyCmd := exec.Command("pbcopy")
	copyCmd.Stdin = strings.NewReader(text)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("pbcopy failed: %w", err)
	}

	// Step 2: Paste via Cmd+V
	script := `tell application "System Events" to keystroke "v" using {command down}`
	return runAppleScript(script)
}

// escapeAppleScript escapes special characters for AppleScript string literals
func escapeAppleScript(s string) string {
	// AppleScript uses backslash escape: \\ for backslash, \" for double quote
	escaped := strings.ReplaceAll(s, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	// Handle line breaks
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	escaped = strings.ReplaceAll(escaped, "\r", "\\r")
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")
	return escaped
}

// NewPlatformKeyboard returns the platform keyboard implementation
func NewPlatformKeyboard() KeyboardController { return NewMacKeyboard() }
