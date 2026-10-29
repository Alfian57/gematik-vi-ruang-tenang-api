package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WhatsAppSender interface {
	Send(ctx context.Context, number, message string) error
	IsConfigured() bool
}

type fonnteClient struct {
	token    string
	endpoint string
	client   *http.Client
}

func NewFonnteClient(token string) WhatsAppSender {
	// WhatsApp OTP delivery must always connect directly to api.fonnte.com and never
	// route through container/system HTTP_PROXY or HTTPS_PROXY (e.g. Cloudflare WARP).
	directTransport := &http.Transport{
		Proxy: nil, // bypass any HTTP_PROXY / HTTPS_PROXY
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &fonnteClient{
		token:    strings.TrimSpace(token),
		endpoint: "https://api.fonnte.com/send",
		client:   &http.Client{Timeout: 15 * time.Second, Transport: directTransport},
	}
}

func (c *fonnteClient) IsConfigured() bool { return c.token != "" }

func (c *fonnteClient) Send(ctx context.Context, number, message string) error {
	if !c.IsConfigured() {
		return errors.New("WhatsApp sender is not configured (FONNTE_TOKEN is empty)")
	}
	target := strings.TrimPrefix(strings.TrimSpace(number), "+")
	form := url.Values{"target": {target}, "message": {message}, "countryCode": {"62"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Fonnte request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return fmt.Errorf("failed to read Fonnte response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("WhatsApp provider HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Status bool   `json:"status"`
		Reason string `json:"reason"`
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("failed to parse Fonnte response: %w", err)
	}
	if !result.Status {
		msg := result.Reason
		if msg == "" {
			msg = result.Detail
		}
		if msg == "" {
			msg = string(bodyBytes)
		}
		return fmt.Errorf("WhatsApp provider rejected: %s", msg)
	}
	return nil
}
