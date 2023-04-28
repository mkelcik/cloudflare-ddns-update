package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/mkelcik/cloudflare-ddns-update/internal"
	"github.com/mkelcik/cloudflare-ddns-update/public_resolvers"
)

type PublicIpResolver interface {
	ResolvePublicIp(ctx context.Context) (net.IP, error)
}

func getResolver(resolverName string) PublicIpResolver {
	switch resolverName {
	// HERE add another resolver if needed
	case public_resolvers.IfConfigMeTag:
		fallthrough
	default:
		return public_resolvers.NewIfConfigMe(&http.Client{
			Timeout: 10 * time.Second,
		})
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := internal.NewConfig()
	if err := config.Validate(); err != nil {
		log.Fatalln(err)
	}

	currentPublicIP, err := getResolver(config.PublicIpResolverTag).ResolvePublicIp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current public ip `%s`", currentPublicIP)

	api, err := cloudflare.NewWithAPIToken(config.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch user details on the account
	zoneID, err := api.ZoneIDByName(config.CloudflareZone)
	if err != nil {
		log.Fatal(err)
	}

	dns, err := allDNSRecords(ctx, api, cloudflare.ZoneIdentifier(zoneID))
	if err != nil {
		log.Fatal(err)
	}

	for _, dnsRecord := range dns {
		if internal.Contains(config.DnsRecordsToCheck, dnsRecord.Name) {
			log.Printf("Checking record `%s` with current value `%s` ...", dnsRecord.Name, dnsRecord.Content)
			if currentPublicIP.String() == dnsRecord.Content {
				log.Println("OK")
				continue // no update needed
			}

			update := cloudflare.UpdateDNSRecordParams{
				ID:      dnsRecord.ID,
				Content: currentPublicIP.String(),
			}

			if config.OnChangeComment != "" {
				update.Comment = config.OnChangeComment
			}

			if _, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), update); err != nil {
				log.Printf("error updating dns record: %s", err)
			} else {
				log.Printf("Updated to `%s`", currentPublicIP)
			}
		}
	}
}

func allDNSRecords(ctx context.Context, api *cloudflare.API, rc *cloudflare.ResourceContainer) ([]cloudflare.DNSRecord, error) {
	out := make([]cloudflare.DNSRecord, 0, 100)
	params := cloudflare.ListDNSRecordsParams{
		ResultInfo: cloudflare.ResultInfo{Page: 1},
	}
	for {
		page, res, err := api.ListDNSRecords(ctx, rc, params)
		if err != nil {
			return nil, err
		}
		out = append(out, page...)

		if res.Page >= res.TotalPages {
			break
		}
		params.Page++
	}
	return out, nil
}
