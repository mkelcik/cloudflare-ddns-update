package internal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultCheckInterval = 5 * 60

	envKeyDnsToCheck           = "CLOUDFLARE_DNS_TO_CHECK"
	envKeyPublicIpResolverTag  = "PUBLIC_IP_RESOLVER"
	envKeyCloudflareApiKey     = "CLOUDFLARE_API_KEY"
	envKeyCloudflareZone       = "CLOUDFLARE_ZONE"
	envKeyOnChangeComment      = "ON_CHANGE_COMMENT"
	envKeyCheckIntervalSeconds = "CHECK_INTERVAL_SECONDS"
	envKeyNotifiers            = "NOTIFIERS"
)

type Config struct {
	DnsRecordsToCheck   []string
	PublicIpResolverTag string
	ApiToken            string
	CloudflareZone      string
	OnChangeComment     string
	Notifiers           []string
	CheckInterval       time.Duration
}

func (c Config) Validate() error {
	if c.ApiToken == "" {
		return fmt.Errorf("empty api token env key %s", envKeyCloudflareApiKey)
	}

	if c.CloudflareZone == "" {
		return fmt.Errorf("empty zone in env key %s", envKeyCloudflareZone)
	}

	if len(c.DnsRecordsToCheck) == 0 {
		return fmt.Errorf("no dns to check defined in env key %s", envKeyDnsToCheck)
	}

	return nil
}

func NewConfig() Config {
	checkInterval, err := strconv.ParseInt(os.Getenv(envKeyCheckIntervalSeconds), 10, 64)
	if err != nil {
		log.Printf("wrong `%s` value. Check interval set default(%ds)", envKeyCheckIntervalSeconds, defaultCheckInterval)
		checkInterval = defaultCheckInterval
	}

	return Config{
		DnsRecordsToCheck:   parseCommaDelimited(os.Getenv(envKeyDnsToCheck)),
		PublicIpResolverTag: os.Getenv(envKeyPublicIpResolverTag),
		ApiToken:            os.Getenv(envKeyCloudflareApiKey),
		CloudflareZone:      os.Getenv(envKeyCloudflareZone),
		OnChangeComment:     os.Getenv(envKeyOnChangeComment),
		Notifiers:           parseCommaDelimited(os.Getenv(envKeyNotifiers)),
		CheckInterval:       time.Duration(checkInterval) * time.Second,
	}
}
