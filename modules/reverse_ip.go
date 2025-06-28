package modules

import (
    "fmt"
    "net"
)

func ReverseIP(domain string) {
    fmt.Println("\n[+] Reverse IP Lookup (basic simulation):")
    ips, err := net.LookupIP(domain)
    if err != nil || len(ips) == 0 {
        fmt.Println("Could not resolve domain.")
        return
    }
    fmt.Println("IP:", ips[0].String())
    fmt.Println("Note: Use Shodan API or external service for full reverse IP results.")
}
