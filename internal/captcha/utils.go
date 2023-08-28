package captcha

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// hmac-sha256 加密
func hmacEncode(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
