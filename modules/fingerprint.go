package modules

import (
    "fmt"
    "net/http"
)

func TechFingerprint(domain string) {
    fmt.Println("\n[+] Technology Fingerprint (headers):")
    url := "http://" + domain
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Request failed:", err)
        return
    }
    defer resp.Body.Close()

    for k, v := range resp.Header {
        fmt.Printf("%s: %s\n", k, v)
    }
}
