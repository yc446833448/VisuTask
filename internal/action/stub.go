package action

import "github.com/yc446833448/VisuTask/internal/model"

// Stub implementations for development

type StubMouse struct{}

func NewStubMouse() *StubMouse { return &StubMouse{} }

func (s *StubMouse) Click(x, y int) error           { return nil }
func (s *StubMouse) DoubleClick(x, y int) error      { return nil }
func (s *StubMouse) RightClick(x, y int) error       { return nil }
func (s *StubMouse) Drag(from, to model.Point) error  { return nil }
func (s *StubMouse) Scroll(x, y, delta int) error     { return nil }
func (s *StubMouse) MoveTo(x, y int) error            { return nil }

type StubKeyboard struct{}

func NewStubKeyboard() *StubKeyboard { return &StubKeyboard{} }

func (s *StubKeyboard) Type(text string) error      { return nil }
func (s *StubKeyboard) HotKey(keys ...string) error  { return nil }
func (s *StubKeyboard) KeyPress(key string) error    { return nil }

type StubWindow struct{}

func NewStubWindow() *StubWindow { return &StubWindow{} }

func (s *StubWindow) Find(title string) (*model.WindowInfo, error) {
	return nil, nil // TODO: implement
}

func (s *StubWindow) Focus(handle string) error { return nil }

func (s *StubWindow) List() ([]model.WindowInfo, error) {
	return []model.WindowInfo{}, nil // TODO: implement
}

func (s *StubWindow) Move(handle string, x, y int) error { return nil }
