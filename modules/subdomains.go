package modules

import (
    "fmt"
)

func FindSubdomains(domain string) {
    fmt.Println("\n[+] Simulated Subdomain Finder:")
    sampleSubs := []string{"www", "mail", "blog", "dev"}
    for _, sub := range sampleSubs {
        fmt.Printf("- %s.%s\n", sub, domain)
    }
    fmt.Println("Note: Use DNS brute-forcing or APIs for real results.")
}
