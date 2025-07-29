package main

import (
	"email-forwarder/internal/config"
	"email-forwarder/internal/email"
	"email-forwarder/internal/logger"
	"email-forwarder/internal/webhook"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

type Backend struct {
	cfg        *config.Config
	dispatcher *webhook.Dispatcher
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{
		cfg:        bkd.cfg,
		dispatcher: bkd.dispatcher,
	}, nil
}

type Session struct {
	cfg        *config.Config
	dispatcher *webhook.Dispatcher
}

func (s *Session) AuthPlain(username, password string) error {
	return smtp.ErrAuthUnsupported
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	return nil
}

func (s *Session) Data(r io.Reader) error {
	parsedEmail, err := email.Parse(r, s.cfg)
	if err != nil {
		log.Printf("Error parsing email: %v", err)
		return err
	}

	s.dispatcher.AddJob(parsedEmail)
	log.Printf("Email from %v to %v queued for webhook", parsedEmail.From, parsedEmail.To)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	logger.Init()
	cfg := config.Load()

	dispatcher := webhook.NewDispatcher(cfg.WebhookConcurrency, cfg.Webhooks)
	dispatcher.Run()

	be := &Backend{
		cfg:        cfg,
		dispatcher: dispatcher,
	}

	s := smtp.NewServer(be)

	s.Addr = fmt.Sprintf(":%d", cfg.Port)
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024 * 5 // 5MB
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting SMTP server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
