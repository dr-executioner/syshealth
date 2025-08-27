package checks

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func Antivirus() CheckResult {
	res := CheckResult{Name: "antivirus"}
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "Get-MpComputerStatus | Select-Object -Property AMServiceEnabled")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "True") {
			res.OK = true
			res.Detail = "Windows Defender active"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res
	case "darwin":
		// macOS doesn’t have built-in AV by default, check for XProtect
		cmd := exec.Command("defaults", "read", "/System/Library/CoreServices/XProtect.bundle/Contents/Info", "CFBundleShortVersionString")
		out, err := cmd.CombinedOutput()
		if err == nil && len(strings.TrimSpace(string(out))) > 0 {
			res.OK = true
			res.Detail = "XProtect present"
			return res
		}
		res.OK = false
		res.Detail = "no AV detected"
		return res
	case "linux":
		cmd := exec.Command("systemctl", "is-active", "clamav-daemon")
		out, err := cmd.CombinedOutput()
		if err == nil && strings.Contains(string(out), "active") {
			res.OK = true
			res.Detail = "ClamAV active"
			return res
		}
		res.OK = false
		res.Detail = "no active antivirus service"
		return res
	default:
		res.OK = false
		res.Detail = "unsupported OS"
		return res
	}
}
