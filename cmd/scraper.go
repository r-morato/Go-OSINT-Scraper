package cmd

import (
    "fmt"
    "go-osint-scraper/modules"
    "github.com/spf13/cobra"
)

var (
    domain   string
    username string
)

var rootCmd = &cobra.Command{
    Use:   "osint",
    Short: "OSINT scraper for gathering public intel",
    Run: func(cmd *cobra.Command, args []string) {
        if domain != "" {
            modules.LookupDNS(domain)
            modules.GetWHOIS(domain)
            modules.ReverseIP(domain)
            modules.FindSubdomains(domain)
            modules.EmailHarvester(domain)
            modules.TechFingerprint(domain)
        }
        if username != "" {
            modules.CheckUsername(username)
            modules.CheckUsernameOtherSites(username)
        }
        modules.ExportResults(domain, username)
    },
}

func Execute() {
    rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Domain to scan")
    rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "Username to check")
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("Error:", err)
    }
}
