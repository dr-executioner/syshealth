package checks

import (
	"runtime"
	"strings"
	util "syshealth/internal/utils"
	"syshealth/pkg/structs"
)

func DiskEncryption() structs.CheckResult {
	switch runtime.GOOS {
	case "windows":
		return util.RunCheck(
			"disk_encryption",
			"manage-bde",
			[]string{"-status", "C:"},
			[]string{"Percentage Encrypted: 100%"},
			"BitLocker 100%",
			"Disk not fully encrypted",
		)
	case "darwin":
		return util.RunCheck(
			"disk_encryption",
			"fdesetup",
			[]string{"status"},
			[]string{"FileVault is On."},
			"FileVault On",
			"FileVault Off",
		)
	case "linux":
		// Custom: because linux needs "crypt"/"luks" in lsblk output
		res := util.RunCheck(
			"disk_encryption",
			"lsblk",
			[]string{"-o", "NAME,TYPE"},
			[]string{"crypt", "luks"},
			"LUKS/crypt mapper found",
			"no encrypted volumes detected",
		)
		// make case-insensitive check
		if strings.Contains(strings.ToLower(res.Detail), "crypt") {
			res.OK = true
			res.Detail = "LUKS/crypt mapper found"
		}
		return res
	default:
		return structs.CheckResult{Name: "disk_encryption", OK: false, Detail: "unsupported OS"}
	}
}
