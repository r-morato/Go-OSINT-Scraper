package modules

// securityHeaders are the response headers checked for, in the order they're
// reported. This mirrors the checks a tool like securityheaders.com runs.
var securityHeaders = []string{
	"Strict-Transport-Security",
	"Content-Security-Policy",
	"X-Frame-Options",
	"X-Content-Type-Options",
	"Referrer-Policy",
	"Permissions-Policy",
}

// SecurityHeadersGrade fetches the site over HTTPS (falling back to HTTP)
// and grades how many recommended security headers it sends back.
func SecurityHeadersGrade(domain string, result *DomainResult) {
	resp, err := fetch("https://" + domain)
	if err != nil {
		resp, err = fetch("http://" + domain)
	}
	if err != nil {
		result.SecHeadersError = err.Error()
		return
	}
	defer resp.Body.Close()

	for _, header := range securityHeaders {
		if resp.Header.Get(header) == "" {
			result.MissingHeaders = append(result.MissingHeaders, header)
		}
	}

	present := len(securityHeaders) - len(result.MissingHeaders)
	switch {
	case present == len(securityHeaders):
		result.SecurityGrade = "A"
	case present >= 5:
		result.SecurityGrade = "B"
	case present >= 4:
		result.SecurityGrade = "C"
	case present >= 3:
		result.SecurityGrade = "D"
	default:
		result.SecurityGrade = "F"
	}
}
