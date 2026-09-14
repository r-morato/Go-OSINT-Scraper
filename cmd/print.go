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

func printDomainResult(d modules.DomainResult) {
	title.Printf("\n=== Domain: %s ===\n", d.Domain)

	section("DNS Records")
	if d.DNSError != "" {
		bad.Println("  " + d.DNSError)
	}
	for _, ip := range d.DNSRecords {
		ok.Println("  " + ip)
	}

	section("Reverse DNS")
	if d.ReverseErr != "" {
		warn.Println("  " + d.ReverseErr)
	}
	for _, entry := range d.ReverseDNS {
		ok.Println("  " + entry)
	}

	section("Subdomains")
	found := 0
	for _, s := range d.Subdomains {
		if s.Found {
			ok.Println("  " + s.Name)
			found++
		}
	}
	if found == 0 {
		bad.Println("  none found")
	}

	section("Emails")
	if len(d.Emails) == 0 {
		bad.Println("  none found")
	}
	for _, e := range d.Emails {
		ok.Println("  " + e)
	}

	section("Tech Fingerprint (HTTP headers)")
	if d.FPError != "" {
		bad.Println("  " + d.FPError)
	}
	for k, v := range d.Headers {
		fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
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
			bad.Printf("  %-8s error: %s\n", p.Platform, p.Error)
		case p.Found:
			ok.Printf("  %-8s found     %s\n", p.Platform, p.URL)
		default:
			bad.Printf("  %-8s not found  %s\n", p.Platform, p.URL)
		}
	}
}
