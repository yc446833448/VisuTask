//go:build darwin

package action

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yc446833448/VisuTask/internal/model"
)

// ─── MacWindow ───

// MacWindow implements WindowController using AppleScript/JXA System Events.
// Requires Accessibility permissions to be granted (System Preferences > Security & Privacy).
type MacWindow struct{}

func NewMacWindow() *MacWindow { return &MacWindow{} }

// Find locates a window by fuzzy-matching its title (case-insensitive)
func (w *MacWindow) Find(title string) (*model.WindowInfo, error) {
	windows, err := w.List()
	if err != nil {
		return nil, err
	}
	for _, win := range windows {
		if strings.Contains(strings.ToLower(win.Title), strings.ToLower(title)) {
			return &win, nil
		}
	}
	return nil, fmt.Errorf("window not found: %s", title)
}

// Focus brings a window to the foreground by its handle (window ID)
func (w *MacWindow) Focus(handle string) error {
	// Parse handle: "ProcessName_windowID" or "windowID"
	parts := strings.SplitN(handle, "_", 2)
	if len(parts) == 2 {
		processName := parts[0]
		winID := parts[1]
		script := fmt.Sprintf(`
tell application "System Events"
	tell process "%s"
		set frontmost to true
		repeat with w in every window
			if (id of w as string) is equal to "%s" then
				perform action "AXRaise" of w
				exit repeat
			end if
		end repeat
	end tell
end tell`, escapeAppleScriptStr(processName), winID)
		return runAppleScript(script)
	}
	// Fallback: try to activate by window title
	return fmt.Errorf("cannot focus by title-only handle on macOS, use full handle: %s", handle)
}

// List returns all visible windows across all running applications.
// Uses JXA (JavaScript for Automation) which returns JSON natively.
func (w *MacWindow) List() ([]model.WindowInfo, error) {
	script := `
var sys = Application("System Events");
var windows = [];
var processes = sys.processes();

for (var i = 0; i < processes.length; i++) {
	try {
		var p = processes[i];
		if (p.backgroundOnly() === true) continue;
		var wins = p.windows();
		for (var j = 0; j < wins.length; j++) {
			try {
				var win = wins[j];
				var title = win.name();
				if (!title || title.length === 0) continue;
				// Only include windows with a non-empty title
				windows.push({
					handle: String(p.name()) + "_" + String(win.id()),
					title: title,
					process: p.name() + ":" + String(p.id())
				});
			} catch(e) {}
		}
	} catch(e) {}
}
JSON.stringify(windows);
`
	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %v — %s", err, string(output))
	}

	var windows []model.WindowInfo
	if err := json.Unmarshal(output, &windows); err != nil {
		return nil, fmt.Errorf("failed to parse window list: %v (output: %s)", err, string(output))
	}
	return windows, nil
}

// Move moves a window to the specified position
func (w *MacWindow) Move(handle string, x, y int) error {
	parts := strings.SplitN(handle, "_", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid window handle format, expected: ProcessName_windowID")
	}
	processName := parts[0]
	winID := parts[1]

	script := fmt.Sprintf(`
tell application "System Events"
	tell process "%s"
		repeat with w in every window
			if (id of w as string) is equal to "%s" then
				set position of w to {%d, %d}
				exit repeat
			end if
		end repeat
	end tell
end tell`, escapeAppleScriptStr(processName), winID, x, y)
	return runAppleScript(script)
}

// escapeAppleScriptStr escapes a string for AppleScript double-quoted literals
func escapeAppleScriptStr(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// NewPlatformWindow returns the platform window implementation
func NewPlatformWindow() WindowController { return NewMacWindow() }
