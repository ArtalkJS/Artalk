package email

import (
	"bytes"
	"net/mail"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/utils"
	"gopkg.in/gomail.v2"
)

const fallbackMessageIDDomain = "artalk.local"

func getCookedEmail(email *Email) *gomail.Message {
	m := gomail.NewMessage()

	// 发送人
	m.SetHeader("From", m.FormatAddress(email.FromAddr, email.FromName))
	// 接收人
	m.SetHeader("To", email.ToAddr)
	// 抄送人
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	// 主题
	m.SetHeader("Subject", email.Subject)
	// 唯一邮件标识
	m.SetHeader("Message-ID", newMessageID(email.FromAddr))
	// 内容
	m.SetBody("text/html", email.Body)
	// 附件
	//m.Attach("./file.png")

	return m
}

func newMessageID(fromAddr string) string {
	domain := fallbackMessageIDDomain
	if addr, err := mail.ParseAddress(fromAddr); err == nil {
		if i := strings.LastIndexByte(addr.Address, '@'); i >= 0 && i < len(addr.Address)-1 {
			domain = addr.Address[i+1:]
		}
	}

	return "<" + utils.RandomString(32) + "@" + domain + ">"
}

func getEmailMineTxt(email *Email) string {
	emailBuffer := bytes.NewBuffer([]byte{})
	getCookedEmail(email).WriteTo(emailBuffer)
	return string(emailBuffer.Bytes()[:])
}
