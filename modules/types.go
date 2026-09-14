package modules

// DomainResult holds every finding gathered while scanning a domain.
type DomainResult struct {
	Domain string

	DNSRecords []string
	DNSError   string

	MXRecords  []string
	NSRecords  []string
	TXTRecords []string
	CNAME      string

	SPFRecord   string
	DMARCRecord string

	TLSIssuer string
	TLSExpiry string
	TLSSANs   []string
	TLSError  string

	CTSubdomains []string
	CTError      string

	SecurityGrade   string
	MissingHeaders  []string
	SecHeadersError string

	WaybackFirst string
	WaybackLast  string
	WaybackError string

	GeoIPs     []GeoIPInfo
	GeoIPError string

	WHOIS      string
	WHOISError string

	ReverseDNS []string
	ReverseErr string

	Subdomains []SubdomainCheck

	OpenPorts     []int
	PortScanError string

	Emails        []string
	PhoneNumbers  []string
	CryptoWallets []string
	SocialLinks   []string

	Headers map[string][]string
	FPError string
}

// GeoIPInfo is the geolocation of one resolved IP address.
type GeoIPInfo struct {
	IP      string
	Country string
	City    string
	ISP     string
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
