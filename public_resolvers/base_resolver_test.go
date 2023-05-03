package public_resolvers

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_baseResolver_ResolvePublicIp(t *testing.T) {

	testUrl := "http://my-test-url.url"
	testIp := `192.168.0.100`

	client := NewTestClient(func(req *http.Request) *http.Response {

		if req.URL.String() != testUrl {
			return &http.Response{
				StatusCode: 500,
				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(`invalid url`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: io.NopCloser(bytes.NewBufferString(testIp)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	type fields struct {
		client Doer
		url    string
		fn     ipParserFunc
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			name: "check parse ip4",
			fields: fields{
				client: client,
				url:    testUrl,
				fn:     defaultIpParser,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    net.ParseIP(testIp),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := baseResolver{
				client:   tt.fields.client,
				url:      tt.fields.url,
				ipParser: tt.fields.fn,
			}
			got, err := i.ResolvePublicIp(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePublicIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolvePublicIp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
