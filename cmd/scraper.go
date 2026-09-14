package cmd

import (
	"fmt"
	"sync"

	"go-osint-scraper/modules"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	domain     string
	username   string
	outputPath string
	scanPorts  bool
)

const banner = `
   ___  _____ _____ _  _ _____
  / _ \/ ___/|_   _| \| |_   _|
 | (_) \___ \  | | | .' | | |
  \___/|____/  |_| |_|\_| |_|   go-osint-scraper
`

var (
	title = color.New(color.FgCyan, color.Bold)
	ok    = color.New(color.FgGreen)
	bad   = color.New(color.FgHiBlack)
	warn  = color.New(color.FgYellow)
)

var rootCmd = &cobra.Command{
	Use:   "osint",
	Short: "OSINT scraper for gathering public intel",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" && username == "" {
			fmt.Println("Nothing to do: pass --domain and/or --user. Run with --help for usage.")
			return
		}

		fmt.Print(banner)

		var domainResult modules.DomainResult
		var usernameResult modules.UsernameResult
		var wg sync.WaitGroup

		if domain != "" {
			domainResult.Domain = domain
			wg.Add(1)
			go func() {
				defer wg.Done()
				scanDomain(domain, &domainResult)
			}()
		}
		if username != "" {
			usernameResult.Username = username
			wg.Add(1)
			go func() {
				defer wg.Done()
				scanUsername(username, &usernameResult)
			}()
		}

		wg.Wait()

		if domain != "" {
			printDomainResult(domainResult)
		}
		if username != "" {
			printUsernameResult(usernameResult)
		}

		if err := modules.ExportResults(domainResult, usernameResult, outputPath); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("\n[+] Full results saved to %s\n", outputPath)
	},
}

// scanDomain runs every domain module concurrently and waits for them all to
// finish before returning.
func scanDomain(domain string, result *modules.DomainResult) {
	var wg sync.WaitGroup

	run := func(fn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}

	run(func() { modules.LookupDNS(domain, result) })
	run(func() { modules.GetWHOIS(domain, result) })
	run(func() { modules.ReverseIP(domain, result) })
	run(func() { modules.FindSubdomains(domain, result) })
	run(func() { modules.HarvestPage(domain, result) })
	run(func() { modules.TechFingerprint(domain, result) })
	run(func() { modules.DNSSweep(domain, result) })
	run(func() { modules.TLSCertInfo(domain, result) })
	run(func() { modules.CTSubdomainSearch(domain, result) })
	run(func() { modules.SecurityHeadersGrade(domain, result) })
	run(func() { modules.WaybackHistory(domain, result) })
	run(func() { modules.GeoIPLookup(domain, result) })
	if scanPorts {
		run(func() { modules.PortScan(domain, result) })
	}

	wg.Wait()
}

func scanUsername(username string, result *modules.UsernameResult) {
	modules.CheckUsername(username, result)
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Domain to scan")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "Username to check")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "osint_report.txt", "File to write the full report to")
	rootCmd.PersistentFlags().BoolVarP(&scanPorts, "ports", "p", false, "Also scan common TCP ports on --domain (only use against hosts you're authorized to scan)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
