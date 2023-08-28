package email

import (
	"fmt"
	"sync"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/log"
)

type EmailConf struct {
	config.EmailConf
	OnSendSuccess func(email *Email)
}

type EmailQueue struct {
	conf   *EmailConf
	ch     chan *Email
	mux    sync.Mutex
	closed bool
}

// Initialize Email Queue
func NewQueue(conf *EmailConf) *EmailQueue {
	queue := &EmailQueue{
		conf: conf,
	}

	// init email queue
	queue.ch = make(chan *Email, conf.Queue.BufferSize)

	log.Debug("[Email] Email Queue initialize complete")

	// init queue worker
	go queue.worker()

	return queue
}

func (q *EmailQueue) worker() {
	for email := range q.ch {
		q.handleEmail(email)
	}
}

func (q *EmailQueue) handleEmail(email *Email) {
	sender := NewSender(q.conf)

	log.Debug(fmt.Sprintf("[Email] Sending an email %+v: ", email))

	// send email
	if isOK := sender.Send(email); isOK {
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
