package anti_spam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/sirupsen/logrus"
)

type AkismetParams struct {
	Blog      string `name:"blog"`       // required
	UserIP    string `name:"user_ip"`    // required
	UserAgent string `name:"user_agent"` // required

	CommentType        string `name:"comment_type"`
	CommentAuthor      string `name:"comment_author"`
	CommentAuthorEmail string `name:"comment_author_email"`
	CommentAuthorURL   string `name:"comment_author_url"`
	CommentContent     string `name:"comment_content"`

	UserRole    string `name:"user_role"`
	Referrer    string `name:"referrer"`
	Permalink   string `name:"permalink"`
	BlogLang    string `name:"blog_lang"`
	BlogCharset string `name:"blog_charset"`
}

// @link https://akismet.com/development/api/#comment-check
func Akismet(params *AkismetParams, key string) (isPass bool, err error) {
	form := url.Values{}

	v := reflect.ValueOf(*params)
	t := v.Type()
	for i := 0; i < v.Type().NumField(); i++ {
		if v.Field(i).String() != "" {
			form.Add(t.Field(i).Tag.Get("name"), v.Field(i).String())
		}
	}

	client := &http.Client{}

	reqBody := strings.NewReader(form.Encode())
	api := fmt.Sprintf("https://%s.rest.akismet.com/1.1/comment-check", key)
	req, err := http.NewRequest("POST", api, reqBody)
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	respStr := string(respBody)

	if config.Instance.Debug {
		logrus.Info("akismet 垃圾检测响应 ", respStr)
	}

	switch respStr {
	case "true":
		// 是垃圾评论
		isPass = false
		return isPass, nil
	case "false":
		// 不是垃圾评论
		isPass = true
		return isPass, nil
	}

	return false, errors.New(respStr)
}
