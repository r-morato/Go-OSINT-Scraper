package modules

import "net"

// ReverseIP resolves the domain to its IP addresses and then performs a real
// reverse DNS (PTR) lookup on each one.
func ReverseIP(domain string, result *DomainResult) {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		result.ReverseErr = "could not resolve domain"
		return
	}

	seen := make(map[string]bool)
	for _, ip := range ips {
		names, err := net.LookupAddr(ip.String())
		if err != nil || len(names) == 0 {
			continue
		}
		for _, name := range names {
			entry := ip.String() + " -> " + name
			if !seen[entry] {
				seen[entry] = true
				result.ReverseDNS = append(result.ReverseDNS, entry)
			}
		}
	}

	if len(result.ReverseDNS) == 0 {
		result.ReverseErr = "no PTR records found for resolved IP(s)"
	}
}
