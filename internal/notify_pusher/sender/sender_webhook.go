package sender

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/artalkjs/artalk/v2/internal/log"
)

type NotifyWebHookReqBody struct {
	NotifySubject string      `json:"notify_subject"`
	NotifyBody    string      `json:"notify_body"`
	Comment       interface{} `json:"comment"`
	ParentComment interface{} `json:"parent_comment"`
}

// WebHook 发送
func SendWebHook(url string, reqData *NotifyWebHookReqBody) {
	jsonByte, _ := json.Marshal(reqData)
	result, err := http.Post(url, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		log.Error("[WebHook Push] ", "Failed to send msg:", err)
		return
	}

	if result.StatusCode != 200 {
		body, _ := io.ReadAll(result.Body)
		log.Error("[WebHook Push] Failed to send msg:", string(body))
	}

	defer result.Body.Close()
}
