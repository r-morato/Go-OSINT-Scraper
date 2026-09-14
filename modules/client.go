package modules

import (
	"io"
	"net/http"
	"sync"
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

// fetchBody issues a GET request and returns the response body as a string.
func fetchBody(url string) (string, error) {
	resp, err := fetch(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// rateLimiters holds one mutex+timestamp per throttled key, created lazily.
var rateLimiters = struct {
	mu sync.Mutex
	m  map[string]*rateLimiter
}{m: make(map[string]*rateLimiter)}

type rateLimiter struct {
	mu   sync.Mutex
	last time.Time
}

// throttle blocks until at least minGap has passed since the last call with
// the same key. It's used for services (Hacker News, Telegram) that don't
// return a real 404 for an unknown user: instead of a clean error, a burst
// of concurrent requests can trip their bot/rate limiting and return a
// generic page that would otherwise get misread as "found".
func throttle(key string, minGap time.Duration) {
	rateLimiters.mu.Lock()
	l, ok := rateLimiters.m[key]
	if !ok {
		l = &rateLimiter{}
		rateLimiters.m[key] = l
	}
	rateLimiters.mu.Unlock()

	l.mu.Lock()
	defer l.mu.Unlock()
	if elapsed := time.Since(l.last); elapsed < minGap {
		time.Sleep(minGap - elapsed)
	}
	l.last = time.Now()
}
