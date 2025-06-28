package modules

import (
    "fmt"
    "os/exec"
    "runtime"
)

func GetWHOIS(domain string) {
    fmt.Println("\n[+] WHOIS info for:", domain)

    var whoisCmd string
    if runtime.GOOS == "windows" {
        whoisCmd = "whois.exe" // Assumes whois.exe is available in PATH on Windows
    } else {
        whoisCmd = "whois"
    }

    out, err := exec.Command(whoisCmd, domain).CombinedOutput()
    if err != nil {
        fmt.Println("Error running whois command:", err)
        return
    }
    fmt.Println(string(out))
}
