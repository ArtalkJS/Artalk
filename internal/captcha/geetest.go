package captcha

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/tidwall/gjson"
)

const GEETEST_API string = "http://gcaptcha4.geetest.com"

type GeetestCaptcha struct {
	CaptchaID  string
	CaptchaKey string
}

var _ Captcha = (*GeetestCaptcha)(nil)

func NewGeetestCaptcha(conf *config.GeetestConf) *GeetestCaptcha {
	return &GeetestCaptcha{
		CaptchaID:  conf.CaptchaID,
		CaptchaKey: conf.CaptchaKey,
	}
}

type GeetestParams struct {
	CaptchaID     string `json:"captcha_id"`
	LotNumber     string `json:"lot_number"`
	PassToken     string `json:"pass_token"`
	GenTime       string `json:"gen_time"`
	CaptchaOutput string `json:"captcha_output"`
}

func (c *GeetestCaptcha) Check(p CaptchaPayload) (bool, error) {
	var ck GeetestParams
	if err := json.Unmarshal([]byte(p.CheckValue), &ck); err != nil {
		return false, err
	}

	// 生成签名
	signToken := hmacEncode(c.CaptchaKey, ck.LotNumber)

	// 向极验转发前端数据 + sign_token 签名
	values := make(url.Values)
	values.Add("lot_number", ck.LotNumber)
	values.Add("captcha_output", ck.CaptchaOutput)
	values.Add("pass_token", ck.PassToken)
	values.Add("gen_time", ck.GenTime)
	values.Add("sign_token", signToken)

	// 发起 POST 请求
	url := GEETEST_API + "/validate?captcha_id=" + c.CaptchaID
	cli := http.Client{Timeout: time.Second * 10} // 10s 超时
	resp, err := cli.PostForm(url, values)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}
	defer resp.Body.Close()

	// 处理响应结果
	respBuf, _ := io.ReadAll(resp.Body)

	if gjson.GetBytes(respBuf, "result").String() == "success" {
		// 验证成功
		return true, nil
	} else {
		// 验证失败
		return false, errors.New("err reason: " + gjson.GetBytes(respBuf, "reason").String())
	}
}

func (c *GeetestCaptcha) PageParams() Map {
	return Map{
		"gt_id": c.CaptchaID,
	}
}
