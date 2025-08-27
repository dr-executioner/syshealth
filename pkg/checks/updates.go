package checks

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func OSUpdates() CheckResult {
	res := CheckResult{Name: "os_updates"}
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject Microsoft.Update.Session).CreateUpdateSearcher().Search(\"IsInstalled=0\").Updates.Count")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.TrimSpace(string(out)) == "0" {
			res.OK = true
			res.Detail = "no pending updates"
			return res
		}
		res.OK = false
		res.Detail = fmt.Sprintf("pending updates: %s", strings.TrimSpace(string(out)))
		return res

	case "darwin":
		cmd := exec.Command("softwareupdate", "-l")
		out, err := cmd.CombinedOutput()
		if err != nil {
			// softwareupdate may return non-zero when no updates
			if strings.Contains(string(out), "No new software available.") {
				res.OK = true
				res.Detail = "no updates"
				return res
			}
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "No new software available.") {
			res.OK = true
			res.Detail = "no updates"
			return res
		}
		res.OK = false
		res.Detail = string(out)
		return res

	case "linux":
		cmd := exec.Command("apt-get", "-s", "upgrade")
		out, err := cmd.CombinedOutput()
		if err != nil {
			res.OK = false
			res.Detail = fmt.Sprintf("error: %v; output: %s", err, string(out))
			return res
		}
		if strings.Contains(string(out), "0 upgraded, 0 newly installed") {
			res.OK = true
			res.Detail = "no updates (apt)"
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
