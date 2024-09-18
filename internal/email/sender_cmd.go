package email

import (
	"io"
	"os/exec"

	"github.com/artalkjs/artalk/v2/internal/log"
)

// CmdSender implements Sender
type CmdSender struct {
}

var _ Sender = (*CmdSender)(nil)

// NewCmdSender sendmail
func NewCmdSender() *CmdSender {
	return &CmdSender{}
}

func (s *CmdSender) Send(email *Email) bool {
	LogTag := "[EMAIL] [sendmail] "
	msg := getEmailMineTxt(email)

	// 调用系统 sendmail
	sendmail := exec.Command("/usr/sbin/sendmail", "-t", "-oi")
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		log.Error(LogTag, err)
		return false
	}

	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		log.Error(LogTag, err)
		return false
	}

	if err := sendmail.Start(); err != nil {
		log.Error(LogTag, err)
		return false
	}

	if _, err := stdin.Write([]byte(msg)); err != nil {
		log.Error(LogTag, err)
		return false
	}

	if err := stdin.Close(); err != nil {
		log.Error(LogTag, err)
		return false
	}

	sentBytes, _ := io.ReadAll(stdout)
	if err := sendmail.Wait(); err != nil {
		log.Error(LogTag, err)
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Errorf("["+LogTag+"] Exit code is %d", exitError.ExitCode())
		}
		return false
	}

	log.Debug(string(sentBytes))

	return true
}
