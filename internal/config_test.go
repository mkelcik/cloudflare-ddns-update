package internal

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		DnsRecordsToCheck   []string
		PublicIpResolverTag string
		ApiToken            string
		CloudflareZone      string
		OnChangeComment     string
		CheckInterval       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all ok",
			fields: fields{
				DnsRecordsToCheck: []string{"domain1", "domain2"},
				ApiToken:          "some_api_key",
				CloudflareZone:    "some_zone",
			},
			wantErr: false,
		},
		{
			name: "empty api token",
			fields: fields{
				DnsRecordsToCheck: []string{"domain1", "domain2"},
				ApiToken:          "",
				CloudflareZone:    "some_zone",
			},
			wantErr: true,
		},
		{
			name: "empty check dns record",
			fields: fields{
				DnsRecordsToCheck: []string{},
				ApiToken:          "",
				CloudflareZone:    "some_zone",
			},
			wantErr: true,
		},
		{
			name: "empty zone",
			fields: fields{
				DnsRecordsToCheck: []string{"domain1", "domain2"},
				ApiToken:          "some_api_key",
				CloudflareZone:    "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				DnsRecordsToCheck:   tt.fields.DnsRecordsToCheck,
				PublicIpResolverTag: tt.fields.PublicIpResolverTag,
				ApiToken:            tt.fields.ApiToken,
				CloudflareZone:      tt.fields.CloudflareZone,
				OnChangeComment:     tt.fields.OnChangeComment,
				CheckInterval:       tt.fields.CheckInterval,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
