package webhook

import (
	"bytes"
	"email-forwarder/internal/email"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Job struct {
	Email *email.ParsedEmail
}

type Worker struct {
	id         int
	jobQueue   chan Job
	webhookURL string
	client     *http.Client
}

func NewWorker(id int, jobQueue chan Job, webhookURL string) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   jobQueue,
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobQueue {
			log.Printf("Worker %d: processing job", w.id)
			err := w.send(job.Email)
			if err != nil {
				log.Printf("Worker %d: error sending webhook: %v", w.id, err)
			}
		}
	}()
}

func (w *Worker) send(email *email.ParsedEmail) error {
	jsonData, err := json.Marshal(email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.webhookURL, bytes.NewBuffer(jsonData))
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
	webhookURL string
}

func NewDispatcher(maxWorkers int, webhookURL string) *Dispatcher {
	return &Dispatcher{
		jobQueue:   make(chan Job, 100),
		maxWorkers: maxWorkers,
		webhookURL: webhookURL,
	}
}

func (d *Dispatcher) Run() {
	for i := 1; i <= d.maxWorkers; i++ {
		worker := NewWorker(i, d.jobQueue, d.webhookURL)
		worker.Start()
	}
}

func (d *Dispatcher) AddJob(email *email.ParsedEmail) {
	d.jobQueue <- Job{Email: email}
}
