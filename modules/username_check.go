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

func CheckUsernameOtherSites(username string) {
    sites := []string{
        "https://twitter.com/",
        "https://reddit.com/user/",
    }
    fmt.Println("\n[+] Checking username on other platforms:")
    for _, site := range sites {
        url := site + username
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error checking", url)
            continue
        }
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            fmt.Println("Found:", url)
        } else {
            fmt.Println("Not Found:", url)
        }
    }
}
