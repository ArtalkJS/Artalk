package aliyun

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestGreenTextModeration(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	for _, tc := range []struct {
		name, response string
		allowed        bool
		wantError      bool
	}{
		{"allow", `{"code":200,"data":[{"results":[{"suggestion":"pass"}]}]}`, true, false},
		{"review", `{"code":200,"data":[{"results":[{"suggestion":"review"}]}]}`, false, false},
		{"missing result", `{"code":200,"data":[]}`, false, true},
		{"upstream rejection", `{"code":400,"msg":"invalid request"}`, false, true},
		{"unexpected body", "not JSON", false, true},
		{"offline", "", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultTransport = transportFunc(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "https://green.cn-shanghai.aliyuncs.com/green/text/scan", r.URL.String())
				require.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "acs fixture-id:"))
				require.NotEmpty(t, r.Header.Get("Content-MD5"))
				require.Equal(t, "2018-05-09", r.Header.Get("X-Acs-Version"))
				var data map[string]any
				require.NoError(t, json.NewDecoder(r.Body).Decode(&data))
				require.Equal(t, []any{"antispam"}, data["scenes"])
				require.Equal(t, map[string]any{"dataId": "17", "content": "A comment"}, data["tasks"])
				if tc.response == "" {
					return nil, errors.New("provider offline")
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(tc.response))}, nil
			})

			ok, err := GreenText(GreenTextConf{AccessKeyId: "fixture-id", AccessKeySecret: "fixture-key", Content: "A comment", DataID: "17"})
			require.Equal(t, tc.allowed, ok)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
