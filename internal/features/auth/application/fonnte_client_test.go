package application

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (fn transportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return fn(req) }

func TestNormalizeWhatsAppNumber(t *testing.T) {
	for _, input := range []string{"081234567890", "6281234567890", "+6281234567890"} {
		got, err := normalizeWhatsAppNumber(input)
		if err != nil || got != "6281234567890" {
			t.Fatalf("normalize %q: %q, %v", input, got, err)
		}
	}
	for _, input := range []string{"", "0212345678", "+62812abc7890", "6281"} {
		if _, err := normalizeWhatsAppNumber(input); err == nil {
			t.Fatalf("accepted invalid number %q", input)
		}
	}
}

func TestFonnteSendChecksProviderResponse(t *testing.T) {
	accepted := true
	client := NewFonnteClient("test-token").(*fonnteClient)
	client.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "test-token" || r.FormValue("target") != "6281234567890" || r.FormValue("message") != "kode uji" {
			t.Errorf("unexpected Fonnte request")
		}
		body, _ := json.Marshal(map[string]bool{"status": accepted})
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}, nil
	})
	if err := client.Send(context.Background(), "6281234567890", "kode uji"); err != nil {
		t.Fatal(err)
	}
	accepted = false
	if err := client.Send(context.Background(), "6281234567890", "kode uji"); err == nil {
		t.Fatal("expected provider rejection")
	}
}

func TestFonnteClientUsesDirectTransport(t *testing.T) {
	client := NewFonnteClient("test-token").(*fonnteClient)
	transport, ok := client.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.client.Transport)
	}
	if transport.Proxy != nil {
		t.Fatal("expected transport.Proxy to be nil (direct connection)")
	}
}
