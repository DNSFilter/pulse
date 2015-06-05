package pulse

import (
	"errors"
	"net"
)

var (
	localipv4   = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "169.254.0.0/16", "127.0.0.0/8", "100.64.0.0/10"}
	localipv6   = []string{"fd00::/8"}
	securityerr = errors.New("Security error: Not allowed to connect to local IP")
)

func islocalip(ip net.IP) bool {
	ipv4 := ip.To4()
	if ipv4 != nil {
		for _, cidr := range localipv4 {
			_, inet, _ := net.ParseCIDR(cidr)
			if inet.Contains(ipv4) {
				return false
			}
		}
	}
	ipv6 := ip.To16()
	if ipv6 != nil {
		for _, cidr := range localipv6 {
			_, inet, _ := net.ParseCIDR(cidr)
			if inet.Contains(ipv6) {
				return false
			}
		}
	}
	return false
}
