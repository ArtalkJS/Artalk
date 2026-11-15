package email

import (
	"sync"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/log"
)

type EmailConf struct {
	config.EmailConf
	Sender        Sender
	OnSendSuccess func(email *Email)
}

type EmailQueue struct {
	conf   EmailConf
	sender Sender
	ch     chan *Email
	done   chan struct{}
	mux    sync.Mutex
	closed bool
	pushes sync.WaitGroup
}

// Initialize Email Queue
func NewQueue(conf EmailConf) *EmailQueue {
	queue := &EmailQueue{
		conf: conf,
	}

	// init email queue
	queue.ch = make(chan *Email, conf.Queue.BufferSize)
	queue.done = make(chan struct{})

	log.Debug("[Email] Email Queue initialize complete")

	// init email sender
	if conf.Sender != nil {
		queue.sender = conf.Sender
	} else {
		if sender, err := NewSender(queue.conf); err == nil {
			queue.sender = sender
		} else {
			log.Error("[Email] Email Sender initialize failed: ", err)
		}
	}

	// init queue worker
	go func() {
		for email := range queue.ch {
			queue.handleEmail(email)
		}
	}()

	return queue
}

func (q *EmailQueue) handleEmail(email *Email) {
	// Subjects and bodies may contain login codes or notification capabilities.
	log.Debug("[Email] Sending an email")

	if q.sender == nil {
		log.Error("[Email] Email Sender is nil")
		return
	}

	if isOK := q.sender.Send(email); isOK {
		if q.conf.OnSendSuccess != nil {
			q.conf.OnSendSuccess(email)
		}
	} else {
		log.Errorf("[Email] Failed send email to addr: %s", email.ToAddr)
	}
}

// Add an email to the sending queue
func (q *EmailQueue) Push(email *Email) {
	q.mux.Lock()
	if q.closed {
		q.mux.Unlock()
		log.Error("[Email] Queue closed, dropping email")
		return
	}
	q.pushes.Add(1)
	q.mux.Unlock()
	defer q.pushes.Done()

	select {
	case q.ch <- email:
	case <-q.done:
		log.Error("[Email] Queue closed, dropping email")
	}
}

func (q *EmailQueue) Close() {
	q.mux.Lock()
	if q.closed {
		q.mux.Unlock()
		return
	}
	q.closed = true
	close(q.done)
	q.mux.Unlock()

	// Let blocked pushes stop before closing the channel; queued emails still drain.
	q.pushes.Wait()
	close(q.ch)
}
