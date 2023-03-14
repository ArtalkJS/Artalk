package aliyun

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AliyunApiRequestConf struct {
	AccessKeyId     string
	AccessKeySecret string

	BaseURL string
	Path    string
	Version string

	Data map[string]interface{}
}

func RequestApi(conf AliyunApiRequestConf) ([]byte, error) {
	dataJson := jsonEncode(conf.Data)
	date := time.Now().UTC().Format(time.RFC1123)
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	contentMD5 := md5Base64(dataJson)
	headers := map[string]string{
		"Accept":                  "application/json",
		"Content-Type":            "application/json",
		"Content-MD5":             contentMD5,
		"Date":                    date,
		"x-acs-version":           conf.Version,
		"x-acs-signature-nonce":   nonce,
		"x-acs-signature-version": "1.0",
		"x-acs-signature-method":  "HMAC-SHA1",
	}

	rawData := strings.Join([]string{
		"POST",
		"application/json",
		contentMD5,
		"application/json",
		date,
		getSignHeaders(headers),
		conf.Path,
	}, "\n")

	headers["Authorization"] = fmt.Sprintf("acs %v:%v", conf.AccessKeyId, hmacSHA1Base64(conf.AccessKeySecret, rawData))

	client := &http.Client{}
	req, err := http.NewRequest("POST", conf.BaseURL+conf.Path, strings.NewReader(dataJson))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getSignHeaders(headers map[string]string) string {
	tmp := map[string]string{}
	for k, v := range headers {
		if strings.HasPrefix(k, "x-acs-") {
			tmp[k] = v
		}
	}

	keys := []string{}
	for k := range tmp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result []string
	for _, k := range keys {
		result = append(result, k+":"+tmp[k])
	}

	return strings.Join(result, "\n")
}

func hmacSHA1Base64(key string, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func md5Base64(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func jsonEncode(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
