package modules

import (
	"fmt"
	"net"
)

func LookupDNS(domain string) {
	fmt.Println("\n[+] DNS Records for:", domain)
	records, err := net.LookupHost(domain)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, ip := range records {
		fmt.Println("-", ip)
	}
}
