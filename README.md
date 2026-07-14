# Go OSINT Scraper

A command-line tool written in Go for gathering Open Source Intelligence (OSINT) on domains and usernames. Useful for reconnaissance in cybersecurity, penetration testing, and privacy research.

##  Features

This tool includes:

-  **DNS Lookup** – Resolves the A records for a domain
-  **WHOIS Lookup** – Retrieves domain registration info (requires `whois` installed)
-  **Reverse IP Lookup** – Gets the resolved IP (basic simulation; use APIs for full data)
-  **Subdomain Enumeration** – Simulated subdomain generation
-  **Email Harvester** – Extracts email addresses from the homepage
-  **Technology Fingerprinting** – Shows HTTP headers for basic tech insight
-  **Username Check (GitHub)** – Checks if a GitHub profile exists
-  **Username Check (Twitter, Reddit)** – Basic existence check on additional platforms
-  **Export Results** – Saves results to `osint_report.txt`

## Installation

Ensure you have Go installed (`go version`), then:

```bash
git clone https://github.com/yourusername/go-osint-scraper.git
cd go-osint-scraper
go mod tidy
````

##  Usage

```bash
go run main.go --domain example.com --user johndoe
```

### Options

* `--domain`, `-d` — Domain to scan (e.g., `example.com`)
* `--user`, `-u` — Username to check across platforms (e.g., `johndoe`)

##  Output

Results are printed to the console and exported to a file called:

```text
osint_report.txt
```

## Dependencies

* [`github.com/spf13/cobra`](https://github.com/spf13/cobra) — CLI framework
* `whois` command must be installed on your system (Linux/macOS: `brew install whois` or `apt install whois`)

