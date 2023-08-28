package sender

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/log"
)

// 飞书发送
func SendLark(webhookURL string, title string, msg string) {
	if title != "" {
		msg = title + "\n\n" + msg
	}

	sendData := fmt.Sprintf(`{"msg_type":"text","content":{"text":%s}}`, strconv.Quote(msg))
	result, err := http.Post(webhookURL, "application/json", strings.NewReader(sendData))
	if err != nil {
		log.Error("[飞书] ", "Failed to send msg:", err)
		return
	}

	defer result.Body.Close()
}
