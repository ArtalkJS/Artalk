package email

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/markbates/pkger"
)

func RenderEmailTpl(from model.CookedCommentForEmail, to model.CookedCommentForEmail) string {
	tplName := config.Instance.Email.MailTpl
	tpl := ""
	if _, err := os.Stat(tplName); os.IsExist(err) {
		tpl = GetExternalEmailTpl(tplName)
	} else {
		tpl = GetInternalEmailTpl(tplName)
	}

	getValue := func(k string, v interface{}) string {
		val := fmt.Sprintf("%v", v)
		if k != "content" {
			val = html.EscapeString(val)
		}
		return val
	}

	flatFrom := lib.StructToFlatDotMap(&from)
	flatTo := lib.StructToFlatDotMap(&to)
	for k, v := range flatFrom {
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{{from.%s}}", k), getValue(k, v))
	}
	for k, v := range flatTo {
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{{to.%s}}", k), getValue(k, v))
	}

	tpl = RenderConfig(tpl)
	tpl = strings.ReplaceAll(tpl, "{{reply_link}}", lib.AddQueryToURL(from.Link, map[string]string{"artalk_comment": fmt.Sprintf("%d", from.ID)}))

	return tpl
}

func RenderConfig(str string) string {
	for k, v := range config.Flat {
		str = strings.ReplaceAll(str, fmt.Sprintf("{{config.%s}}", k), fmt.Sprintf("%v", v))
	}
	str = strings.ReplaceAll(str, "{{site_name}}", config.Instance.SiteName)
	return str
}

func GetExternalEmailTpl(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(buf)
}

func GetInternalEmailTpl(tplName string) string {
	filename := fmt.Sprintf("/email-tpl/%s.html", tplName)
	f, err := pkger.Open(filename)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return ""
	}
	contents := buf.String()

	return contents
}
