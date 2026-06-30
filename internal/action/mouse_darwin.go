//go:build darwin

package action

import (
	"fmt"
	"os/exec"

	"github.com/yc446833448/VisuTask/internal/model"
)

// ─── AppleScript helpers ───

// runAppleScript executes an AppleScript and returns any error
func runAppleScript(script string) error {
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript error: %v — %s", err, string(output))
	}
	return nil
}

// ─── MacMouse ───

// MacMouse implements MouseController using AppleScript + CoreGraphics framework.
// Requires Accessibility permissions to be granted (System Preferences > Security & Privacy).
type MacMouse struct{}

func NewMacMouse() *MacMouse { return &MacMouse{} }

// Click performs a left mouse click at (x, y)
func (m *MacMouse) Click(x, y int) error {
	script := fmt.Sprintf(`
use framework "CoreGraphics"
set pos to current application's CGPointMake(%d, %d)
-- Move cursor
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), pos)
delay 0.02
-- Left button down
set downEvent to current application's CGEventCreateMouseEvent(missing value, 1, pos, 0)
current application's CGEventPost(0, downEvent)
delay 0.02
-- Left button up
set upEvent to current application's CGEventCreateMouseEvent(missing value, 2, pos, 0)
current application's CGEventPost(0, upEvent)
`, x, y)
	return runAppleScript(script)
}

// DoubleClick performs a double left click at (x, y)
func (m *MacMouse) DoubleClick(x, y int) error {
	if err := m.Click(x, y); err != nil {
		return err
	}
	return m.Click(x, y)
}

// RightClick performs a right mouse click at (x, y)
func (m *MacMouse) RightClick(x, y int) error {
	script := fmt.Sprintf(`
use framework "CoreGraphics"
set pos to current application's CGPointMake(%d, %d)
-- Move cursor
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), pos)
delay 0.02
-- Right button down (kCGMouseButtonRight = 1)
set downEvent to current application's CGEventCreateMouseEvent(missing value, 3, pos, 1)
current application's CGEventPost(0, downEvent)
delay 0.02
-- Right button up
set upEvent to current application's CGEventCreateMouseEvent(missing value, 4, pos, 1)
current application's CGEventPost(0, upEvent)
`, x, y)
	return runAppleScript(script)
}

// Drag performs a drag operation from one point to another
func (m *MacMouse) Drag(from, to model.Point) error {
	script := fmt.Sprintf(`
use framework "CoreGraphics"
set fromPos to current application's CGPointMake(%d, %d)
set toPos to current application's CGPointMake(%d, %d)
-- Move to start position
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), fromPos)
delay 0.05
-- Left button down
set downEvent to current application's CGEventCreateMouseEvent(missing value, 1, fromPos, 0)
current application's CGEventPost(0, downEvent)
delay 0.05
-- Move to end position
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), toPos)
delay 0.05
-- Left button up
set upEvent to current application's CGEventCreateMouseEvent(missing value, 2, toPos, 0)
current application's CGEventPost(0, upEvent)
`, from.X, from.Y, to.X, to.Y)
	return runAppleScript(script)
}

// Scroll performs a scroll wheel action at (x, y) with given delta
func (m *MacMouse) Scroll(x, y, delta int) error {
	// macOS scroll: kCGScrollEventUnitLine = 0, kCGScrollEventUnitPixel = 1
	// We scale delta by 10 for reasonable scroll amount (matches ~3 lines per 1 delta unit)
	scrollAmount := delta * 10
	script := fmt.Sprintf(`
use framework "CoreGraphics"
set pos to current application's CGPointMake(%d, %d)
-- Move cursor to target
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), pos)
delay 0.02
-- Create scroll event (type 22 = kCGEventScrollWheel)
set scrollEvent to current application's CGEventCreateScrollWheelEvent(missing value, 0, 1, %d)
current application's CGEventPost(0, scrollEvent)
`, x, y, scrollAmount)
	return runAppleScript(script)
}

// MoveTo moves the mouse cursor to (x, y)
func (m *MacMouse) MoveTo(x, y int) error {
	script := fmt.Sprintf(`
use framework "CoreGraphics"
set pos to current application's CGPointMake(%d, %d)
current application's CGDisplayMoveCursorToPoint(current application's CGMainDisplayID(), pos)
`, x, y)
	return runAppleScript(script)
}

// NewPlatformMouse returns the platform mouse implementation
func NewPlatformMouse() MouseController { return NewMacMouse() }
