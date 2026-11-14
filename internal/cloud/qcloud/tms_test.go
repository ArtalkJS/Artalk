package qcloud

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTextModeration(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	for _, tc := range []struct {
		name, response string
		allowed        bool
		wantError      bool
	}{
		{"allow", `{"Response":{"Suggestion":"Pass"}}`, true, false},
		{"block", `{"Response":{"Suggestion":"Block"}}`, false, false},
		{"missing suggestion", `{"Response":{}}`, false, true},
		{"upstream rejection", `{"Response":{"RequestId":"r1","Error":{"Message":"invalid request"}}}`, false, true},
		{"unexpected body", "not JSON", false, true},
		{"offline", "", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultTransport = transportFunc(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "tms.tencentcloudapi.com", r.URL.Host)
				require.Equal(t, "TextModeration", r.Header.Get("X-TC-Action"))
				require.Equal(t, "ap-guangzhou", r.Header.Get("X-TC-Region"))
				require.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "TC3-HMAC-SHA256 Credential=fixture-id/"))
				require.Equal(t, "QSBjb21tZW50", r.URL.Query().Get("Content"))
				require.Equal(t, "17", r.URL.Query().Get("DataId"))
				if tc.response == "" {
					return nil, errors.New("provider offline")
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(tc.response))}, nil
			})

			ok, err := TMS(TmsConf{SecretID: "fixture-id", SecretKey: "fixture-key", Content: "A comment", DataID: "17"})
			require.Equal(t, tc.allowed, ok)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
