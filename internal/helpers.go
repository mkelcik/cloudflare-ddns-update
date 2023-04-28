package internal

import "strings"

func parseDNSToCheck(data string) []string {
	out := make([]string, 0, strings.Count(data, ",")+1)
	for _, dns := range strings.Split(data, ",") {
		out = append(out, strings.TrimSpace(dns))
	}
	return out
}

func Contains[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
