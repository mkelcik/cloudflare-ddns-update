package public_resolvers

import (
	"net/http"
	"time"
)

const (
	V4IdentMeTag = "v4.ident.me"
	v4IdentMeUrl = "https://v4.ident.me/"
)

type V4IdentMe struct {
	baseResolver
}

func NewV4IdentMeDefault() *V4IdentMe {
	return NewV4IdentMe(&http.Client{
		Timeout: 10 * time.Second,
	})
}

func NewV4IdentMe(client Doer) *V4IdentMe {
	return &V4IdentMe{
		baseResolver: baseResolver{
			client:   client,
			url:      v4IdentMeUrl,
			ipParser: defaultIpParser,
		},
	}
}
