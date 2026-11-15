package email

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSendmailCommandDelivery(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sendmail requires a Unix host")
	}
	sender, err := NewSender(EmailConf{EmailConf: config.EmailConf{SendType: config.TypeSendmail}})
	require.NoError(t, err)
	require.IsType(t, &CmdSender{}, sender)
	message := &Email{FromAddr: "noreply@example.org", ToAddr: "reader@example.org", Subject: "Test delivery", Body: "<p>Reply</p>"}
	output := filepath.Join(t.TempDir(), "message")
	command := exec.Command("sh", "-c", `cat > "$1"`, "sendmail-fixture", output)
	require.True(t, sendEmailCommand(message, command))
	data, err := os.ReadFile(output)
	require.NoError(t, err)
	parsed := readEmailMessage(t, string(data))
	require.Equal(t, message.Subject, parsed.Header.Get("Subject"))
	require.Equal(t, message.ToAddr, parsed.Header.Get("To"))
	require.True(t, command.ProcessState.Exited())
	for _, command := range []*exec.Cmd{exec.Command(filepath.Join(t.TempDir(), "missing")), exec.Command("sh", "-c", "exit 1")} {
		require.False(t, sendEmailCommand(message, command))
		if command.ProcessState != nil {
			require.True(t, command.ProcessState.Exited())
		}
	}
}
