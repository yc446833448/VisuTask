//go:build windows

package action

import (
	"strings"
	"unicode/utf16"
	"unsafe"
)

var (
	procKeybdEvent = user32.NewProc("keybd_event")
	procSendInput  = user32.NewProc("SendInput")
)

const (
	keyEventDown    = 0x0000
	keyEventUp      = 0x0002
	keyEventUnicode = 0x0004
	inputKeyboard   = 1
)

// keyboardInput matches Windows KEYBDINPUT + INPUT struct layout
type keyboardInput struct {
	Type uint32
	// KEYBDINPUT
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

// WinKeyboard implements KeyboardController using Windows syscalls
type WinKeyboard struct{}

func NewWinKeyboard() *WinKeyboard { return &WinKeyboard{} }

func (k *WinKeyboard) Type(text string) error {
	if hasCJK(text) {
		return typeUnicode(text)
	}
	// ASCII: use keybd_event per character
	for _, ch := range text {
		vk := charToVK(ch)
		shift := needsShift(ch)
		if shift {
			procKeybdEvent.Call(0x10, 0, uintptr(keyEventDown), 0) // VK_SHIFT down
		}
		procKeybdEvent.Call(uintptr(vk), 0, uintptr(keyEventDown), 0)
		procKeybdEvent.Call(uintptr(vk), 0, uintptr(keyEventUp), 0)
		if shift {
			procKeybdEvent.Call(0x10, 0, uintptr(keyEventUp), 0) // VK_SHIFT up
		}
	}
	return nil
}

func (k *WinKeyboard) HotKey(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	vkCodes := make([]uint16, len(keys))
	for i, key := range keys {
		vkCodes[i] = keyToVK(key)
	}
	for _, vk := range vkCodes {
		procKeybdEvent.Call(uintptr(vk), 0, uintptr(keyEventDown), 0)
	}
	for i := len(vkCodes) - 1; i >= 0; i-- {
		procKeybdEvent.Call(uintptr(vkCodes[i]), 0, uintptr(keyEventUp), 0)
	}
	return nil
}

func (k *WinKeyboard) KeyPress(key string) error {
	vk := keyToVK(key)
	procKeybdEvent.Call(uintptr(vk), 0, uintptr(keyEventDown), 0)
	procKeybdEvent.Call(uintptr(vk), 0, uintptr(keyEventUp), 0)
	return nil
}

// typeUnicode types Unicode text using SendInput with KEYEVENTF_UNICODE
func typeUnicode(text string) error {
	chars := utf16.Encode([]rune(text))

	for _, ch := range chars {
		// Key down
		ki := keyboardInput{
			Type:    inputKeyboard,
			wScan:   ch,
			dwFlags: keyEventUnicode,
		}
		procSendInput.Call(1, uintptr(unsafe.Pointer(&ki)), unsafe.Sizeof(ki))

		// Key up
		ki.dwFlags = keyEventUnicode | keyEventUp
		procSendInput.Call(1, uintptr(unsafe.Pointer(&ki)), unsafe.Sizeof(ki))
	}
	return nil
}

func charToVK(ch rune) uint16 {
	if ch >= 'a' && ch <= 'z' {
		return uint16(ch - 32) // A=0x41
	}
	if ch >= 'A' && ch <= 'Z' {
		return uint16(ch)
	}
	if ch >= '0' && ch <= '9' {
		return uint16(ch)
	}
	switch ch {
	case ' ':
		return 0x20
	case '\n', '\r':
		return 0x0D
	case '\t':
		return 0x09
	}
	return 0
}

func needsShift(ch rune) bool {
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	// Special characters that need shift
	return strings.ContainsRune("!@#$%^&*()_+{}|:\"<>?~", ch)
}

// keyToVK maps key names to Windows Virtual Key codes
func keyToVK(key string) uint16 {
	switch strings.ToLower(key) {
	case "ctrl", "control":
		return 0x11
	case "alt":
		return 0x12
	case "shift":
		return 0x10
	case "win", "super", "command":
		return 0x5B
	case "enter", "return":
		return 0x0D
	case "tab":
		return 0x09
	case "escape", "esc":
		return 0x1B
	case "space":
		return 0x20
	case "backspace":
		return 0x08
	case "delete", "del":
		return 0x2E
	case "home":
		return 0x24
	case "end":
		return 0x23
	case "pageup":
		return 0x21
	case "pagedown":
		return 0x22
	case "up":
		return 0x26
	case "down":
		return 0x28
	case "left":
		return 0x25
	case "right":
		return 0x27
	case "f1":
		return 0x70
	case "f2":
		return 0x71
	case "f3":
		return 0x72
	case "f4":
		return 0x73
	case "f5":
		return 0x74
	case "f6":
		return 0x75
	case "f7":
		return 0x76
	case "f8":
		return 0x77
	case "f9":
		return 0x78
	case "f10":
		return 0x79
	case "f11":
		return 0x7A
	case "f12":
		return 0x7B
	default:
		if len(key) == 1 {
			ch := rune(key[0])
			return charToVK(ch)
		}
		return 0
	}
}

// NewPlatformKeyboard returns the platform keyboard implementation
func NewPlatformKeyboard() KeyboardController { return NewWinKeyboard() }
