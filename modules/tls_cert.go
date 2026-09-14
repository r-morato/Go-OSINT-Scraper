package modules

import (
	"crypto/tls"
	"net"
	"time"
)

// TLSCertInfo connects to the domain on port 443 and reads its TLS
// certificate: who issued it, when it expires, and every hostname it
// covers. The Subject Alternative Names often list extra subdomains that
// aren't found any other way.
func TLSCertInfo(domain string, result *DomainResult) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", domain+":443", &tls.Config{
		ServerName: domain,
	})
	if err != nil {
		result.TLSError = err.Error()
		return
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		result.TLSError = "no certificate presented"
		return
	}

	cert := certs[0]
	result.TLSIssuer = cert.Issuer.CommonName
	result.TLSExpiry = cert.NotAfter.Format("2006-01-02")
	result.TLSSANs = cert.DNSNames
}
