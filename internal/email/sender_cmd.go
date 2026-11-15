package email

import (
	"os/exec"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/log"
)

// CmdSender sends messages through the system sendmail command.
type CmdSender struct{}

var _ Sender = (*CmdSender)(nil)

func NewCmdSender() *CmdSender { return &CmdSender{} }

func (s *CmdSender) Send(email *Email) bool {
	return sendEmailCommand(email, exec.Command("/usr/sbin/sendmail", "-t", "-oi"))
}

func sendEmailCommand(email *Email, command *exec.Cmd) bool {
	command.Stdin = strings.NewReader(getEmailMineTxt(email))
	// CombinedOutput closes the pipes and waits for the process on all paths.
	output, err := command.CombinedOutput()
	if err != nil {
		log.Error("[Email] Email sending failed via sendmail: ", err)
		return false
	}
	log.Debug(string(output))
	return true
}
