//go:build windows

package action

import (
	"syscall"
	"unsafe"

	"github.com/yc446833448/VisuTask/internal/model"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procMouseEvent   = user32.NewProc("mouse_event")
	procGetCursorPos = user32.NewProc("GetCursorPos")
)

const (
	mouseEventLeftDown  = 0x0002
	mouseEventLeftUp    = 0x0004
	mouseEventRightDown = 0x0008
	mouseEventRightUp   = 0x0010
	mouseEventWheel     = 0x0800
)

// WinMouse implements MouseController using Windows user32.dll syscalls
type WinMouse struct{}

func NewWinMouse() *WinMouse { return &WinMouse{} }

func (m *WinMouse) Click(x, y int) error {
	if err := m.MoveTo(x, y); err != nil {
		return err
	}
	procMouseEvent.Call(uintptr(mouseEventLeftDown), 0, 0, 0, 0)
	procMouseEvent.Call(uintptr(mouseEventLeftUp), 0, 0, 0, 0)
	return nil
}

func (m *WinMouse) DoubleClick(x, y int) error {
	m.Click(x, y)
	m.Click(x, y)
	return nil
}

func (m *WinMouse) RightClick(x, y int) error {
	if err := m.MoveTo(x, y); err != nil {
		return err
	}
	procMouseEvent.Call(uintptr(mouseEventRightDown), 0, 0, 0, 0)
	procMouseEvent.Call(uintptr(mouseEventRightUp), 0, 0, 0, 0)
	return nil
}

func (m *WinMouse) Drag(from, to model.Point) error {
	if err := m.MoveTo(from.X, from.Y); err != nil {
		return err
	}
	procMouseEvent.Call(uintptr(mouseEventLeftDown), 0, 0, 0, 0)
	m.MoveTo(to.X, to.Y)
	procMouseEvent.Call(uintptr(mouseEventLeftUp), 0, 0, 0, 0)
	return nil
}

func (m *WinMouse) Scroll(x, y, delta int) error {
	if err := m.MoveTo(x, y); err != nil {
		return err
	}
	procMouseEvent.Call(uintptr(mouseEventWheel), 0, 0, uintptr(delta*120), 0)
	return nil
}

func (m *WinMouse) MoveTo(x, y int) error {
	ret, _, err := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return err
	}
	return nil
}

// point is used internally for cursor position
type point struct {
	X int32
	Y int32
}

func getCursorPos() (int, int) {
	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return int(p.X), int(p.Y)
}

// NewPlatformMouse returns the platform mouse implementation
func NewPlatformMouse() MouseController { return NewWinMouse() }
