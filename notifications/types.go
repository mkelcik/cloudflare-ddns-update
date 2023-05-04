package notifications

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	configDelimiter = "@"
)

type Notifiers []Notifier

func (n Notifiers) NotifyWithLog(ctx context.Context, notification Notification) error {
	var outErr error
	for _, notifier := range n {
		if err := notifier.Notify(ctx, notification); err != nil {
			outErr = errors.Join(outErr, err)
			continue
		}
		log.Printf("Notification sent via %s\n", notifier.Tag())
	}
	return outErr
}

type Notification struct {
	OldIp       net.IP    `json:"old_ip,omitempty"`
	NewIp       net.IP    `json:"new_ip"`
	CheckedAt   time.Time `json:"checked_at"`
	ResolverTag string    `json:"resolver_tag"`
	Domain      string    `json:"domain"`
}

func (n Notification) ToSlice() []string {
	return []string{n.OldIp.String(), n.NewIp.String(), n.CheckedAt.Format(time.RFC3339), n.ResolverTag, n.Domain}
}

var Available = map[string]func(string) (Notifier, error){
	webhookTag: func(config string) (Notifier, error) {
		parts := strings.Split(config, configDelimiter)

		if len(parts) < 2 {
			return nil, fmt.Errorf("wrong webhook config, missing url part")
		}

		return NewWebhookNotification(WebhookConfig{Url: parts[1]}, &http.Client{
			Timeout: 10 * time.Second,
		}), nil
	},
}

type Notifier interface {
	Tag() string
	Notify(ctx context.Context, notification Notification) error
}

func GetNotifiers(tags []string) Notifiers {
	out := Notifiers{}
	for _, t := range tags {
		if initFn, ok := Available[strings.Split(t, configDelimiter)[0]]; ok {
			notifier, err := initFn(t)
			if err != nil {
				log.Println(err)
				continue
			}
			out = append(out, notifier)
		}
	}
	return out
}
