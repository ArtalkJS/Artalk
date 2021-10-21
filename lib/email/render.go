package email

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/ArtalkJS/ArtalkGo/pkged"
)

func RenderEmailTpl(notify *model.Notify, from model.CookedCommentForEmail, to model.CookedCommentForEmail) string {
	tplName := config.Instance.Email.MailTpl
	tpl := ""
	if _, err := os.Stat(tplName); errors.Is(err, os.ErrNotExist) {
		tpl = GetInternalEmailTpl(tplName)
	} else {
		tpl = GetExternalEmailTpl(tplName)
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

	tpl = RenderConfig(tpl, notify, from, to)
	tpl = strings.ReplaceAll(tpl, "{{reply_link}}", notify.GetReadLink())

	return tpl
}

func RenderConfig(str string, notify *model.Notify, from model.CookedCommentForEmail, to model.CookedCommentForEmail) string {
	for k, v := range config.Flat {
		str = strings.ReplaceAll(str, fmt.Sprintf("{{config.%s}}", k), fmt.Sprintf("%v", v))
	}
	str = strings.ReplaceAll(str, "{{site_name}}", from.SiteName)
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
	f, err := pkged.Open(filename)
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
