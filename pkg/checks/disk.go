package checks

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// CheckResult contains name, status and optional info
type CheckResult struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

func DiskEncryption() CheckResult {
	res := CheckResult{Name: "disk_encryption"}
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("manage-bde", "-status", "C:")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "Percentage Encrypted:") && strings.Contains(string(out), "100%") {
			res.OK = true
			res.Detail = "BitLocker 100%"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res

	case "darwin":
		cmd := exec.Command("fdesetup", "status")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "FileVault is On.") {
			res.OK = true
			res.Detail = "FileVault On"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res

	case "linux":
		// Simplified: detect "crypt" mapper or luks
		cmd := exec.Command("lsblk", "-o", "NAME,TYPE")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "crypt") || strings.Contains(strings.ToLower(string(out)), "luks") {
			res.OK = true
			res.Detail = "LUKS/crypt mapper found"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res
	default:
		res.OK = false
		res.Detail = "unsupported OS"
		return res
	}
}
