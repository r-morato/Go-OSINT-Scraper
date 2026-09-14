package cmd

import (
	"fmt"
	"strings"

	"go-osint-scraper/modules"
)

func section(name string) {
	fmt.Println()
	title.Println("[+] " + name)
}

func printList(items []string, emptyMsg string) {
	if len(items) == 0 {
		bad.Println("  " + emptyMsg)
		return
	}
	for _, item := range items {
		ok.Println("  " + item)
	}
}

func printDomainResult(d modules.DomainResult) {
	title.Printf("\n=== Domain: %s ===\n", d.Domain)

	section("DNS Records")
	if d.DNSError != "" {
		bad.Println("  " + d.DNSError)
	}
	printList(d.DNSRecords, "none found")

	section("DNS Sweep (MX / NS / TXT / CNAME)")
	if len(d.MXRecords) > 0 {
		fmt.Println("  MX:", strings.Join(d.MXRecords, ", "))
	}
	if len(d.NSRecords) > 0 {
		fmt.Println("  NS:", strings.Join(d.NSRecords, ", "))
	}
	if d.CNAME != "" {
		fmt.Println("  CNAME:", d.CNAME)
	}
	if d.SPFRecord != "" {
		ok.Println("  SPF: " + d.SPFRecord)
	} else {
		warn.Println("  SPF: not configured (domain is easier to spoof in email)")
	}
	if d.DMARCRecord != "" {
		ok.Println("  DMARC: " + d.DMARCRecord)
	} else {
		warn.Println("  DMARC: not configured (domain is easier to spoof in email)")
	}

	section("TLS Certificate")
	if d.TLSError != "" {
		bad.Println("  " + d.TLSError)
	} else {
		fmt.Println("  Issuer:", d.TLSIssuer)
		fmt.Println("  Expires:", d.TLSExpiry)
		if len(d.TLSSANs) > 0 {
			fmt.Println("  SANs:", strings.Join(d.TLSSANs, ", "))
		}
	}

	section("Reverse DNS")
	if d.ReverseErr != "" {
		warn.Println("  " + d.ReverseErr)
	}
	printList(d.ReverseDNS, "none found")

	section("IP Geolocation")
	if d.GeoIPError != "" {
		warn.Println("  " + d.GeoIPError)
	}
	for _, g := range d.GeoIPs {
		fmt.Printf("  %s — %s, %s (%s)\n", g.IP, g.City, g.Country, g.ISP)
	}

	section("Subdomains (DNS probe)")
	probed := 0
	for _, s := range d.Subdomains {
		if s.Found {
			ok.Println("  " + s.Name)
			probed++
		}
	}
	if probed == 0 {
		bad.Println("  none found")
	}

	section("Subdomains (crt.sh certificate transparency)")
	if d.CTError != "" {
		warn.Println("  " + d.CTError)
	}
	printList(d.CTSubdomains, "none found")

	section("Emails")
	printList(d.Emails, "none found")

	section("Phone Numbers")
	printList(d.PhoneNumbers, "none found")
	if len(d.PhoneNumbers) > 0 {
		warn.Println("  (regex-based — verify manually, false positives are common)")
	}

	section("Crypto Wallets")
	printList(d.CryptoWallets, "none found")

	section("Social Links")
	printList(d.SocialLinks, "none found")

	section("Tech Fingerprint (HTTP headers)")
	if d.FPError != "" {
		bad.Println("  " + d.FPError)
	}
	for k, v := range d.Headers {
		fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
	}

	section("Security Headers Grade")
	if d.SecHeadersError != "" {
		bad.Println("  " + d.SecHeadersError)
	} else {
		gradeColor := ok
		if d.SecurityGrade == "D" || d.SecurityGrade == "F" {
			gradeColor = bad
		} else if d.SecurityGrade != "A" {
			gradeColor = warn
		}
		gradeColor.Println("  Grade: " + d.SecurityGrade)
		if len(d.MissingHeaders) > 0 {
			fmt.Println("  Missing:", strings.Join(d.MissingHeaders, ", "))
		}
	}

	section("Wayback Machine History")
	if d.WaybackError != "" {
		warn.Println("  " + d.WaybackError)
	} else {
		fmt.Println("  First seen:", d.WaybackFirst)
		fmt.Println("  Last seen:", d.WaybackLast)
	}

	if scanPorts {
		section("Open Ports (common TCP ports)")
		if d.PortScanError != "" {
			bad.Println("  " + d.PortScanError)
		} else if len(d.OpenPorts) == 0 {
			bad.Println("  none open")
		} else {
			for _, port := range d.OpenPorts {
				ok.Printf("  %d\n", port)
			}
		}
	}

	section("WHOIS")
	if d.WHOISError != "" {
		bad.Println("  " + d.WHOISError)
	} else {
		fmt.Println(d.WHOIS)
	}
}

func printUsernameResult(u modules.UsernameResult) {
	title.Printf("\n=== Username: %s ===\n", u.Username)
	for _, p := range u.Platforms {
		switch {
		case p.Error != "":
			bad.Printf("  %-12s error: %s\n", p.Platform, p.Error)
		case p.Found:
			ok.Printf("  %-12s found     %s\n", p.Platform, p.URL)
		default:
			bad.Printf("  %-12s not found  %s\n", p.Platform, p.URL)
		}
	}

	if len(u.Variations) > 0 {
		section(fmt.Sprintf("Variations checked (%d)", len(u.Variations)))
		if len(u.VariationMatches) == 0 {
			bad.Println("  no matches found")
		} else {
			for _, m := range u.VariationMatches {
				ok.Printf("  %-20s %-12s %s\n", m.Username, m.Platform, m.URL)
			}
		}
	}
}
