package modules

import (
	"io"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Requires a country code or at least two separators between digit
	// groups, so it doesn't light up on SVG path data, CSS values, or other
	// long digit runs found in raw HTML/JS/CSS.
	phoneRegex   = regexp.MustCompile(`\+\d{1,3}[-.\s]\(?\d{2,4}\)?[-.\s]\d{3,4}[-.\s]?\d{3,4}|\(\d{2,4}\)[-.\s]\d{3,4}[-.\s]?\d{3,4}`)
	telHrefRegex = regexp.MustCompile(`href="tel:([^"]+)"`)

	btcRegex    = regexp.MustCompile(`\b(bc1[a-z0-9]{25,39}|[13][a-km-zA-HJ-NP-Z1-9]{25,34})\b`)
	ethRegex    = regexp.MustCompile(`\b0x[a-fA-F0-9]{40}\b`)
	socialRegex = regexp.MustCompile(`href="(https?://(?:www\.)?(?:twitter\.com|x\.com|linkedin\.com/in|instagram\.com|facebook\.com|youtube\.com/(?:c/|channel/|@))[^"]+)"`)

	// Elements whose content is never human-readable text, stripped before
	// harvesting so they don't get mistaken for phone numbers, emails, etc.
	scriptTagRegex = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	styleTagRegex  = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	svgTagRegex    = regexp.MustCompile(`(?is)<svg\b[^>]*>.*?</svg>`)
)

// HarvestPage fetches the domain's homepage once and pulls out every email
// address, phone number, crypto wallet address, and social media profile
// link it can find on the page.
func HarvestPage(domain string, result *DomainResult) {
	resp, err := fetch("http://" + domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := stripNonTextElements(string(body))

	result.Emails = dedupe(emailRegex.FindAllString(html, -1))

	phones := phoneRegex.FindAllString(html, -1)
	for _, m := range telHrefRegex.FindAllStringSubmatch(html, -1) {
		phones = append(phones, m[1])
	}
	result.PhoneNumbers = dedupe(phones)

	result.SocialLinks = dedupe(firstSubmatches(socialRegex.FindAllStringSubmatch(html, -1)))

	wallets := append([]string{}, btcRegex.FindAllString(html, -1)...)
	wallets = append(wallets, ethRegex.FindAllString(html, -1)...)
	result.CryptoWallets = dedupe(wallets)
}

// stripNonTextElements removes <script>, <style>, and <svg> blocks (tags and
// content) since they're never real page copy and are a major source of
// false positives for the regexes above.
func stripNonTextElements(html string) string {
	html = scriptTagRegex.ReplaceAllString(html, "")
	html = styleTagRegex.ReplaceAllString(html, "")
	html = svgTagRegex.ReplaceAllString(html, "")
	return html
}

func firstSubmatches(matches [][]string) []string {
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			out = append(out, m[1])
		}
	}
	return out
}

func dedupe(items []string) []string {
	seen := make(map[string]bool, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	return out
}
