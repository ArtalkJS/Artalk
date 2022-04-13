package email

import (
	"bytes"
	"io/ioutil"
	"os/exec"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
	ali_dm "github.com/qwqcode/go-aliyun-email"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Email struct {
	FromAddr     string
	FromName     string
	ToAddr       string
	Subject      string
	Body         string
	LinkedNotify *model.Notify
}

func SendBySMTP(email Email) bool {
	smtp := config.Instance.Email.SMTP

	m := GetCookedEmail(email)
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		logrus.Error("[EMAIL] SMTP 邮件发送失败 ", err)
		return false
	}
	return true
}

func SendByAliDM(email Email) bool {
	client := ali_dm.NewClient(
		config.Instance.Email.AliDM.AccessKeyId,
		config.Instance.Email.AliDM.AccessKeySecret,
		config.Instance.Email.AliDM.AccountName,
		email.FromName,
		ali_dm.RegionCNHangZhou,
	)
	req := &ali_dm.SingleRequest{
		ReplyToAddress: true,
		AddressType:    1,
		ToAddress:      email.ToAddr,
		Subject:        email.Subject,
		HtmlBody:       email.Body,
	}

	resp, err := client.SingleRequest(req)
	if err != nil {
		logrus.Error("[EMAIL] ali_dm 邮件发送失败 ", err)
		return false
	}

	if config.Instance.Debug {
		logrus.Debug(resp)
	}

	return true
}

func SendByUsingSystemCMD(email Email) bool {
	LogTag := "[EMAIL] [sendmail] "
	msg := GetEmailMineTxt(email)

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

	sendmail.Start()
	stdin.Write([]byte(msg))
	stdin.Close()
	sentBytes, _ := ioutil.ReadAll(stdout)
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

func GetCookedEmail(email Email) *gomail.Message {
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

func GetEmailMineTxt(email Email) string {
	emailBuffer := bytes.NewBuffer([]byte{})
	GetCookedEmail(email).WriteTo(emailBuffer)
	return string(emailBuffer.Bytes()[:])
}
