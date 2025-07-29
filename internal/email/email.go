package email

import (
	"bytes"
	"email-forwarder/internal/config"
	"io"
	"net/mail"

	"github.com/mnako/letters"
)

type ParsedEmail struct {
	From        []*mail.Address `json:"from"`
	To          []*mail.Address `json:"to"`
	Subject     string          `json:"subject"`
	HTMLBody    string          `json:"html_body"`
	TextBody    string          `json:"text_body"`
	Attachments []*Attachment   `json:"attachments"`
}

func Parse(r io.Reader, cfg *config.Config) (*ParsedEmail, error) {
	email, err := letters.ParseEmail(r)
	if err != nil {
		return nil, err
	}

	attachments := make([]*Attachment, 0)
	for _, a := range email.AttachedFiles {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, bytes.NewReader(a.Data)); err != nil {
			return nil, err
		}

		path, err := SaveAttachment(buf, a.ContentType.Params["name"])
		if err != nil {
			// Log the error but continue processing
			continue
		}

		attachments = append(attachments, &Attachment{
			Filename:    a.ContentType.Params["name"],
			ContentType: a.ContentType.ContentType,
			Size:        len(a.Data),
			S3Url:       path,
		})
	}

	return &ParsedEmail{
		From:        email.Headers.From,
		To:          email.Headers.To,
		Subject:     email.Headers.Subject,
		HTMLBody:    email.HTML,
		TextBody:    email.Text,
		Attachments: attachments,
	}, nil
}
