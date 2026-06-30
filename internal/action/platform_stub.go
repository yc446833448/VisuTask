//go:build !windows && !darwin

package action

// NewPlatformMouse returns a stub mouse implementation for unsupported platforms
func NewPlatformMouse() MouseController { return NewStubMouse() }

// NewPlatformKeyboard returns a stub keyboard implementation for unsupported platforms
func NewPlatformKeyboard() KeyboardController { return NewStubKeyboard() }

// NewPlatformWindow returns a stub window implementation for unsupported platforms
func NewPlatformWindow() WindowController { return NewStubWindow() }
