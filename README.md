# Go OSINT Scraper

A command-line tool written in Go for gathering Open Source Intelligence (OSINT) on domains and usernames. Useful for reconnaissance in cybersecurity, penetration testing, and privacy research.

## Features

All domain checks run concurrently, so a full scan takes about as long as the slowest single check:

- **DNS Lookup** – Resolves A/AAAA records for a domain
- **DNS Sweep** – Resolves MX/NS/TXT/CNAME records and parses SPF/DMARC policy from TXT
- **TLS Certificate Inspector** – Reads the site's certificate: issuer, expiry, and every hostname in its Subject Alternative Names
- **WHOIS Lookup** – Retrieves domain registration info (requires `whois` installed)
- **Reverse DNS Lookup** – Resolves the domain's IP(s) and performs a real PTR lookup on each
- **IP Geolocation** – Country/city/ISP for each resolved IP (via the free ip-api.com)
- **Subdomain Enumeration** – Checks a list of common subdomains and reports which actually resolve
- **Certificate Transparency Search** – Queries crt.sh for every subdomain that has ever appeared on a public TLS certificate for the domain
- **Page Harvester** – Extracts emails, phone numbers, crypto wallet addresses, and social media links from the homepage
- **Technology Fingerprinting** – Shows HTTP headers for basic tech insight
- **Security Headers Grade** – Grades the site A–F on recommended security headers (HSTS, CSP, etc.), similar to securityheaders.com
- **Wayback Machine History** – First/last archive.org snapshot for the domain
- **Port Scan** *(opt-in via `--ports`)* – Concurrent TCP connect scan of ~15 common ports
- **Username Check** – Checks whether a profile exists on GitHub, X, and Reddit
- **Export Results** – Saves the full scan output to a report file

## Installation

Ensure you have Go installed (`go version`), then:

```bash
git clone https://github.com/r-morato/Go-OSINT-Scraper.git
cd Go-OSINT-Scraper
go mod tidy
```

## Usage

```bash
go run . --domain example.com --user octocat
```

Include a port scan (only scan hosts you're authorized to scan):

```bash
go run . --domain example.com --ports
```

### Options

- `--domain`, `-d` — Domain to scan (e.g., `example.com`)
- `--user`, `-u` — Username to check across platforms (e.g., `octocat`)
- `--output`, `-o` — File to write the full report to (default: `osint_report.txt`)
- `--ports`, `-p` — Also scan common TCP ports on `--domain`

## Output

Results are printed to the console (with color-coded status) and exported in full to the report file.

> A few caveats worth knowing:
> - **Username checks** rely on HTTP status codes. Some sites render profiles client-side and can return `200 OK` even for a page that doesn't really exist, so treat a "Found" result as a lead worth verifying, not proof.
> - **Phone numbers** are extracted with a regex over page text — treat matches as leads, not verified numbers.
> - **crt.sh and the Wayback Machine** are free, rate-limited third-party services and occasionally time out or are slow; the tool reports this gracefully rather than failing the whole scan.

## Dependencies

- [`github.com/spf13/cobra`](https://github.com/spf13/cobra) — CLI framework
- [`github.com/fatih/color`](https://github.com/fatih/color) — Colored terminal output
- `whois` command must be installed on your system (Linux/macOS: `brew install whois` or `apt install whois`)
