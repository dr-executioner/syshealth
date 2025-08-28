package util

import (
	"fmt"
	"os/exec"
	"strings"
	"syshealth/pkg/structs"
)

func RunCheck(name string, cmd string, args []string, successIndicators []string, successMsg string, failMsg string) structs.CheckResult {
	res := structs.CheckResult{Name: name}

	out, err := exec.Command(cmd, args...).CombinedOutput()
	output := string(out)

	if err != nil {
		res.OK = false
		res.Detail = fmt.Sprintf("error: %v; output: %s", err, output)
		return res
	}

	for _, indicator := range successIndicators {
		if strings.Contains(output, indicator) {
			res.OK = true
			res.Detail = successMsg
			return res
		}
	}

	res.OK = false
	if failMsg != "" {
		res.Detail = failMsg
	} else {
		res.Detail = output
	}
	return res
}
