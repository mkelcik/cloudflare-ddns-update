package notifications

import (
	"io"
	"reflect"
	"testing"
)

func TestWebhookNotification_getRequestBody(t *testing.T) {
	type fields struct {
		config WebhookConfig
	}
	type args struct {
		notification Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.Reader
		wantErr bool
	}{
		{
			name:    "text",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WebhookNotification{
				config: tt.fields.config,
			}
			got, err := w.getRequestBody(tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRequestBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRequestBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
