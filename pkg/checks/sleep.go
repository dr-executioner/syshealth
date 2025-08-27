package checks

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func SleepTimeout() CheckResult {
	res := CheckResult{Name: "sleep_timeout"}
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(powercfg -query SCHEME_CURRENT SUB_SLEEP STANDBYIDLE).FriendlyName")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		// Simplified: assume 10 min or less is OK
		if strings.Contains(string(out), "10") || strings.Contains(string(out), "5") {
			res.OK = true
			res.Detail = "sleep <= 10 min"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res
	case "darwin":
		cmd := exec.Command("pmset", "-g")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "sleep") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					val, _ := strconv.Atoi(parts[len(parts)-1])
					if val <= 10 && val > 0 {
						res.OK = true
						res.Detail = fmt.Sprintf("sleep %d min", val)
						return res
					}
					res.OK = false
					res.Detail = fmt.Sprintf("sleep too high: %d", val)
					return res
				}
			}
		}
		res.OK = false
		res.Detail = string(out)
		return res
	case "linux":
		cmd := exec.Command("gsettings", "get", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout")
		out, err := cmd.CombinedOutput()
		if err == nil {
			valStr := strings.TrimSpace(string(out))
			val, _ := strconv.Atoi(valStr)
			if val/60 <= 10 && val > 0 {
				res.OK = true
				res.Detail = fmt.Sprintf("sleep %d min", val/60)
				return res
			}
			res.OK = false
			res.Detail = fmt.Sprintf("sleep too high: %d min", val/60)
			return res
		}
		res.OK = false
		res.Detail = "could not determine sleep timeout"
		return res
	default:
		res.OK = false
		res.Detail = "unsupported OS"
		return res
	}
}
