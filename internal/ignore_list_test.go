package internal

import "testing"

func Test_checkAddress(t *testing.T) {
	type args struct {
		address string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				address: "192.168.0.1",
				pattern: "",
			},
			want: false,
		}, {
			name: "true match 192.*",
			args: args{
				address: "192.168.0.1",
				pattern: "192.*",
			},
			want: true,
		}, {
			name: "false match 193.*",
			args: args{
				address: "192.168.0.1",
				pattern: "193.*",
			},
			want: false,
		},
		{
			name: "true match 192.168.0.1",
			args: args{
				address: "192.168.0.1",
				pattern: "192.168.0.1",
			},
			want: true,
		},
		{
			name: "false not match 192.168.0.2",
			args: args{
				address: "192.168.0.1",
				pattern: "192.168.0.2",
			},
			want: false,
		}, {
			name: "true match 192.168.0.*",
			args: args{
				address: "192.168.0.10",
				pattern: "192.168.0.*",
			},
			want: true,
		}, {
			name: "false match 192.168.0.*",
			args: args{
				address: "192.168.1.10",
				pattern: "192.168.0.*",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAddress(tt.args.address, tt.args.pattern); got != tt.want {
				t.Errorf("checkAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
