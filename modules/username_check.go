package modules

import (
    "fmt"
    "net/http"
)

func CheckUsername(username string) {
    fmt.Println("\n[+] Checking username on GitHub:")
    url := fmt.Sprintf("https://github.com/%s", username)
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusOK {
        fmt.Println("Found:", url)
    } else {
        fmt.Println("Not Found:", url)
    }
}
