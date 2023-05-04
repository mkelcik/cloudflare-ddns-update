package internal

import "strings"

func parseCommaDelimited(data string) []string {
	out := make([]string, 0, strings.Count(data, ",")+1)
	for _, item := range strings.Split(data, ",") {
		if w := strings.TrimSpace(item); w != "" {
			out = append(out, w)
		}
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
