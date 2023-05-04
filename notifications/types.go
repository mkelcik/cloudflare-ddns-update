package notifications

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type Notifiers []Notifier

func (n Notifiers) NotifyWithLog(ctx context.Context, notification Notification) error {
	var outErr error
	for _, notifier := range n {
		if err := notifier.Notify(ctx, notification); err != nil {
			outErr = errors.Join(outErr, err)
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

var Available = map[string]func() (Notifier, error){
	webhookTag: func() (Notifier, error) {
		return NewWebhookNotification(NewWebhookConfigFromEnv(), &http.Client{
			Timeout: 10 * time.Second,
		}), nil
	},
}

type Notifier interface {
	Tag() string
	Notify(ctx context.Context, notification Notification) error
}
