package modules

import (
	"net"
	"sort"
	"sync"
)

// candidateSubdomains is a small, common wordlist. Real engagements would use
// a much bigger list or an API; this keeps the tool dependency-free.
var candidateSubdomains = []string{
	"www", "mail", "ftp", "blog", "dev", "api", "staging", "test", "admin", "shop", "vpn", "cdn",
}

// FindSubdomains checks a list of common subdomain candidates and reports
// which ones actually resolve.
func FindSubdomains(domain string, result *DomainResult) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, sub := range candidateSubdomains {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			host := sub + "." + domain
			_, err := net.LookupHost(host)

			mu.Lock()
			result.Subdomains = append(result.Subdomains, SubdomainCheck{
				Name:  host,
				Found: err == nil,
			})
			mu.Unlock()
		}(sub)
	}

	wg.Wait()

	sort.Slice(result.Subdomains, func(i, j int) bool {
		return result.Subdomains[i].Name < result.Subdomains[j].Name
	})
}
