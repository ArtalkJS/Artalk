package captcha

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/tidwall/gjson"
)

// geetest 服务地址
const GEETEST_API_SERVER string = "http://gcaptcha4.geetest.com"

// hmac-sha256 加密：  CAPTCHA_KEY,lot_number
func hmac_encode(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

type GeetestParams struct {
	CaptchaID     string `json:"captcha_id"`
	LotNumber     string `json:"lot_number"`
	PassToken     string `json:"pass_token"`
	GenTime       string `json:"gen_time"`
	CaptchaOutput string `json:"captcha_output"`
}

func GeetestCheck(paramsJSON string) (isPass bool, reason string, err error) {
	geetestConf := config.Instance.Captcha.Geetest

	var p GeetestParams
	err = json.Unmarshal([]byte(paramsJSON), &p)
	if err != nil {
		return false, "", errors.New("Reqest params json parse err: " + err.Error())
	}

	// 生成签名
	signToken := hmac_encode(geetestConf.CaptchaKey, p.LotNumber)

	// 向极验转发前端数据 + sign_token 签名
	form_data := make(url.Values)
	form_data["lot_number"] = []string{p.LotNumber}
	form_data["captcha_output"] = []string{p.CaptchaOutput}
	form_data["pass_token"] = []string{p.PassToken}
	form_data["gen_time"] = []string{p.GenTime}
	form_data["sign_token"] = []string{signToken}

	// 发起 POST 请求
	url := GEETEST_API_SERVER + "/validate" + "?captcha_id=" + geetestConf.CaptchaID
	cli := http.Client{Timeout: time.Second * 5} // 5s 超时
	resp, err := cli.PostForm(url, form_data)
	if err != nil || resp.StatusCode != 200 {
		return false, "", errors.New("服务接口异常: " + err.Error())
	}

	// 处理响应结果
	resJsonBuf, _ := ioutil.ReadAll(resp.Body)
	resJson := string(resJsonBuf)

	gResult := gjson.Get(resJson, "result")
	if !gResult.Exists() {
		return false, "", errors.New("响应结果不符合预期: " + resJson)
	}
	result := gResult.String()

	if result == "success" {
		// 验证成功
		return true, "", nil
	} else {
		// 验证失败
		gReason := gjson.Get(resJson, "reason")
		if gReason.Exists() {
			reason = gReason.String()
		}

		return false, reason, nil
	}
}
