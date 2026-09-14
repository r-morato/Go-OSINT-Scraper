package modules

import (
	"encoding/json"
	"sort"
	"strings"
)

type crtshEntry struct {
	NameValue string `json:"name_value"`
}

// CTSubdomainSearch queries crt.sh's free Certificate Transparency log
// search for every hostname that has ever appeared on a public TLS
// certificate for the domain. This tends to surface real subdomains that a
// wordlist-based guess would never find.
func CTSubdomainSearch(domain string, result *DomainResult) {
	url := "https://crt.sh/?q=%25." + domain + "&output=json"

	resp, err := fetchSlow(url)
	if err != nil {
		result.CTError = err.Error()
		return
	}
	defer resp.Body.Close()

	var entries []crtshEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		result.CTError = "crt.sh returned an unexpected response (it may be rate-limiting or slow right now)"
		return
	}

	seen := make(map[string]bool)
	for _, entry := range entries {
		for _, name := range strings.Split(entry.NameValue, "\n") {
			name = strings.ToLower(strings.TrimSpace(name))
			name = strings.TrimPrefix(name, "*.")
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			result.CTSubdomains = append(result.CTSubdomains, name)
		}
	}

	sort.Strings(result.CTSubdomains)
}
