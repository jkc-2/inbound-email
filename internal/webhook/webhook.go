package webhook

import (
	"email-forwarder/internal/config"
	"email-forwarder/internal/email"
	"log"
	"net/http"
	"strings"
	"time"
)

type Job struct {
	Email *email.ParsedEmail
}

type Worker struct {
	id         int
	jobQueue   chan Job
	webhooks   []config.WebhookConfig
	client     *http.Client
}

func NewWorker(id int, jobQueue chan Job, webhooks []config.WebhookConfig) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   jobQueue,
		webhooks:   webhooks,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobQueue {
			log.Printf("Worker %d: processing job", w.id)
			w.processJob(job)
		}
	}()
}

func (w *Worker) processJob(job Job) {
	for _, webhook := range w.webhooks {
		if w.shouldSend(job.Email, webhook) {
			err := w.send(job.Email, webhook)
			if err != nil {
				log.Printf("Worker %d: error sending webhook '%s': %v", w.id, webhook.Name, err)
			}
		}
	}
}

func (w *Worker) shouldSend(email *email.ParsedEmail, webhook config.WebhookConfig) bool {
	senderMatch := false
	for _, s := range webhook.Senders {
		if s == "*" {
			senderMatch = true
			break
		}
		for _, from := range email.From {
			if strings.Contains(from.Address, s) {
				senderMatch = true
				break
			}
		}
		if senderMatch {
			break
		}
	}

	recipientMatch := false
	for _, r := range webhook.Recipients {
		if r == "*" {
			recipientMatch = true
			break
		}
		for _, to := range email.To {
			if strings.Contains(to.Address, r) {
				recipientMatch = true
				break
			}
		}
		if recipientMatch {
			break
		}
	}

	return senderMatch && recipientMatch
}

func (w *Worker) send(email *email.ParsedEmail, webhook config.WebhookConfig) error {
	payload, err := FormatPayload(email, webhook.Template)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhook.URL, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type Dispatcher struct {
	jobQueue   chan Job
	maxWorkers int
	webhooks   []config.WebhookConfig
}

func NewDispatcher(maxWorkers int, webhooks []config.WebhookConfig) *Dispatcher {
	return &Dispatcher{
		jobQueue:   make(chan Job, 100),
		maxWorkers: maxWorkers,
		webhooks:   webhooks,
	}
}

func (d *Dispatcher) Run() {
	for i := 1; i <= d.maxWorkers; i++ {
		worker := NewWorker(i, d.jobQueue, d.webhooks)
		worker.Start()
	}
}

func (d *Dispatcher) AddJob(email *email.ParsedEmail) {
	d.jobQueue <- Job{Email: email}
}
