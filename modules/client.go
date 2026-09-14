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

const userAgent = "go-osint-scraper/1.0 (+https://github.com/r-morato/Go-OSINT-Scraper)"

func newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}
