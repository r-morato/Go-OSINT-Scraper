package modules

import (
    "fmt"
    "os"
)

func ExportResults(domain, username string) {
    file, err := os.Create("osint_report.txt")
    if err != nil {
        fmt.Println("Error writing report:", err)
        return
    }
    defer file.Close()

    file.WriteString("OSINT Report\n")
    file.WriteString("=================\n")
    file.WriteString("Domain: " + domain + "\n")
    file.WriteString("Username: " + username + "\n")
    fmt.Println("\n[+] Results saved to osint_report.txt")
}
