package captcha

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

type verifierTransport func(*http.Request) (*http.Response, error)

func (f verifierTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRemoteCaptchaVerification(t *testing.T) {
	user := &User{IP: "192.0.2.42"}
	const userToken = "issued-test-token"
	const secret = "fixture-secret"
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })

	for _, tc := range []struct {
		name, endpoint, input, accepted, rejected string
		checker                                   Checker
		fields                                    url.Values
	}{
		{
			name: "hCaptcha", endpoint: HCAPTCHA_API, input: userToken,
			accepted: `{"success":true}`, rejected: `{"success":false,"error-codes":["invalid-input-response"]}`,
			checker: NewCaptchaChecker(&CheckerConf{CaptchaConf: config.CaptchaConf{CaptchaType: config.TypeHCaptcha, HCaptcha: config.HCaptchaConf{SecretKey: secret}}, User: *user}),
			fields:  url.Values{"secret": {secret}, "response": {userToken}, "remoteip": {user.IP}},
		},
		{
			name: "reCAPTCHA", endpoint: RECAPTCHA_API, input: userToken,
			accepted: `{"success":true}`, rejected: `{"success":false,"error-codes":["invalid-input-response"]}`,
			checker: NewCaptchaChecker(&CheckerConf{CaptchaConf: config.CaptchaConf{CaptchaType: config.TypeReCaptcha, ReCaptcha: config.ReCaptchaConf{SecretKey: secret}}, User: *user}),
			fields:  url.Values{"secret": {secret}, "response": {userToken}, "remoteip": {user.IP}},
		},
		{
			name: "Turnstile", endpoint: TURNSTILE_API, input: userToken,
			accepted: `{"success":true}`, rejected: `{"success":false,"error-codes":["invalid-input-response"]}`,
			checker: NewCaptchaChecker(&CheckerConf{CaptchaConf: config.CaptchaConf{CaptchaType: config.TypeTurnstile, Turnstile: config.TurnstileConf{SecretKey: secret}}, User: *user}),
			fields:  url.Values{"secret": {secret}, "response": {userToken}, "remoteip": {user.IP}},
		},
		{
			name: "Geetest", endpoint: GEETEST_API + "/validate?captcha_id=test-site",
			input:    `{"lot_number":"test-lot","captcha_output":"fixture","pass_token":"issued-test-token","gen_time":"123"}`,
			accepted: `{"result":"success"}`, rejected: `{"result":"fail","reason":"invalid result"}`,
			checker: NewCaptchaChecker(&CheckerConf{CaptchaConf: config.CaptchaConf{CaptchaType: config.TypeGeetest, Geetest: config.GeetestConf{CaptchaID: "test-site", CaptchaKey: secret}}, User: *user}),
			fields: url.Values{
				"lot_number": {"test-lot"}, "captcha_output": {"fixture"}, "pass_token": {userToken},
				"gen_time": {"123"}, "sign_token": {"c672b146a5a8d4a98569121df6a3c01dbca9ca6df46ba10ceaf74d3ac25b83c7"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, IFrame, tc.checker.Type())
			for _, result := range []struct {
				body string
				pass bool
			}{{tc.accepted, true}, {tc.rejected, false}} {
				requests := 0
				http.DefaultTransport = verifierTransport(func(r *http.Request) (*http.Response, error) {
					requests++
					require.Equal(t, http.MethodPost, r.Method)
					require.Equal(t, tc.endpoint, r.URL.String())
					require.NoError(t, r.ParseForm())
					require.Equal(t, tc.fields, r.PostForm)
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(result.body))}, nil
				})
				ok, err := tc.checker.Check(tc.input)
				require.Equal(t, result.pass, ok)
				if result.pass {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
				require.Equal(t, 1, requests)
			}

			http.DefaultTransport = verifierTransport(func(*http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("verification service offline")
			})
			ok, err := tc.checker.Check(tc.input)
			require.False(t, ok)
			require.ErrorContains(t, err, "verification service offline")
		})
	}
}
