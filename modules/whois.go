package modules

import (
	"os/exec"
	"runtime"
)

// GetWHOIS retrieves domain registration info via the system's whois binary.
func GetWHOIS(domain string, result *DomainResult) {
	whoisCmd := "whois"
	if runtime.GOOS == "windows" {
		whoisCmd = "whois.exe"
	}

	out, err := exec.Command(whoisCmd, domain).CombinedOutput()
	if err != nil {
		result.WHOISError = err.Error()
		return
	}
	result.WHOIS = string(out)
}
