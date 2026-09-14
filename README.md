# Go OSINT Scraper

A command-line tool written in Go for gathering Open Source Intelligence (OSINT) on domains and usernames. Useful for reconnaissance in cybersecurity, penetration testing, and privacy research.

## Features

All domain checks run concurrently, so a full scan takes about as long as the slowest single check:

- **DNS Lookup** – Resolves A/AAAA records for a domain
- **WHOIS Lookup** – Retrieves domain registration info (requires `whois` installed)
- **Reverse DNS Lookup** – Resolves the domain's IP(s) and performs a real PTR lookup on each
- **Subdomain Enumeration** – Checks a list of common subdomains and reports which actually resolve
- **Email Harvester** – Extracts email addresses from the homepage
- **Technology Fingerprinting** – Shows HTTP headers for basic tech insight
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

### Options

- `--domain`, `-d` — Domain to scan (e.g., `example.com`)
- `--user`, `-u` — Username to check across platforms (e.g., `octocat`)
- `--output`, `-o` — File to write the full report to (default: `osint_report.txt`)

## Output

Results are printed to the console (with color-coded found/not-found status) and exported in full to the report file.

> Username checks rely on HTTP status codes. Some sites render profiles client-side and can return `200 OK` even for a page that doesn't really exist, so treat a "Found" result as a lead worth verifying, not proof.

## Dependencies

- [`github.com/spf13/cobra`](https://github.com/spf13/cobra) — CLI framework
- [`github.com/fatih/color`](https://github.com/fatih/color) — Colored terminal output
- `whois` command must be installed on your system (Linux/macOS: `brew install whois` or `apt install whois`)
