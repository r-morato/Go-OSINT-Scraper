package modules

import (
	"encoding/json"
	"strconv"
	"sync"
)

// WaybackHistory checks the Internet Archive's Wayback Machine for the
// earliest and most recent snapshot of the domain, which gives a rough
// sense of how long it's been publicly online. The two lookups run
// concurrently since each can take several seconds on its own.
func WaybackHistory(domain string, result *DomainResult) {
	var first, last string
	var firstErr, lastErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		first, firstErr = waybackTimestamp(domain, 1)
	}()
	go func() {
		defer wg.Done()
		last, lastErr = waybackTimestamp(domain, -1)
	}()
	wg.Wait()

	if firstErr != nil {
		result.WaybackError = firstErr.Error()
		return
	}
	if lastErr != nil {
		result.WaybackError = lastErr.Error()
		return
	}
	if first == "" {
		result.WaybackError = "no snapshots found"
		return
	}

	result.WaybackFirst = formatWaybackTimestamp(first)
	result.WaybackLast = formatWaybackTimestamp(last)
}

func waybackTimestamp(domain string, limit int) (string, error) {
	url := "https://web.archive.org/cdx/search/cdx?url=" + domain +
		"&output=json&fl=timestamp&filter=statuscode:200&limit=" + strconv.Itoa(limit)

	resp, err := fetchSlow(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rows [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return "", err
	}

	// The first row is a header ("timestamp"); the actual data is row 2+.
	if len(rows) < 2 {
		return "", nil
	}
	return rows[1][0], nil
}

func formatWaybackTimestamp(ts string) string {
	if len(ts) < 8 {
		return ts
	}
	return ts[0:4] + "-" + ts[4:6] + "-" + ts[6:8]
}
