package modules

import (
	"net"
	"strings"
)

// DNSSweep looks up MX, NS, TXT, and CNAME records for the domain, and pulls
// out the SPF and DMARC policies from the TXT records if present. SPF/DMARC
// tell you how well the domain is protected against email spoofing.
func DNSSweep(domain string, result *DomainResult) {
	if mxRecords, err := net.LookupMX(domain); err == nil {
		for _, mx := range mxRecords {
			// A lone "." host is an explicit RFC 7505 null MX record
			// ("this domain accepts no mail"), not a real mail server.
			if mx.Host == "." {
				continue
			}
			result.MXRecords = append(result.MXRecords, strings.TrimSuffix(mx.Host, "."))
		}
	}

	if nsRecords, err := net.LookupNS(domain); err == nil {
		for _, ns := range nsRecords {
			result.NSRecords = append(result.NSRecords, strings.TrimSuffix(ns.Host, "."))
		}
	}

	if cname, err := net.LookupCNAME(domain); err == nil {
		cname = strings.TrimSuffix(cname, ".")
		if cname != domain+"." && cname != domain {
			result.CNAME = cname
		}
	}

	if txtRecords, err := net.LookupTXT(domain); err == nil {
		result.TXTRecords = txtRecords
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				result.SPFRecord = txt
			}
		}
	}

	if dmarcRecords, err := net.LookupTXT("_dmarc." + domain); err == nil {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(txt, "v=DMARC1") {
				result.DMARCRecord = txt
			}
		}
	}
}
