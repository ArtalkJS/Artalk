package email

import (
	"sync"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	testEmail := Email{
		FromAddr: "artalkjs@gmail.com",
		FromName: "ArtalkJS",
		ToAddr:   "example@emxaple.com",
		Subject:  "TestEmail",
		Body:     "This is a test email",
	}

	t.Run("Send", func(t *testing.T) {
		sendSuccess := make(chan Email, 1)

		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: true,
			},
			OnSendSuccess: func(email *Email) {
				sendSuccess <- *email
			},
		})
		defer q.Close()

		q.Push(&testEmail)

		select {
		case email := <-sendSuccess:
			assert.Equal(t, testEmail, email)
		case <-time.After(time.Second):
			t.Fatal("Send success hook should be triggered")
		}
	})

	t.Run("Send failed", func(t *testing.T) {
		sendAttempted := make(chan Email, 1)
		sendSuccess := make(chan Email, 1)

		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: false,
				OnSend: func(email *Email) {
					sendAttempted <- *email
				},
			},
			OnSendSuccess: func(email *Email) {
				sendSuccess <- *email
			},
		})
		defer q.Close()

		q.Push(&testEmail)

		select {
		case email := <-sendAttempted:
			assert.Equal(t, testEmail, email)
		case <-time.After(time.Second):
			t.Fatal("Email send should be attempted")
		}

		select {
		case <-sendSuccess:
			t.Fatal("Send success hook should not be triggered")
		default:
		}
	})

	t.Run("Close", func(t *testing.T) {
		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: true,
			},
		})

		q.Close()

		assert.True(t, q.closed, "Queue should be closed")
	})
}

// -------------------------------------------------------------------
//  Mock Sender
// -------------------------------------------------------------------

type mockSender struct {
	SendResult bool
	OnSend     func(email *Email)
}

var _ Sender = (*mockSender)(nil)

func (s *mockSender) Send(email *Email) bool {
	if s.OnSend != nil {
		s.OnSend(email)
	}
	return s.SendResult
}

func TestQueueConcurrentCloseAndPush(t *testing.T) {
	q := NewQueue(EmailConf{Sender: &mockSender{SendResult: true}})
	start := make(chan struct{})
	var workers sync.WaitGroup
	for range 32 {
		workers.Go(func() { <-start; q.Push(&Email{ToAddr: "reader@example.org"}) })
	}
	workers.Go(func() { <-start; q.Close() })
	close(start)
	workers.Wait()
	q.Close()
	q.Push(&Email{}) // Closed queues must reject new work without a panic.
	assert.True(t, q.closed)
}
