package qcloud

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func RequestApi(conf TencentApiRequestConf) ([]byte, error) {
	if conf.Region == "" {
		conf.Region = "ap-guangzhou"
	}

	var (
		host      = fmt.Sprintf("%s.tencentcloudapi.com", conf.Product)
		endpoint  = fmt.Sprintf("https://%s", host)
		algorithm = "TC3-HMAC-SHA256"
		timestamp = time.Now().Unix()
		date      = time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	)

	// ************* 步骤 1：拼接规范请求串 *************
	var (
		httpReqMethod  = conf.Method
		canonicalURI   = "/"
		payload        string
		canonicalQuery string
		contentType    string
	)

	switch conf.Method {
	case "GET":
		payload = ""
		canonicalQuery = join(conf.Data)
		contentType = "application/x-www-form-urlencoded; charset=utf-8"
	case "POST":
		payloadBytes, _ := json.Marshal(conf.Data)
		payload = string(payloadBytes)
		canonicalQuery = ""
		contentType = "application/json; charset=utf-8"
	}

	// fmt.Println(canonicalQuery)
	var (
		canonicalHeaders = fmt.Sprintf("content-type:%s\nhost:%s\n", contentType, host)
		signedHeaders    = "content-type;host"
		hashedReqPayload = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	)
	canonicalReq := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		httpReqMethod, canonicalURI, canonicalQuery,
		canonicalHeaders, signedHeaders, hashedReqPayload,
	)

	// ************* 步骤 2：拼接待签名字符串 *************
	var (
		credentialScope    = fmt.Sprintf("%s/%s/tc3_request", date, conf.Product)
		hashedCanonicalReq = fmt.Sprintf("%x", sha256.Sum256([]byte(canonicalReq)))
	)
	stringToSign := fmt.Sprintf(
		"%s\n%d\n%s\n%s",
		algorithm, timestamp, credentialScope, hashedCanonicalReq,
	)

	// ************* 步骤 3：计算签名 *************
	var (
		secretDate    = sign([]byte(fmt.Sprintf("TC3%s", conf.SecretKey)), date)
		secretService = sign(secretDate, conf.Product)
		secretSigning = sign(secretService, "tc3_request")
		signature     = hex.EncodeToString(sign(secretSigning, stringToSign))
	)

	// ************* 步骤 4：拼接 Authorization *************
	authorization := fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, conf.SecretID, credentialScope, signedHeaders, signature,
	)
	headers := map[string]string{
		"Authorization":  authorization,
		"Content-Type":   contentType,
		"Host":           host,
		"X-TC-Action":    conf.Action,
		"X-TC-Timestamp": fmt.Sprintf("%d", timestamp),
		"X-TC-Version":   conf.Version,
		"X-TC-Region":    conf.Region,
	}

	// 发送 Http 请求
	var (
		client  = &http.Client{}
		reqUrl  string
		reqBody *bytes.Buffer
	)
	switch httpReqMethod {
	case "GET":
		reqUrl = endpoint + "?" + join(conf.Data)
		reqBody = bytes.NewBuffer([]byte(""))
	case "POST":
		reqUrl = endpoint
		jsonData, _ := json.Marshal(conf.Data)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, _ := http.NewRequest(httpReqMethod, reqUrl, reqBody)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// fmt.Println(req.Header)
	// fmt.Println(resp.Status)

	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type TencentApiRequestConf struct {
	SecretID  string
	SecretKey string
	Region    string

	Product string
	Action  string
	Version string
	Method  string

	Data map[string]string
}

func join(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values.Encode()
}

func sign(key []byte, msg string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(msg))
	return h.Sum(nil)
}

func base64EncodeStr(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
