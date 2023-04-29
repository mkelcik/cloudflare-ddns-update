package public_resolvers

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	IfConfigMeTag = "ifconfig.me"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	ifConfigMeUrl = "https://ifconfig.me"
)

type IfConfigMe struct {
	client Doer
}

func NewDefaultIfConfigMe() *IfConfigMe {
	return NewIfConfigMe(&http.Client{
		Timeout: 10 * time.Second,
	})
}

func NewIfConfigMe(c Doer) *IfConfigMe {
	return &IfConfigMe{client: c}
}

func (i IfConfigMe) ResolvePublicIp(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ifConfigMeUrl, nil)
	if err != nil {
		return net.IP{}, fmt.Errorf("error creating ifconfig request: %w", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return net.IP{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return net.IP{}, fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}

	ipText, err := io.ReadAll(resp.Body)
	if err != nil {
		return net.IP{}, fmt.Errorf("error reading body: %w", err)
	}

	return net.ParseIP(string(ipText)), nil
}
