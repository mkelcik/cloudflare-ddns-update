package public_resolvers

import (
	"bytes"
	"io"
	"testing"
)

func Test_cloudflareTraceResponseParser(t *testing.T) {

	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				reader: bytes.NewBuffer([]byte(`fl=31f118
h=1.1.1.1
ip=94.113.142.206
ts=1683145336.383
visit_scheme=https
uag=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36
colo=PRG
sliver=none
http=http/2
loc=CZ
tls=TLSv1.3
sni=off
warp=off
gateway=off
rbi=off
kex=X25519`)),
			},
			want:    "94.113.142.206",
			wantErr: false,
		},
		{
			name: "no ip in response",
			args: args{
				reader: bytes.NewBuffer([]byte(`fl=31f118
h=1.1.1.1
ts=1683145336.383
visit_scheme=https
uag=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36
colo=PRG
sliver=none
http=http/2
loc=CZ
tls=TLSv1.3
sni=off
warp=off
gateway=off
rbi=off
kex=X25519`)),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cloudflareTraceResponseParser(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("cloudflareTraceResponseParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cloudflareTraceResponseParser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
