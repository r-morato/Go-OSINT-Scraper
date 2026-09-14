package modules

import (
	"encoding/json"
	"net"
)

type ipAPIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Country string `json:"country"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
	Query   string `json:"query"`
}

// GeoIPLookup resolves the domain and geolocates each unique IP using the
// free ip-api.com service (no API key required).
func GeoIPLookup(domain string, result *DomainResult) {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		result.GeoIPError = "could not resolve domain"
		return
	}

	seen := make(map[string]bool)
	for _, ip := range ips {
		ipStr := ip.String()
		if seen[ipStr] {
			continue
		}
		seen[ipStr] = true

		info, err := lookupIPInfo(ipStr)
		if err != nil {
			continue
		}
		result.GeoIPs = append(result.GeoIPs, info)
	}

	if len(result.GeoIPs) == 0 && result.GeoIPError == "" {
		result.GeoIPError = "no geolocation data available"
	}
}

func lookupIPInfo(ip string) (GeoIPInfo, error) {
	resp, err := fetch("http://ip-api.com/json/" + ip + "?fields=status,message,country,city,isp,query")
	if err != nil {
		return GeoIPInfo{}, err
	}
	defer resp.Body.Close()

	var data ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return GeoIPInfo{}, err
	}
	if data.Status != "success" {
		return GeoIPInfo{}, nil
	}

	return GeoIPInfo{
		IP:      data.Query,
		Country: data.Country,
		City:    data.City,
		ISP:     data.ISP,
	}, nil
}
