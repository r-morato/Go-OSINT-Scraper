package modules

import "net"

// LookupDNS resolves the A/AAAA records for a domain.
func LookupDNS(domain string, result *DomainResult) {
	records, err := net.LookupHost(domain)
	if err != nil {
		result.DNSError = err.Error()
		return
	}
	result.DNSRecords = records
}
