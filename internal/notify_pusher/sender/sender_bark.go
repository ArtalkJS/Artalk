package sender

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/log"
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

	defer result.Body.Close()
}
