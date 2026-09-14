package modules

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// commonPorts is a short, well-known wordlist rather than a full 1-65535
// sweep, so a scan stays quick and only checks the obvious services.
var commonPorts = []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 993, 995, 3306, 3389, 8080, 8443}

// PortScan does a concurrent TCP connect scan of commonPorts against the
// domain's first resolved IP.
//
// Only scan hosts you're authorized to scan.
func PortScan(domain string, result *DomainResult) {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		result.PortScanError = "could not resolve domain"
		return
	}
	ip := ips[0].String()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, port := range commonPorts {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := net.JoinHostPort(ip, fmt.Sprint(port))
			conn, err := net.DialTimeout("tcp", address, 1500*time.Millisecond)
			if err != nil {
				return
			}
			conn.Close()

			mu.Lock()
			result.OpenPorts = append(result.OpenPorts, port)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	sort.Ints(result.OpenPorts)
}
