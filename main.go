package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/mkelcik/cloudflare-ddns-update/internal"
	"github.com/mkelcik/cloudflare-ddns-update/notifications"
	"github.com/mkelcik/cloudflare-ddns-update/public_resolvers"
)

type PublicIpResolver interface {
	ResolvePublicIp(ctx context.Context) (net.IP, error)
}

func getNotifiers(tags []string) notifications.Notifiers {
	out := notifications.Notifiers{}
	for _, t := range tags {
		if initFn, ok := notifications.Available[t]; ok {
			notifier, err := initFn()
			if err != nil {
				log.Println(err)
				continue
			}
			out = append(out, notifier)
		}
	}
	return out
}

func getResolver(resolverName string) (PublicIpResolver, string) {
	switch resolverName {
	// HERE add another resolver if needed
	case public_resolvers.CloudflareTraceTag:
		return public_resolvers.NewDefaultCloudflareTrace(), public_resolvers.CloudflareTraceTag
	case public_resolvers.V4IdentMeTag:
		return public_resolvers.NewV4IdentMeDefault(), public_resolvers.V4IdentMeTag
	case public_resolvers.IfConfigMeTag:
		fallthrough
	default:
		return public_resolvers.NewDefaultIfConfigMe(), public_resolvers.IfConfigMeTag
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := internal.NewConfig()
	if err := config.Validate(); err != nil {
		log.Fatalln(err)
	}

	api, err := cloudflare.NewWithAPIToken(config.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch user details on the account
	zoneID, err := api.ZoneIDByName(config.CloudflareZone)
	if err != nil {
		log.Fatal(err)
	}

	notifiers := getNotifiers(config.Notifiers)

	// public ip resolver
	publicIpResolver, resolverTag := getResolver(config.PublicIpResolverTag)

	checkFunc := func() {
		currentPublicIP, err := publicIpResolver.ResolvePublicIp(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Current public ip `%s` (resolver: %s)", currentPublicIP, resolverTag)

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
					continue
				}

				if err := notifiers.NotifyWithLog(ctx, notifications.Notification{
					OldIp:       net.ParseIP(dnsRecord.Content),
					NewIp:       currentPublicIP,
					CheckedAt:   time.Now(),
					ResolverTag: resolverTag,
					Domain:      dnsRecord.Name,
				}); err != nil {
					log.Printf("errors in notifications: %s", err)
				}
				log.Printf("Updated to `%s`", currentPublicIP)
			}
		}
	}

	log.Printf("checking ...")
	checkFunc()

	log.Println("waiting for check tick ...")
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("tick received checking ...")
			checkFunc()
		case <-ctx.Done():
			break
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
