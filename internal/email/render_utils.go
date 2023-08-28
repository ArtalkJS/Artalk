package email

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/utils"
)

//go:embed email_tpl/*
//go:embed notify_tpl/*
var internalTpl embed.FS

// 替换 {{ key }} 为 val
func ReplaceAllMustache(data string, dict map[string]interface{}) string {
	return utils.RenderMustaches(data, dict, func(k string, v interface{}) string {
		return GetPurifiedValue(k, v)
	})
}

// 净化文本，防止 XSS
func GetPurifiedValue(k string, v interface{}) string {
	val := fmt.Sprintf("%v", v)

	// 白名单
	ignoreEscapeKeys := []string{"reply_content", "content", "link_to_reply"}
	if utils.ContainsStr(ignoreEscapeKeys, k) ||
		strings.HasSuffix(k, ".content") || // 排除 entity.CookedComment.content
		strings.HasSuffix(k, ".content_raw") {
		return val
	}

	val = html.EscapeString(val)
	return val
}

func HandleEmoticonsImgTagsForNotify(str string) string {
	r := regexp.MustCompile(`<img\s[^>]*?atk-emoticon=["]([^"]*?)["][^>]*?>`)
	return r.ReplaceAllStringFunc(str, func(m string) string {
		ms := r.FindStringSubmatch(m)
		if len(ms) < 2 {
			return m
		}
		if ms[1] == "" {
			return "[表情]"
		}
		return "[" + ms[1] + "]"
	})
}

func GetMailTpl(tplName string) string {
	// 配置文件未指定邮件模板路径，使用内置默认模板
	if tplName == "" {
		tplName = "default"
	}

	var tpl string
	if !utils.CheckFileExist(tplName) {
		tpl = GetInternalEmailTpl(tplName)
	} else {
		// TODO 反复文件 IO 操作会导致性能下降，
		// 之后优化可以改成程序启动时加载模板文件到内存中
		tpl = GetExternalTpl(tplName)
	}

	return tpl
}

func GetNotifyTpl(tplName string) string {
	if tplName == "" {
		tplName = "default"
	}

	var tpl string
	if !utils.CheckFileExist(tplName) {
		tpl = GetInternalNotifyTpl(tplName)
	} else {
		tpl = GetExternalTpl(tplName)
	}

	return tpl
}

// 获取内建邮件模版
func GetInternalEmailTpl(tplName string) string {
	return GetInternalTpl("email_tpl", tplName)
}

// 获取内建通知模版
func GetInternalNotifyTpl(tplName string) string {
	return GetInternalTpl("notify_tpl", tplName)
}

// 获取内建模版
func GetInternalTpl(basePath string, tplName string) string {
	filename := fmt.Sprintf("%s/%s.html", basePath, tplName)
	f, err := internalTpl.Open(filename)
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

// 获取外置模版
func GetExternalTpl(filename string) string {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(buf)
}
