package public_resolvers

import (
	"net/http"
	"time"
)

const (
	IfConfigMeTag = "ifconfig.me"
	ifConfigMeUrl = "https://ifconfig.me"
)

type IfConfigMe struct {
	baseResolver
}

func NewDefaultIfConfigMe() *IfConfigMe {
	return NewIfConfigMe(&http.Client{
		Timeout: 10 * time.Second,
	})
}

func NewIfConfigMe(client Doer) *IfConfigMe {
	return &IfConfigMe{
		baseResolver: baseResolver{
			client:   client,
			url:      ifConfigMeUrl,
			ipParser: defaultIpParser,
		},
	}
}
