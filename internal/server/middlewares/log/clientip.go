package log

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

func extractClientIP(request *http.Request) netip.Addr {
	if request == nil {
		return netip.Addr{}
	}

	remoteAddress := removeAllSpaces(request.RemoteAddr)
	xRealIP := removeAllSpaces(request.Header.Get("X-Real-IP"))
	xForwardedFor := request.Header.Values("X-Forwarded-For")
	for i := range xForwardedFor {
		xForwardedFor[i] = removeAllSpaces(xForwardedFor[i])
	}

	// No header so it can only be remoteAddress
	if xRealIP == "" && len(xForwardedFor) == 0 {
		ip, err := getIPFromHostPort(remoteAddress)
		if err == nil {
			return ip
		}
		return netip.Addr{}
	}

	// remoteAddress is the last proxy server forwarding the traffic
	// so we look into the HTTP headers to get the client IP
	xForwardedIPs := parseAllValidIPStrings(xForwardedFor)
	publicXForwardedIPs := extractPublicIPs(xForwardedIPs)
	if len(publicXForwardedIPs) > 0 {
		// first public XForwardedIP should be the client IP
		return publicXForwardedIPs[0]
	}

	// If all forwarded IP addresses are private we use the x-real-ip
	// address if it exists
	if xRealIP != "" {
		ip, err := getIPFromHostPort(xRealIP)
		if err == nil {
			return ip
		}
	}

	// Client IP is the first private IP address in the chain
	return xForwardedIPs[0]
}

func removeAllSpaces(header string) string {
	header = strings.ReplaceAll(header, " ", "")
	header = strings.ReplaceAll(header, "\t", "")
	return header
}

func extractPublicIPs(ips []netip.Addr) (publicIPs []netip.Addr) {
	publicIPs = make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		if ip.IsPrivate() {
			continue
		}
		publicIPs = append(publicIPs, ip)
	}
	return publicIPs
}

func parseAllValidIPStrings(stringIPs []string) (ips []netip.Addr) {
	ips = make([]netip.Addr, 0, len(stringIPs))
	for _, s := range stringIPs {
		ip, err := netip.ParseAddr(s)
		if err == nil {
			ips = append(ips, ip)
		}
	}
	return ips
}

func getIPFromHostPort(address string) (ip netip.Addr, err error) {
	// address can be in the form ipv4:port, ipv6:port, ipv4 or ipv6
	ipString, _, err := splitHostPort(address)
	if err != nil {
		ipString = address
	}
	return netip.ParseAddr(ipString)
}

func splitHostPort(address string) (ip, port string, err error) {
	if strings.ContainsRune(address, '[') && strings.ContainsRune(address, ']') {
		// should be an IPv6 address with brackets
		return net.SplitHostPort(address)
	}
	const ipv4MaxColons = 1
	if strings.Count(address, ":") > ipv4MaxColons {
		// could be an IPv6 without brackets
		i := strings.LastIndex(address, ":")
		port = address[i+1:]
		ip = address[0:i]
		_, err = netip.ParseAddr(ip)
		if err != nil {
			return net.SplitHostPort(address)
		}
		return ip, port, nil
	}
	// IPv4 address
	return net.SplitHostPort(address)
}
