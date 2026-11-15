package email

import (
	"fmt"
	"io"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSMTPDelivery(t *testing.T) {
	for _, accept := range []bool{true, false} {
		t.Run(map[bool]string{true: "accepted", false: "recipient rejected"}[accept], func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			t.Cleanup(func() { listener.Close() })
			require.NoError(t, listener.(*net.TCPListener).SetDeadline(time.Now().Add(5*time.Second)))
			delivered := make(chan string, 1)
			done := make(chan error, 1)
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					done <- err
					return
				}
				defer conn.Close()
				conn.SetDeadline(time.Now().Add(5 * time.Second))
				text := textproto.NewConn(conn)
				text.PrintfLine("220 localhost SMTP")
				for {
					line, err := text.ReadLine()
					if err != nil {
						done <- err
						return
					}
					switch {
					case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"), strings.HasPrefix(line, "MAIL FROM:"):
						text.PrintfLine("250 OK")
					case strings.HasPrefix(line, "RCPT TO:"):
						if !accept {
							text.PrintfLine("550 recipient rejected")
							done <- nil
							return
						}
						text.PrintfLine("250 OK")
					case line == "DATA":
						text.PrintfLine("354 send message")
						raw, err := text.ReadDotBytes()
						if err != nil {
							done <- err
							return
						}
						delivered <- string(raw)
						text.PrintfLine("250 queued")
					case line == "QUIT":
						text.PrintfLine("221 bye")
						done <- nil
						return
					default:
						done <- fmt.Errorf("unexpected SMTP command: %s", line)
						return
					}
				}
			}()
			sender, err := NewSender(EmailConf{EmailConf: config.EmailConf{SendType: config.TypeSMTP, SMTP: config.SMTPConf{Host: "127.0.0.1", Port: listener.Addr().(*net.TCPAddr).Port}}})
			require.NoError(t, err)
			email := &Email{FromAddr: "noreply@example.org", FromName: "Artalk", ToAddr: "reader@example.org", Subject: "Reply notification", Body: "<p>A new reply</p>"}
			require.Equal(t, accept, sender.Send(email))
			require.NoError(t, <-done)
			if accept {
				message, err := mail.ReadMessage(strings.NewReader(<-delivered))
				require.NoError(t, err)
				require.Equal(t, email.ToAddr, message.Header.Get("To"))
				require.Equal(t, email.Subject, message.Header.Get("Subject"))
				body, err := io.ReadAll(quotedprintable.NewReader(message.Body))
				require.NoError(t, err)
				require.Equal(t, email.Body+"\n", string(body))
			} else {
				require.Empty(t, delivered)
			}
		})
	}
}
