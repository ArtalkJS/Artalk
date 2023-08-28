package email

import (
	"bytes"

	"gopkg.in/gomail.v2"
)

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
	// 内容
	m.SetBody("text/html", email.Body)
	// 附件
	//m.Attach("./file.png")

	return m
}

func getEmailMineTxt(email *Email) string {
	emailBuffer := bytes.NewBuffer([]byte{})
	getCookedEmail(email).WriteTo(emailBuffer)
	return string(emailBuffer.Bytes()[:])
}
