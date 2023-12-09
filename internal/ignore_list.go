package internal

import (
	"net"
	"regexp"
)

func checkAddress(address, pattern string) bool {
	pattern = "^" + pattern + "$"
	re := regexp.MustCompile(pattern)
	return re.MatchString(address)
}

func IgnoredIpChange(ip net.IP, ignored []string) bool {
	if len(ignored) == 0 {
		return false
	}

	for _, i := range ignored {
		if checkAddress(ip.String(), i) {
			return true
		}
	}

	return false
}
