package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	webhookTag             = "webhook"
	webhookRequestTypeJson = "JSON"
	envWebhookUrl          = "WEBHOOK_RL"
	envWebhookRequestType  = "WEBHOOK_REQ_TYPE"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type WebhookConfig struct {
	Url  string
	Json bool
}

func NewWebhookConfigFromEnv() WebhookConfig {
	return WebhookConfig{
		Url:  os.Getenv(envWebhookUrl),
		Json: strings.ToUpper(os.Getenv(envWebhookRequestType)) == webhookRequestTypeJson,
	}
}

type WebhookNotification struct {
	config WebhookConfig
	client Doer
}

func (w WebhookNotification) Tag() string {
	return webhookTag
}

func NewWebhookNotification(config WebhookConfig, client Doer) *WebhookNotification {
	return &WebhookNotification{config: config, client: client}
}

func (w WebhookNotification) getRequestBody(notification Notification) (io.Reader, error) {
	out := bytes.NewBuffer(notification.NewIp)
	if w.config.Json {
		if err := json.NewEncoder(out).Encode(notification); err != nil {
			return nil, fmt.Errorf("error encoding notification body: %w", err)
		}
		return out, nil
	}
	return out, nil
}

func (w WebhookNotification) Notify(ctx context.Context, notification Notification) error {
	body, err := w.getRequestBody(notification)
	if err != nil {
		return fmt.Errorf("WebhookNotification::NotifyWithLog error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.config.Url, body)
	if err != nil {
		return fmt.Errorf("WebhookNotification::NotifyWithLog error creating request: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("WebhookNotification::NotifyWithLog error while sending notification: %w", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("WebhookNotification::NotifyWithLog unexpected non 2xx code %d returned", resp.StatusCode)
	}
	return nil
}
