package modules

// TechFingerprint records the HTTP response headers for a domain, which can
// hint at the server software, framework, or CDN in use.
func TechFingerprint(domain string, result *DomainResult) {
	resp, err := fetch("http://" + domain)
	if err != nil {
		result.FPError = err.Error()
		return
	}
	defer resp.Body.Close()

	result.Headers = resp.Header
}
