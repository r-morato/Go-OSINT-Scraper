package modules

// DomainResult holds every finding gathered while scanning a domain.
type DomainResult struct {
	Domain string

	DNSRecords []string
	DNSError   string

	WHOIS      string
	WHOISError string

	ReverseDNS []string
	ReverseErr string

	Subdomains []SubdomainCheck

	Emails []string

	Headers map[string][]string
	FPError string
}

// SubdomainCheck is the outcome of probing a single candidate subdomain.
type SubdomainCheck struct {
	Name  string
	Found bool
}

// UsernameResult holds every finding gathered while checking a username.
type UsernameResult struct {
	Username  string
	Platforms []PlatformCheck
}

// PlatformCheck is the outcome of checking a username on one platform.
type PlatformCheck struct {
	Platform string
	URL      string
	Found    bool
	Error    string
}
