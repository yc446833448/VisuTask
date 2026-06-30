//go:build darwin

package action

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// CheckAccessibility checks if the app has Accessibility permissions granted.
// Returns true if the permission is granted, false otherwise.
func CheckAccessibility() bool {
	// Try a simple System Events operation that requires accessibility permissions
	cmd := exec.Command("osascript", "-e",
		`tell application "System Events" to get name of first process whose background only is false`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	// If output is empty or contains error indicators, permission is likely missing
	out := strings.TrimSpace(string(output))
	if out == "" || strings.Contains(out, "not allowed") ||
		strings.Contains(out, "not permitted") || strings.Contains(out, "assistive") {
		return false
	}
	return true
}

// RequestAccessibility shows a dialog guiding the user to enable Accessibility
// permissions. If the user agrees, it opens System Settings to the Accessibility pane.
// Returns an error if the dialog fails or the user declines.
func RequestAccessibility() error {
	// Show dialog explaining the requirement
	dialogScript := `display dialog "VisuTask 需要「辅助功能」权限才能模拟鼠标点击、键盘输入和窗口控制。

请在打开的「安全性与隐私」设置中：
1. 点击左下角 🔒 解锁
2. 找到 VisuTask 并勾选启用
3. 首次使用请点击 + 号添加 VisuTask

授权后请重启 VisuTask。" buttons {"打开系统设置", "稍后再说"} default button 1 with title "权限请求" with icon caution`

	cmd := exec.Command("osascript", "-e", dialogScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// User clicked "Not Now" or dialog failed
		out := strings.TrimSpace(string(output))
		if strings.Contains(out, "User canceled") {
			return fmt.Errorf("permission dialog dismissed by user")
		}
		// Dialog itself might have failed, but we can still try to open settings
	}

	// Open System Settings to the Accessibility privacy pane
	openSettings := `open "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"`
	exec.Command("sh", "-c", openSettings).Run()

	return nil
}

// EnsureAccessibility checks permissions and prompts the user if not granted.
// This should be called during app startup.
func EnsureAccessibility() {
	if !CheckAccessibility() {
		log.Println("⚠️  Accessibility permission not granted — requesting permission...")
		if err := RequestAccessibility(); err != nil {
			log.Printf("⚠️  Permission request skipped: %v\n", err)
		}
	} else {
		log.Println("✅ Accessibility permission is granted")
	}
}
