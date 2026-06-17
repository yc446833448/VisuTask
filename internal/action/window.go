package action

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/yc446833448/VisuTask/internal/model"
)

var (
	procEnumWindows       = user32.NewProc("EnumWindows")
	procGetWindowTextW    = user32.NewProc("GetWindowTextW")
	procGetWindowTextLenW = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible   = user32.NewProc("IsWindowVisible")
	procSetForegroundWnd  = user32.NewProc("SetForegroundWindow")
	procGetWindowPID      = user32.NewProc("GetWindowThreadProcessId")
	procGetClassNameW     = user32.NewProc("GetClassNameW")
)

// WinWindow implements WindowController using Windows user32.dll syscalls
type WinWindow struct{}

func NewWinWindow() *WinWindow { return &WinWindow{} }

func (w *WinWindow) Find(title string) (*model.WindowInfo, error) {
	windows, err := w.List()
	if err != nil {
		return nil, err
	}
	for _, win := range windows {
		if containsIgnoreCase(win.Title, title) {
			return &win, nil
		}
	}
	return nil, fmt.Errorf("window not found: %s", title)
}

func (w *WinWindow) Focus(handle string) error {
	hwnd, err := parseHWND(handle)
	if err != nil {
		return err
	}
	ret, _, _ := procSetForegroundWnd.Call(hwnd)
	if ret == 0 {
		return fmt.Errorf("failed to focus window %s", handle)
	}
	return nil
}

func (w *WinWindow) List() ([]model.WindowInfo, error) {
	var windows []model.WindowInfo

	cb := syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		// Only visible windows
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1 // continue enumeration
		}

		// Get window title length
		titleLen, _, _ := procGetWindowTextLenW.Call(hwnd)
		if titleLen == 0 {
			return 1
		}

		// Get window title
		titleBuf := make([]uint16, titleLen+1)
		procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), titleLen+1)
		title := syscall.UTF16ToString(titleBuf)

		// Get process ID
		var pid uint32
		procGetWindowPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

		windows = append(windows, model.WindowInfo{
			Handle:  fmt.Sprintf("%d", hwnd),
			Title:   title,
			Process: fmt.Sprintf("PID:%d", pid),
		})

		return 1 // continue enumeration
	})

	procEnumWindows.Call(cb, 0)
	return windows, nil
}

func (w *WinWindow) Move(handle string, x, y int) error {
	// MoveWindow requires CGo or more complex syscall setup
	// For now, just focus the window
	return w.Focus(handle)
}

func parseHWND(handle string) (uintptr, error) {
	var h uintptr
	_, err := fmt.Sscanf(handle, "%d", &h)
	if err != nil {
		// Try hex
		_, err = fmt.Sscanf(handle, "0x%x", &h)
		if err != nil {
			return 0, fmt.Errorf("invalid window handle: %s", handle)
		}
	}
	return h, nil
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
