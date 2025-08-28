package checks

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type UpdateDetail struct {
	UpdatesAvailable bool     `json:"updates_available"`
	Count            int      `json:"count"`
	Packages         []string `json:"packages"`
}

func OSUpdates() CheckResult {
	switch runtime.GOOS {
	case "linux":
		content, err := os.ReadFile("/etc/os-release")
		if err != nil {
			return failDetail("os_updates", UpdateDetail{}, "cannot detect distro: "+err.Error())
		}
		osrelease := string(content)

		if strings.Contains(osrelease, "Ubuntu") || strings.Contains(osrelease, "Debian") {
			return parseUpdates("apt-get", "-s", "upgrade")
		} else if strings.Contains(osrelease, "Fedora") || strings.Contains(osrelease, "Red Hat") {
			return parseUpdates("dnf", "check-update")
		} else if strings.Contains(osrelease, "Arch") {
			return parseUpdates("checkupdates")
		} else {
			return failDetail("os_updates", UpdateDetail{}, "unsupported Linux distro")
		}

	case "darwin":
		return parseUpdates("softwareupdate", "-l")

	case "windows":
		return parseUpdates("powershell", "winget", "upgrade", "--accept-source-agreements")

	default:
		return failDetail("os_updates", UpdateDetail{}, "unsupported OS: "+runtime.GOOS)
	}
}

func parseUpdates(cmd string, args ...string) CheckResult {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()

	output := out.String()
	if err != nil && output == "" {
		return failDetail("os_updates", UpdateDetail{}, "error running "+cmd+": "+err.Error())
	}

	// Very naive parsing: split lines into "packages"
	lines := strings.Split(output, "\n")
	var pkgs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// heuristics: look for package names depending on command
		if cmd == "apt-get" && strings.Contains(line, "Inst ") {
			// e.g., "Inst firefox [version]"
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pkgs = append(pkgs, fields[1])
			}
		} else if cmd == "dnf" || cmd == "checkupdates" {
			// usually: "pkgname version"
			fields := strings.Fields(line)
			if len(fields) > 0 {
				pkgs = append(pkgs, fields[0])
			}
		} else if cmd == "softwareupdate" && strings.Contains(line, "*") {
			// macOS lines with *
			pkgs = append(pkgs, strings.TrimSpace(strings.TrimPrefix(line, "*")))
		} else if cmd == "powershell" {
			// winget upgrade outputs: "Name Id Version Available ..."
			if !strings.Contains(line, "------") && !strings.HasPrefix(line, "Name") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					pkgs = append(pkgs, fields[0])
				}
			}
		}
	}

	detail := UpdateDetail{
		UpdatesAvailable: len(pkgs) > 0,
		Count:            len(pkgs),
		Packages:         pkgs,
	}

	// JSON-encode detail
	detailJSON, _ := json.Marshal(detail)

	return CheckResult{
		Name:   "os_updates",
		OK:     !detail.UpdatesAvailable,
		Detail: string(detailJSON),
	}
}

func failDetail(name string, d UpdateDetail, msg string) CheckResult {
	detail := UpdateDetail{
		UpdatesAvailable: false,
		Count:            0,
		Packages:         []string{},
	}
	detailJSON, _ := json.Marshal(detail)

	return CheckResult{Name: name, OK: false, Detail: msg + " | " + string(detailJSON)}
}
