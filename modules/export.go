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

func writeDomainSection(b *strings.Builder, d DomainResult) {
	b.WriteString("\nDomain: " + d.Domain + "\n")
	b.WriteString("-----------------\n")

	b.WriteString("\nDNS Records:\n")
	if d.DNSError != "" {
		b.WriteString("  Error: " + d.DNSError + "\n")
	}
	for _, ip := range d.DNSRecords {
		b.WriteString("  - " + ip + "\n")
	}

	b.WriteString("\nReverse DNS:\n")
	if d.ReverseErr != "" {
		b.WriteString("  Note: " + d.ReverseErr + "\n")
	}
	for _, entry := range d.ReverseDNS {
		b.WriteString("  - " + entry + "\n")
	}

	b.WriteString("\nSubdomains found:\n")
	found := 0
	for _, s := range d.Subdomains {
		if s.Found {
			b.WriteString("  - " + s.Name + "\n")
			found++
		}
	}
	if found == 0 {
		b.WriteString("  (none)\n")
	}

	b.WriteString("\nEmails found:\n")
	if len(d.Emails) == 0 {
		b.WriteString("  (none)\n")
	}
	for _, e := range d.Emails {
		b.WriteString("  - " + e + "\n")
	}

	b.WriteString("\nHTTP Headers:\n")
	if d.FPError != "" {
		b.WriteString("  Error: " + d.FPError + "\n")
	}
	for k, v := range d.Headers {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
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
		b.WriteString(fmt.Sprintf("  %-8s %-10s %s\n", p.Platform, status, p.URL))
	}
}
