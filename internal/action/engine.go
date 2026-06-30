package action

import "github.com/yc446833448/VisuTask/internal/model"

// Engine provides system-level input simulation
type Engine struct {
	mouse    MouseController
	keyboard KeyboardController
	window   WindowController
}

func NewEngine(m MouseController, k KeyboardController, w WindowController) *Engine {
	return &Engine{mouse: m, keyboard: k, window: w}
}

// ─── Mouse ───

func (e *Engine) Click(x, y int) error       { return e.mouse.Click(x, y) }
func (e *Engine) DoubleClick(x, y int) error  { return e.mouse.DoubleClick(x, y) }
func (e *Engine) RightClick(x, y int) error   { return e.mouse.RightClick(x, y) }
func (e *Engine) Drag(from, to model.Point) error { return e.mouse.Drag(from, to) }
func (e *Engine) Scroll(x, y, delta int) error { return e.mouse.Scroll(x, y, delta) }
func (e *Engine) MoveTo(x, y int) error       { return e.mouse.MoveTo(x, y) }

// ─── Keyboard ───

func (e *Engine) Type(text string) error           { return e.keyboard.Type(text) }
func (e *Engine) HotKey(keys ...string) error      { return e.keyboard.HotKey(keys...) }
func (e *Engine) KeyPress(key string) error        { return e.keyboard.KeyPress(key) }

// ─── Window ───

func (e *Engine) FindWindow(title string) (*model.WindowInfo, error) { return e.window.Find(title) }
func (e *Engine) FocusWindow(handle string) error                     { return e.window.Focus(handle) }
func (e *Engine) ListWindows() ([]model.WindowInfo, error)            { return e.window.List() }
func (e *Engine) MoveWindow(handle string, x, y int) error            { return e.window.Move(handle, x, y) }

// ─── Interfaces ───

type MouseController interface {
	Click(x, y int) error
	DoubleClick(x, y int) error
	RightClick(x, y int) error
	Drag(from, to model.Point) error
	Scroll(x, y, delta int) error
	MoveTo(x, y int) error
}

type KeyboardController interface {
	Type(text string) error
	HotKey(keys ...string) error
	KeyPress(key string) error
}

type WindowController interface {
	Find(title string) (*model.WindowInfo, error)
	Focus(handle string) error
	List() ([]model.WindowInfo, error)
	Move(handle string, x, y int) error
}

// ─── Shared helpers ───

// hasCJK checks if the string contains CJK (Chinese/Japanese/Korean) characters.
// Used by both Windows and macOS keyboard implementations.
func hasCJK(s string) bool {
	for _, r := range s {
		if r > 0x2E80 {
			return true
		}
	}
	return false
}
