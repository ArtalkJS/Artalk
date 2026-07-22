package email

import (
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailMessageHeaders(t *testing.T) {
	email := &Email{
		FromAddr: "noreply@example.com",
		FromName: "Artalk",
		ToAddr:   "recipient@example.net",
		Subject:  "Test email",
		Body:     "Hello from Artalk",
	}

	first := readEmailMessage(t, getEmailMineTxt(email))
	second := readEmailMessage(t, getEmailMineTxt(email))

	firstMessageID := first.Header.Get("Message-ID")
	secondMessageID := second.Header.Get("Message-ID")

	assert.Regexp(t, `^<[A-Za-z0-9]{32}@example\.com>$`, firstMessageID)
	assert.Regexp(t, `^<[A-Za-z0-9]{32}@example\.com>$`, secondMessageID)
	assert.NotEqual(t, firstMessageID, secondMessageID)
	assert.NotEmpty(t, first.Header.Get("Date"))
	assert.Equal(t, "1.0", first.Header.Get("MIME-Version"))
}

func TestNewMessageIDFallsBackForInvalidFromAddress(t *testing.T) {
	assert.Regexp(t, `^<[A-Za-z0-9]{32}@artalk\.local>$`, newMessageID("invalid address"))
}

func readEmailMessage(t *testing.T, raw string) *mail.Message {
	t.Helper()

	message, err := mail.ReadMessage(strings.NewReader(raw))
	require.NoError(t, err)
	return message
}
