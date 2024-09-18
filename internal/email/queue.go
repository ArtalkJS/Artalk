package email

import (
	"fmt"
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
	mux    sync.Mutex
	closed bool
}

// Initialize Email Queue
func NewQueue(conf EmailConf) *EmailQueue {
	queue := &EmailQueue{
		conf: conf,
	}

	// init email queue
	queue.ch = make(chan *Email, conf.Queue.BufferSize)

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
	log.Debug(fmt.Sprintf("[Email] Sending an email %+v: ", email))

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
	if q.closed {
		log.Error("[Email] Queue closed, dropping email")
		return
	}

	q.ch <- email
}

func (q *EmailQueue) Close() {
	q.mux.Lock()
	defer q.mux.Unlock()

	if !q.closed {
		close(q.ch)
		q.closed = true
	}
}
