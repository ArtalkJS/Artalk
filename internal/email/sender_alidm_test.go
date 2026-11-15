package email

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

type mailTransport func(*http.Request) (*http.Response, error)

func (f mailTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestAliDMDeliveryResult(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	for _, tc := range []struct {
		name, response string
		failure        error
		success        bool
	}{
		{"accepted", `{"EnvId":"queued","RequestId":"fixture"}`, nil, true},
		{"rejected", `{"Code":"InvalidToAddress","Message":"invalid recipient","RequestId":"fixture"}`, nil, false},
		{"invalid response", `not JSON`, nil, false},
		{"network failure", "", errors.New("provider unavailable"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultTransport = mailTransport(func(r *http.Request) (*http.Response, error) {
				require.NoError(t, r.ParseForm())
				require.Equal(t, "SingleSendMail", r.Form.Get("Action"))
				require.Equal(t, "reader@example.org", r.Form.Get("ToAddress"))
				require.Equal(t, "A reply", r.Form.Get("Subject"))
				require.Equal(t, "<p>Reply</p>", r.Form.Get("HtmlBody"))
				if tc.failure != nil {
					return nil, tc.failure
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(tc.response))}, nil
			})
			sender, err := NewSender(EmailConf{EmailConf: config.EmailConf{SendType: config.TypeAliDM}})
			require.NoError(t, err)
			require.Equal(t, tc.success, sender.Send(&Email{ToAddr: "reader@example.org", Subject: "A reply", Body: "<p>Reply</p>"}))
		})
	}
}
