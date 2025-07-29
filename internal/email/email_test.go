package email

import (
	"bytes"
	"email-forwarder/internal/config"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cfg := config.Load()
	// Create a dummy email
	rawEmail := `From: sender@example.com
To: recipient@example.com
Subject: Test Email

This is a test email.`

	email, err := Parse(strings.NewReader(rawEmail), cfg)
	if err != nil {
		t.Fatalf("Failed to parse email: %v", err)
	}

	if email.Subject != "Test Email" {
		t.Errorf("Expected subject 'Test Email', got '%s'", email.Subject)
	}

	if strings.TrimSpace(email.TextBody) != "This is a test email." {
		t.Errorf("Expected text body 'This is a test email.', got '%s'", email.TextBody)
	}
}

func TestSaveAttachment(t *testing.T) {
	// Create a dummy attachment
	attachmentContent := "This is a test attachment."
	attachment := bytes.NewBufferString(attachmentContent)
	filename := "test.txt"

	filePath, err := SaveAttachment(attachment, filename)
	if err != nil {
		t.Fatalf("Failed to save attachment: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Attachment file was not created at '%s'", filePath)
	}

	// Clean up the created file
	os.Remove(filePath)
}
