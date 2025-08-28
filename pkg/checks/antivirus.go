package checks

import (
	"runtime"
	util "syshealth/internal/utils"
	"syshealth/pkg/structs"
)

func Antivirus() structs.CheckResult {
	switch runtime.GOOS {
	case "windows":
		return util.RunCheck(
			"antivirus",
			"powershell",
			[]string{"-Command", "Get-MpComputerStatus | Select-Object -Property AMServiceEnabled"},
			[]string{"True"},
			"Windows Defender active",
			"Windows Defender inactive",
		)
	case "darwin":
		return util.RunCheck(
			"antivirus",
			"defaults",
			[]string{"read", "/System/Library/CoreServices/XProtect.bundle/Contents/Info", "CFBundleShortVersionString"},
			[]string{""}, // any output means present
			"XProtect present",
			"no AV detected",
		)
	case "linux":
		return util.RunCheck(
			"antivirus",
			"systemctl",
			[]string{"is-active", "clamav-daemon"},
			[]string{"active"},
			"ClamAV active",
			"no active antivirus service",
		)
	default:
		return structs.CheckResult{Name: "antivirus", OK: false, Detail: "unsupported OS"}
	}
}
