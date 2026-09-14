package modules

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
)

// platforms maps a display name to a URL template (%s is the username).
var platforms = map[string]string{
	"GitHub": "https://github.com/%s",
	"X":      "https://x.com/%s",
	"Reddit": "https://reddit.com/user/%s",
}

// CheckUsername probes each known platform for the given username
// concurrently and reports whether a profile appears to exist.
//
// Note: some platforms render profiles client-side and return HTTP 200 even
// for pages that don't exist, so results should be treated as a lead, not
// proof.
func CheckUsername(username string, result *UsernameResult) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, urlTemplate := range platforms {
		wg.Add(1)
		go func(name, urlTemplate string) {
			defer wg.Done()
			check := checkPlatform(name, fmt.Sprintf(urlTemplate, username))

			mu.Lock()
			result.Platforms = append(result.Platforms, check)
			mu.Unlock()
		}(name, urlTemplate)
	}

	wg.Wait()

	sort.Slice(result.Platforms, func(i, j int) bool {
		return result.Platforms[i].Platform < result.Platforms[j].Platform
	})
}

func checkPlatform(name, url string) PlatformCheck {
	check := PlatformCheck{Platform: name, URL: url}

	req, err := newRequest("GET", url)
	if err != nil {
		check.Error = err.Error()
		return check
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		check.Error = err.Error()
		return check
	}
	defer resp.Body.Close()

	check.Found = resp.StatusCode == http.StatusOK
	return check
}
