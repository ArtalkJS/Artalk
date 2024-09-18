package anti_spam

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/log"
)

var _ Checker = (*AkismetChecker)(nil)

type AkismetChecker struct {
	key string
}

func NewAkismetChecker(key string) Checker {
	return &AkismetChecker{
		key: key,
	}
}

func (*AkismetChecker) Name() string {
	return "akismet"
}

func (c *AkismetChecker) Check(p *CheckerParams) (bool, error) {
	// @link https://akismet.com/development/api/#comment-check
	form := url.Values{}

	reqParams := newAkismetReqParams(p)
	v := reflect.ValueOf(*reqParams)
	t := v.Type()
	for i := 0; i < v.Type().NumField(); i++ {
		if v.Field(i).String() != "" {
			form.Add(t.Field(i).Tag.Get("name"), v.Field(i).String())
		}
	}

	client := &http.Client{}

	reqBody := strings.NewReader(form.Encode())
	api := fmt.Sprintf("https://%s.rest.akismet.com/1.1/comment-check", c.key)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	respStr := string(respBody)

	log.Debug("akismet Spam Detection Response ", respStr)

	switch respStr {
	case "true":
		// is a spam comment
		return false, nil
	case "false":
		// not a spam comment
		return true, nil
	}

	return false, fmt.Errorf(respStr)
}

type AkismetReqParams struct {
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

func newAkismetReqParams(params *CheckerParams) *AkismetReqParams {
	return &AkismetReqParams{
		Blog: params.BlogURL,

		UserIP:    params.UserIP,
		UserAgent: params.UserAgent,

		CommentAuthor:      params.UserName,
		CommentAuthorEmail: params.UserEmail,
		CommentContent:     params.Content,
	}
}
