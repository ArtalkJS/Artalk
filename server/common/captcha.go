package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/gofiber/fiber/v2"
	imgCaptcha "github.com/steambap/captcha"
)

var (
	CaptchaExpiration = 5 * time.Minute // 验证码 5 分钟内有效
)

// 请求是否需要验证码
func IsReqNeedCaptchaCheck(c *fiber.Ctx) bool {
	captchaConf := config.Instance.Captcha

	// 管理员直接忽略
	if CheckIsAdminReq(c) {
		return false
	}

	// 总是需要验证码模式
	if config.Instance.Captcha.Always {
		return !GetAlwaysCaptchaMode_Pass(c.IP())
	}

	// 不重置计数器模式
	if captchaConf.ActionReset == -1 {
		if getActionCount(c) >= captchaConf.ActionLimit { // 只要操作次数超过
			return true // 就过限
		} else {
			return false // 放行
		}
	}

	// 开启重置计数器功能的情况：在时间范围内，操作次数超过
	if IsActionInTimeFrame(c) && getActionCount(c) >= captchaConf.ActionLimit {
		return true // 过限
	} else {
		return false // 放行
	}
}

// 记录操作
func RecordAction(c *fiber.Ctx) {
	updateActionLastTime(c) // 更新最后操作时间
	addActionCount(c)       // 操作次数 +1
}

// 重置操作记录
func ResetActionRecord(c *fiber.Ctx) {
	ip := c.IP()

	cache.CACHE.Delete(cache.Ctx, "action-time:"+ip)
	cache.CACHE.Delete(cache.Ctx, "action-count:"+ip)
}

// 操作计数是否应该被重置
func IsActionInTimeFrame(c *fiber.Ctx) bool {
	return time.Since(getActionLastTime(c)).Seconds() <= float64(config.Instance.Captcha.ActionReset)
}

// 修改最后操作时间
func updateActionLastTime(c *fiber.Ctx) {
	curtTime := fmt.Sprintf("%v", time.Now().Unix())
	cache.CACHE.Set(cache.Ctx, "action-time:"+c.IP(), curtTime)
}

// 获取最后操作时间
func getActionLastTime(c *fiber.Ctx) time.Time {
	var timestamp int64
	var val string
	if _, err := cache.CACHE.Get(cache.Ctx, "action-time:"+c.IP(), &val); err == nil {
		timestamp, _ = strconv.ParseInt(string(val), 10, 64)
	}
	tm := time.Unix(timestamp, 0)
	return tm
}

// 获取操作次数
func getActionCount(c *fiber.Ctx) int {
	count := 0
	var val string
	if _, err := cache.CACHE.Get(cache.Ctx, "action-count:"+c.IP(), &val); err == nil {
		count, _ = strconv.Atoi(val)
	}

	return count
}

// 修改操作次数
func setActionCount(c *fiber.Ctx, num int) {
	cache.CACHE.Set(cache.Ctx, "action-count:"+c.IP(), fmt.Sprintf("%d", num))
}

// 操作次数 +1
func addActionCount(c *fiber.Ctx) {
	setActionCount(c, getActionCount(c)+1)
}

// 验证成功操作
func OnCaptchaPass(c *fiber.Ctx) {
	ip := c.IP()

	setActionCount(c, config.Instance.Captcha.ActionLimit-1)
	SetAlwaysCaptchaMode_Pass(ip, true) // 允许 always mode pass
}

// 验证失败操作
func OnCaptchaFail(c *fiber.Ctx) {
	ip := c.IP()

	RecordAction(c)                      // 记录操作
	SetAlwaysCaptchaMode_Pass(ip, false) // 取消 always mode pass
}

// #region 图片验证码
// 获取对应 IP 图片验证码正确的值
func GetImageCaptchaRealCode(ip string) string {
	var realVal string
	cache.CACHE.Get(cache.Ctx, "captcha:"+ip, &realVal)
	return strings.ToLower(realVal)
}

// 获取新验证码 base64 格式图片
func GetNewImageCaptchaBase64(ip string) string {
	// generate a image
	pngBuffer := bytes.NewBuffer([]byte{})
	data, _ := imgCaptcha.New(160, 40, func(o *imgCaptcha.Options) {
		o.FontScale = 1
		o.CurveNumber = 2
		o.FontDPI = 85.0
		o.Noise = 0.7
		o.BackgroundColor = color.White
	})
	data.WriteImage(pngBuffer)
	base64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	// save real code
	cache.CACHE.Set(cache.Ctx, "captcha:"+ip, data.Text, store.WithExpiration(CaptchaExpiration))

	return base64
}

// 销毁图片验证码
func DisposeImageCaptcha(ip string) {
	cache.CACHE.Delete(cache.Ctx, "captcha:"+ip)
}

//#endregion

// AlwaysMode 是否能 Pass (for 总是需要验证码的选项)
func GetAlwaysCaptchaMode_Pass(ip string) bool {
	var val string
	_, err := cache.CACHE.Get(cache.Ctx, "captcha-am-pass:"+ip, &val)
	return err == nil && val == "1"
}

// 设置 AlwaysMode 允许 Pass (for 总是需要验证码的选项)
func SetAlwaysCaptchaMode_Pass(ip string, pass bool) {
	val := "0"
	if pass {
		val = "1"
	}

	cache.CACHE.Set(cache.Ctx, "captcha-am-pass:"+ip, val, store.WithExpiration(CaptchaExpiration))
}
