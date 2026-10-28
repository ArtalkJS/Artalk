package server_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
)

func TestServeListenAddresses(t *testing.T) {
	// Use a trusted local certificate so HTTPS tests also verify the handshake.
	tlsServer := httptest.NewTLSServer(nil)
	cert := tlsServer.Certificate()
	key, err := x509.MarshalPKCS8PrivateKey(tlsServer.TLS.Certificates[0].PrivateKey)
	tlsServer.Close()
	require.NoError(t, err)
	certPath, keyPath := filepath.Join(t.TempDir(), "cert.pem"), filepath.Join(t.TempDir(), "key.pem")
	require.NoError(t, os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), 0600))
	require.NoError(t, os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}), 0600))
	roots := x509.NewCertPool()
	roots.AddCert(cert)

	ipv6 := supportsListen(t, "tcp6", "[::1]:0")
	for _, scheme := range []string{"http", "https"} {
		for _, tc := range []struct {
			name, host string
			clients    []string
			rejected   string
		}{
			{"IPv4 loopback", "127.0.0.1", []string{"127.0.0.1"}, "::1"},
			{"IPv4 wildcard", "0.0.0.0", []string{"127.0.0.1"}, ""},
			{"empty host", "", []string{"127.0.0.1"}, ""},
			{"IPv6 loopback", "::1", []string{"::1"}, "127.0.0.1"},
			{"bracketed IPv6 loopback", "[::1]", []string{"::1"}, "127.0.0.1"},
			{"IPv6 wildcard", "::", []string{"::1"}, ""},
			{"bracketed IPv6 wildcard", "[::]", []string{"::1"}, ""},
		} {
			t.Run(scheme+"/"+tc.name, func(t *testing.T) {
				if strings.Contains(tc.host, ":") && !ipv6 {
					t.Skip("IPv6 loopback is unavailable")
				}
				// Reserve a free port, then release it for the actual Artalk server.
				probe, err := net.Listen("tcp", net.JoinHostPort(strings.Trim(tc.host, "[]"), "0"))
				require.NoError(t, err)
				port := probe.Addr().(*net.TCPAddr).Port
				// Check which families this system supports for wildcard listeners.
				clients := append([]string{}, tc.clients...)
				if tc.rejected == "" {
					for _, host := range []string{"127.0.0.1", "::1"} {
						if host == clients[0] {
							continue
						}
						conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), time.Second)
						if err == nil {
							conn.Close()
							clients = append(clients, host)
						}
					}
				}
				require.NoError(t, probe.Close())

				app, err := test.NewTestApp()
				require.NoError(t, err)
				// Close the listener before the fixture database and cache.
				app.OnTerminate().RemoveAll()
				t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
				app.Conf().Host, app.Conf().Port = tc.host, port
				app.Conf().SSL.Enabled = scheme == "https"
				app.Conf().SSL.CertPath, app.Conf().SSL.KeyPath = certPath, keyPath
				done := make(chan error, 1)
				go func() { _, err := server.Serve(app.App); done <- err }()
				exited := false
				t.Cleanup(func() {
					require.NoError(t, app.OnTerminate().Trigger(&core.TerminateEvent{App: app.App}))
					if !exited {
						select {
						case err := <-done:
							require.NoError(t, err)
						case <-time.After(5 * time.Second):
							t.Error("Artalk did not stop")
						}
					}
				})

				transport := &http.Transport{DisableKeepAlives: true, TLSClientConfig: &tls.Config{RootCAs: roots}}
				defer transport.CloseIdleConnections()
				client := &http.Client{Transport: transport, Timeout: time.Second}
				endpoint := func(host string) string {
					return scheme + "://" + net.JoinHostPort(host, strconv.Itoa(port)) + "/api/v2/version"
				}
				var startupErr error
				require.Eventually(t, func() bool {
					select {
					case startupErr = <-done:
						exited = true
						return true
					default:
					}
					response, err := client.Get(endpoint(clients[0]))
					if err != nil {
						return false
					}
					response.Body.Close()
					return response.StatusCode == http.StatusOK
				}, 5*time.Second, 20*time.Millisecond, "Artalk did not become reachable")
				require.NoError(t, startupErr)
				require.False(t, exited, "Artalk stopped before accepting requests")
				for _, host := range clients {
					response, err := client.Get(endpoint(host))
					require.NoError(t, err, host)
					body, err := io.ReadAll(response.Body)
					response.Body.Close()
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, response.StatusCode)
					require.Contains(t, string(body), `"version"`)
				}
				if tc.rejected != "" {
					conn, err := net.DialTimeout("tcp", net.JoinHostPort(tc.rejected, strconv.Itoa(port)), time.Second)
					if conn != nil {
						conn.Close()
					}
					require.Error(t, err, "loopback binding must not accept the other address family")
				}
			})
		}
	}
}

func supportsListen(t *testing.T, network, address string) bool {
	t.Helper()
	listener, err := net.Listen(network, address)
	if err != nil {
		t.Logf("%s %s unavailable: %v", network, address, err)
		return false
	}
	require.NoError(t, listener.Close())
	return true
}
