package webhook

import (
	"email-forwarder/internal/config"
	"email-forwarder/internal/email"
	"net/mail"
	"testing"
)

func TestShouldSend(t *testing.T) {
	webhooks := []config.WebhookConfig{
		{
			Name:      "test-webhook",
			URL:       "http://localhost:8080",
			Senders:   []string{"test@example.com"},
			Recipients: []string{"recipient@example.com"},
			Template:  "{{.Subject}}",
		},
		{
			Name:      "catch-all",
			URL:       "http://localhost:8080",
			Senders:   []string{"*"},
			Recipients: []string{"*"},
			Template:  "{{.Subject}}",
		},
	}

	parsedEmail := &email.ParsedEmail{
		From: []*mail.Address{
			{Address: "test@example.com"},
		},
		To: []*mail.Address{
			{Address: "recipient@example.com"},
		},
		Subject: "Test Email",
	}

	worker := NewWorker(1, make(chan Job), webhooks)

	if !worker.shouldSend(parsedEmail, webhooks[0]) {
		t.Error("Expected shouldSend to return true for test-webhook")
	}

	if !worker.shouldSend(parsedEmail, webhooks[1]) {
		t.Error("Expected shouldSend to return true for catch-all")
	}

	parsedEmail.From[0].Address = "other@example.com"

	if worker.shouldSend(parsedEmail, webhooks[0]) {
		t.Error("Expected shouldSend to return false for test-webhook")
	}
}

func TestFormatPayload(t *testing.T) {
	parsedEmail := &email.ParsedEmail{
		Subject: "Test Email",
	}

	template := "{{.Subject}}"

	payload, err := FormatPayload(parsedEmail, template)
	if err != nil {
		t.Fatalf("FormatPayload failed: %v", err)
	}

	if payload.String() != "Test Email" {
		t.Errorf("Expected payload 'Test Email', got '%s'", payload.String())
	}
}
