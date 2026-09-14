package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-osint-scraper/modules"
)

var (
	domain   string
	username string
)

var rootCmd = &cobra.Command{
	Use:   "osint",
	Short: "OSINT scraper for gathering public intel",
	Run: func(cmd *cobra.Command, args []string) {
		var domainResult modules.DomainResult
		var usernameResult modules.UsernameResult

		if domain != "" {
			domainResult.Domain = domain
			modules.LookupDNS(domain, &domainResult)
			modules.GetWHOIS(domain, &domainResult)
			modules.ReverseIP(domain, &domainResult)
			modules.FindSubdomains(domain, &domainResult)
			modules.EmailHarvester(domain, &domainResult)
			modules.TechFingerprint(domain, &domainResult)
		}
		if username != "" {
			usernameResult.Username = username
			modules.CheckUsername(username, &usernameResult)
		}

		if err := modules.ExportResults(domainResult, usernameResult, "osint_report.txt"); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("\n[+] Results saved to osint_report.txt")
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Domain to scan")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "Username to check")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
