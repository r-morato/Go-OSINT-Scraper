package modules

import (
	"net/http"
	"time"
)

// httpClient is shared by every module that makes HTTP requests so scans
// never hang forever on a slow or unresponsive host.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// slowClient is for known-slow public APIs (crt.sh, the Wayback Machine)
// that routinely take longer than a normal request but are still worth
// waiting on since there's no alternative source for that data.
var slowClient = &http.Client{
	Timeout: 25 * time.Second,
}

const userAgent = "go-osint-scraper/1.0 (+https://github.com/r-morato/Go-OSINT-Scraper)"

func newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

// fetch issues a GET request through the shared client.
func fetch(url string) (*http.Response, error) {
	req, err := newRequest("GET", url)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

// fetchSlow is like fetch but allows the longer timeout of slowClient.
func fetchSlow(url string) (*http.Response, error) {
	req, err := newRequest("GET", url)
	if err != nil {
		return nil, err
	}
	return slowClient.Do(req)
}
