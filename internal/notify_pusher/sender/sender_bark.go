package sender

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/log"
)

// Bark 发送
func SendBark(serverURL string, title string, msg string) {
	if title == "" {
		title = "Artalk"
	}

	result, err := http.Get(fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(serverURL, "/"), url.QueryEscape(title), url.QueryEscape(msg)))
	if err != nil {
		log.Error("[Bark] ", "Failed to send msg:", err)
		return
	}

	if result.StatusCode != 200 {
		body, _ := io.ReadAll(result.Body)
		log.Error("[Bark] Failed to send msg:", string(body))
	}

	defer result.Body.Close()
}
