package validate

import "net"

var NonPublicIPNet = []net.IPNet{
	*ParseIPNet("10.0.0.0/8"),
	*ParseIPNet("172.16.0.0/12"),
	*ParseIPNet("192.168.0.0/16"),
	*ParseIPNet("127.0.0.0/8"),
	*ParseIPNet("0.0.0.0/8"),
	*ParseIPNet("169.254.0.0/16"),
	*ParseIPNet("192.0.0.0/24"),
	*ParseIPNet("192.0.2.0/24"),
	*ParseIPNet("198.51.100.0/24"),
	*ParseIPNet("203.0.113.0/24"),
	*ParseIPNet("192.88.99.0/24"),
	*ParseIPNet("192.18.0.0/15"),
	*ParseIPNet("224.0.0.0/4"),
	*ParseIPNet("240.0.0.0/4"),
	*ParseIPNet("255.255.255.255/32"),
	*ParseIPNet("100.64.0.0/10"),
	*ParseIPNet("::/128"),
	*ParseIPNet("::1/128"),
	*ParseIPNet("100::/64"),
	*ParseIPNet("2001::/23"),
	*ParseIPNet("2001:2::/48"),
	*ParseIPNet("2001:db8::/32"),
	*ParseIPNet("2001::/32"),
	*ParseIPNet("fc00::/7"),
	*ParseIPNet("fe80::/10"),
	*ParseIPNet("ff00::/8"),
	*ParseIPNet("2002::/16"),
}

func ParseIPNet(cidr string) *net.IPNet {
	if _, ipNet, err := net.ParseCIDR(cidr); err != nil {
		return nil
	} else {
		return ipNet
	}
}

func IsPublicIP(ip net.IP) bool {
	for _, ipNet := range NonPublicIPNet {
		if ipNet.Contains(ip) {
			return false
		}
	}

	return true
}

func IsCIDRContain(cidr string, ip net.IP) bool {
	if ipNet := ParseIPNet(cidr); ipNet != nil {
		return ipNet.Contains(ip)
	}

	return false
}

func IsIPNetContain(ipNet *net.IPNet, ip net.IP) bool {
	if ipNet != nil {
		return ipNet.Contains(ip)
	}

	return false
}
