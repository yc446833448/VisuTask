//go:build !darwin

package action

// CheckAccessibility returns true on non-macOS platforms (no special permissions needed).
func CheckAccessibility() bool {
	return true
}

// RequestAccessibility is a no-op on non-macOS platforms.
func RequestAccessibility() error {
	return nil
}

// EnsureAccessibility is a no-op on non-macOS platforms.
func EnsureAccessibility() {
	// No special permissions needed on this platform
}
