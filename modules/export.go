package modules

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ExportResults writes every finding gathered during the scan to a
// human-readable report file.
func ExportResults(domain DomainResult, username UsernameResult, outputPath string) error {
	var b strings.Builder

	b.WriteString("OSINT Report\n")
	b.WriteString("=================\n")
	b.WriteString("Generated: " + time.Now().Format(time.RFC1123) + "\n")

	if domain.Domain != "" {
		writeDomainSection(&b, domain)
	}
	if username.Username != "" {
		writeUsernameSection(&b, username)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(b.String()); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	return nil
}

func writeList(b *strings.Builder, heading string, items []string, note string) {
	b.WriteString("\n" + heading + ":\n")
	if note != "" {
		b.WriteString("  Note: " + note + "\n")
	}
	if len(items) == 0 {
		b.WriteString("  (none)\n")
		return
	}
	for _, item := range items {
		b.WriteString("  - " + item + "\n")
	}
}

func orNone(s string) string {
	if s == "" {
		return "(not configured)"
	}
	return s
}

func writeDomainSection(b *strings.Builder, d DomainResult) {
	b.WriteString("\nDomain: " + d.Domain + "\n")
	b.WriteString("-----------------\n")

	writeList(b, "DNS Records", d.DNSRecords, d.DNSError)

	b.WriteString("\nDNS Sweep:\n")
	b.WriteString("  MX: " + strings.Join(d.MXRecords, ", ") + "\n")
	b.WriteString("  NS: " + strings.Join(d.NSRecords, ", ") + "\n")
	b.WriteString("  CNAME: " + d.CNAME + "\n")
	b.WriteString("  SPF: " + orNone(d.SPFRecord) + "\n")
	b.WriteString("  DMARC: " + orNone(d.DMARCRecord) + "\n")

	b.WriteString("\nTLS Certificate:\n")
	if d.TLSError != "" {
		b.WriteString("  Error: " + d.TLSError + "\n")
	} else {
		b.WriteString("  Issuer: " + d.TLSIssuer + "\n")
		b.WriteString("  Expires: " + d.TLSExpiry + "\n")
		b.WriteString("  SANs: " + strings.Join(d.TLSSANs, ", ") + "\n")
	}

	writeList(b, "Reverse DNS", d.ReverseDNS, d.ReverseErr)

	b.WriteString("\nIP Geolocation:\n")
	if d.GeoIPError != "" {
		b.WriteString("  Note: " + d.GeoIPError + "\n")
	}
	for _, g := range d.GeoIPs {
		b.WriteString(fmt.Sprintf("  - %s: %s, %s (%s)\n", g.IP, g.City, g.Country, g.ISP))
	}

	var probed []string
	for _, s := range d.Subdomains {
		if s.Found {
			probed = append(probed, s.Name)
		}
	}
	writeList(b, "Subdomains (DNS probe)", probed, "")
	writeList(b, "Subdomains (crt.sh certificate transparency)", d.CTSubdomains, d.CTError)

	writeList(b, "Emails found", d.Emails, "")
	writeList(b, "Phone numbers found (unverified, regex-based)", d.PhoneNumbers, "")
	writeList(b, "Crypto wallets found", d.CryptoWallets, "")
	writeList(b, "Social links found", d.SocialLinks, "")

	b.WriteString("\nHTTP Headers:\n")
	if d.FPError != "" {
		b.WriteString("  Error: " + d.FPError + "\n")
	}
	for k, v := range d.Headers {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
	}

	b.WriteString("\nSecurity Headers Grade:\n")
	if d.SecHeadersError != "" {
		b.WriteString("  Error: " + d.SecHeadersError + "\n")
	} else {
		b.WriteString("  Grade: " + d.SecurityGrade + "\n")
		if len(d.MissingHeaders) > 0 {
			b.WriteString("  Missing: " + strings.Join(d.MissingHeaders, ", ") + "\n")
		}
	}

	b.WriteString("\nWayback Machine History:\n")
	if d.WaybackError != "" {
		b.WriteString("  Note: " + d.WaybackError + "\n")
	} else {
		b.WriteString("  First seen: " + d.WaybackFirst + "\n")
		b.WriteString("  Last seen: " + d.WaybackLast + "\n")
	}

	if len(d.OpenPorts) > 0 || d.PortScanError != "" {
		b.WriteString("\nOpen Ports:\n")
		if d.PortScanError != "" {
			b.WriteString("  Error: " + d.PortScanError + "\n")
		}
		for _, port := range d.OpenPorts {
			b.WriteString(fmt.Sprintf("  - %d\n", port))
		}
	}

	b.WriteString("\nWHOIS:\n")
	if d.WHOISError != "" {
		b.WriteString("  Error: " + d.WHOISError + "\n")
	} else {
		b.WriteString(d.WHOIS + "\n")
	}
}

func writeUsernameSection(b *strings.Builder, u UsernameResult) {
	b.WriteString("\nUsername: " + u.Username + "\n")
	b.WriteString("-----------------\n")
	for _, p := range u.Platforms {
		status := "Not Found"
		if p.Found {
			status = "Found"
		}
		if p.Error != "" {
			status = "Error: " + p.Error
		}
		b.WriteString(fmt.Sprintf("  %-12s %-10s %s\n", p.Platform, status, p.URL))
	}

	if len(u.Variations) > 0 {
		b.WriteString(fmt.Sprintf("\nVariations checked (%d): %s\n", len(u.Variations), strings.Join(u.Variations, ", ")))
		if len(u.VariationMatches) == 0 {
			b.WriteString("  (no matches found)\n")
		} else {
			for _, m := range u.VariationMatches {
				b.WriteString(fmt.Sprintf("  - %-20s %-12s %s\n", m.Username, m.Platform, m.URL))
			}
		}
	}
}
