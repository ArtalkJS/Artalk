package sender

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
)

// 飞书发送
func SendLark(webhookURL string, title string, content string, isCard bool) {
	if title != "" {
		content = title + "\n\n" + content
	}

	var data string
	if isCard {
		texts := strings.Split(strings.TrimSpace(content), "\n")
		link := texts[len(texts)-1]
		if utils.ValidateURL(link) {
			content = strings.Join(texts[1:len(texts)-1], "\n")
		} else {
			content = strings.Join(texts[1:], "\n")
			link = ""
		}
		data = getLarkCardJson(texts[0], content, link)
	} else {
		data = getLarkTextJson(content)
	}

	result, err := http.Post(webhookURL, "application/json", strings.NewReader(data))
	if err != nil {
		log.Error("[飞书] Failed to send msg:", err)
		return
	}

	if result.StatusCode != 200 {
		body, _ := io.ReadAll(result.Body)
		log.Error("[飞书] Failed to send msg:", string(body), " ", string(data))
	}

	defer result.Body.Close()
}

func getLarkTextJson(text string) string {
	return fmt.Sprintf(`{"msg_type":"text","content":{"text":%s}}`, strconv.Quote(text))
}

func getLarkCardJson(title string, content string, replyLink string) string {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	buttonCode := ""
	if replyLink != "" {
		buttonCode = fmt.Sprintf(`,{"tag":"action","actions":[{"tag":"button","text":{"content":%s,"tag":"lark_md"},"url":%s,"type":"default","value":{}}]}`,
			strconv.Quote(i18n.T("Reply")), strconv.Quote(replyLink))
	}

	return fmt.Sprintf(`{"msg_type":"interactive","card":{"header":{"title":{"content":%s,"tag":"plain_text"}},"elements":[{"tag":"div","text":{"content":%s,"tag":"lark_md"}}%s]}}`,
		strconv.Quote(title), strconv.Quote(content), buttonCode)
}
