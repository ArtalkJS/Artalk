package email

import (
	"io"
	"os/exec"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/sirupsen/logrus"
)

// CmdSender implements Sender
type CmdSender struct {
}

var _ Sender = (*CmdSender)(nil)

// NewCmdSender sendmail
func NewCmdSender() *CmdSender {
	return &CmdSender{}
}

func (s *CmdSender) Send(email Email) bool {
	LogTag := "[EMAIL] [sendmail] "
	msg := getEmailMineTxt(email)

	// 调用系统 sendmail
	sendmail := exec.Command("/usr/sbin/sendmail", "-t", "-oi")
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		logrus.Error(LogTag, err)
		return false
	}

	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		logrus.Error(LogTag, err)
		return false
	}

	if err := sendmail.Start(); err != nil {
		logrus.Error(LogTag, err)
		return false
	}

	if _, err := stdin.Write([]byte(msg)); err != nil {
		logrus.Error(LogTag, err)
		return false
	}

	if err := stdin.Close(); err != nil {
		logrus.Error(LogTag, err)
		return false
	}

	sentBytes, _ := io.ReadAll(stdout)
	if err := sendmail.Wait(); err != nil {
		logrus.Error(LogTag, err)
		if exitError, ok := err.(*exec.ExitError); ok {
			logrus.Error(LogTag, "Exit code is %d", exitError.ExitCode())
		}
		return false
	}

	if config.Instance.Debug {
		logrus.Debug(string(sentBytes))
	}

	return true
}
