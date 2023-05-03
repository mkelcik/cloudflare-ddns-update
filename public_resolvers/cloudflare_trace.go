package public_resolvers

import (
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	CloudflareTraceTag = "1.1.1.1"
	CloudflareTraceUrl = "https://1.1.1.1/cdn-cgi/trace"

	ipPrefix = "ip="
)

type CloudflareTrace struct {
	baseResolver
}

func NewDefaultCloudflareTrace() *CloudflareTrace {
	return NewCloudflareTrace(&http.Client{
		Timeout: 10 * time.Second,
	})
}

func cloudflareTraceResponseParser(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	for _, row := range strings.Split(string(data), "\n") {
		if strings.Index(row, ipPrefix) == 0 {
			return strings.TrimSpace(strings.ReplaceAll(row, ipPrefix, "")), nil
		}
	}
	return "", NoIPInResponseError
}

func NewCloudflareTrace(client Doer) *CloudflareTrace {
	return &CloudflareTrace{
		baseResolver: baseResolver{
			client:   client,
			url:      CloudflareTraceUrl,
			ipParser: cloudflareTraceResponseParser,
		},
	}
}
