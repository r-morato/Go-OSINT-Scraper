package modules

import (
	"io"
	"regexp"
)

var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

// EmailHarvester fetches the domain's homepage and extracts email addresses.
func EmailHarvester(domain string, result *DomainResult) {
	req, err := newRequest("GET", "http://"+domain)
	if err != nil {
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	seen := make(map[string]bool)
	for _, email := range emailRegex.FindAllString(string(body), -1) {
		if !seen[email] {
			seen[email] = true
			result.Emails = append(result.Emails, email)
		}
	}
}
