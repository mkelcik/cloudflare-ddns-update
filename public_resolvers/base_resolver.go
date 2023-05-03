package public_resolvers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

var NoIPInResponseError = errors.New("no ip found in response")

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type ipParserFunc func(reader io.Reader) (string, error)

func defaultIpParser(reader io.Reader) (string, error) {
	out, err := io.ReadAll(reader)
	return string(out), err
}

type baseResolver struct {
	client   Doer
	url      string
	ipParser ipParserFunc
}

func (i baseResolver) ResolvePublicIp(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.url, nil)
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

	ipText, err := i.ipParser(resp.Body)
	if err != nil {
		return net.IP{}, fmt.Errorf("error reading body: %w", err)
	}

	return net.ParseIP(ipText), nil
}
