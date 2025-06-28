package modules

import (
    "fmt"
    "net/http"
    "io"
    "regexp"
)

func EmailHarvester(domain string) {
    fmt.Println("\n[+] Searching for email addresses on:", domain)
    resp, err := http.Get("http://" + domain)
    if err != nil {
        fmt.Println("Failed to fetch page:", err)
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading body:", err)
        return
    }

    re := regexp.MustCompile(`[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
    emails := re.FindAllString(string(body), -1)

    if len(emails) == 0 {
        fmt.Println("No emails found.")
    } else {
        for _, email := range emails {
            fmt.Println("Found:", email)
        }
    }
}
