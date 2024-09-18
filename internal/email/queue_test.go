package email

import (
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
		sendSuccessHookTriggered := false

		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: true,
			},
			OnSendSuccess: func(email *Email) {
				assert.Equal(t, testEmail, *email)
				sendSuccessHookTriggered = true
			},
		})

		q.Push(&testEmail)

		time.Sleep(100 * time.Millisecond)

		assert.True(t, sendSuccessHookTriggered, "Send success hook should be triggered")
	})

	t.Run("Send failed", func(t *testing.T) {
		sendSuccessHookTriggered := false

		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: false,
			},
			OnSendSuccess: func(email *Email) {
				assert.Equal(t, testEmail, email)
				sendSuccessHookTriggered = true
			},
		})

		q.Push(&testEmail)

		time.Sleep(100 * time.Millisecond)

		assert.False(t, sendSuccessHookTriggered, "Send success hook should not be triggered")
	})

	t.Run("Close", func(t *testing.T) {
		q := NewQueue(EmailConf{
			EmailConf: config.EmailConf{},
			Sender: &mockSender{
				SendResult: true,
			},
		})

		q.Close()

		time.Sleep(100 * time.Millisecond)

		assert.True(t, q.closed, "Queue should be closed")
	})
}

// -------------------------------------------------------------------
//  Mock Sender
// -------------------------------------------------------------------

type mockSender struct {
	SendResult bool
}

var _ Sender = (*mockSender)(nil)

func (s *mockSender) Send(email *Email) bool {
	return s.SendResult
}
