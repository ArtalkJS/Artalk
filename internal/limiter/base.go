package limiter

import (
	"github.com/artalkjs/artalk/v2/internal/cache/simple_cache"
)

// 操作次数限制器（根据 IP）
//
// 关键函数：
//
//  1. IsPass(ip)
//     请求启用了 Limiter 的页面时，
//     判断操作次数是否过限，若超过限制则响应 need_captcha
//
//  2. MarkVerifyPassed
//     若验证码正确，则放行操作一次（最大操作数 -1）
//
//  3. MarkVerifyFailed
//     若验证码错误，则操作次数 +1，记录最后操作时间
//
//  4. Log(ip)
//     当用户发表评论后，调用 Log 操作次数 +1，记录最后操作时间
type Limiter struct {
	conf  *LimiterConf
	store *simple_cache.Cache // simple_cache is a thread-safe in-memory cache
}

type LimiterConf struct {
	AlwaysMode          bool // 总是需要验证码模式 (允许零次操作时开启)
	MaxActionDuringTime int  // 时间范围内允许多少次操作 (激活验证码所需调用 Log 函数次数)
	ResetTimeout        int  // 重置超时，单位 s (当你想当评论 n 次后一直需要验证码，可以将时间设置得足够大，或是 -1)
}

func NewLimiter(conf *LimiterConf) *Limiter {
	if conf.MaxActionDuringTime <= 0 {
		conf.AlwaysMode = true
	}

	return &Limiter{
		conf:  conf,
		store: simple_cache.New(),
	}
}

// 请求是否需要验证码
//
// Notice: call IsPass will trigger a write operation.
func (l *Limiter) IsPass(ip string) bool {
	// =======================
	//  总是需要验证码模式
	// =======================
	if l.conf.AlwaysMode {
		return l.isVerified(ip)
	}

	// =======================
	//  时间范围内允许多少次操作
	// =======================
	if l.getCount(ip) >= l.conf.MaxActionDuringTime { // 当前已操作次数达到限制

		// 配置了过期时间，并且操作在时间窗口外（操作超时了）则执行 Reset
		// 若未配置超时时间，则永不会超时，即用不会执行 Reset
		if l.isTimeoutSinceLastAction(ip) {
			l.ResetLog(ip) // 重置计数 (write)
			return true    // 放行
		}

		return false // 阻止

	} else {
		return true // 操作未达到限制，放行
	}
}

// 记录操作
// (请勿在 IsNeedVerify 函数被调用之前执行 Log)
func (l *Limiter) Log(ip string) {
	l.increaseCount(ip)      // 操作次数 +1
	l.logTime(ip)            // 更新最后操作时间
	l.setVerified(ip, false) // 重置已验证状态为 false
}

// 重置操作记录
func (l *Limiter) ResetLog(ip string) {
	l.store.Delete(ActionTimeCachePrefix + ip)
	l.store.Delete(ActionCountCachePrefix + ip)
}

// 验证成功操作
func (l *Limiter) MarkVerifyPassed(ip string) {
	l.setCount(ip, l.conf.MaxActionDuringTime-1) // 仅放行操作一次（而不是归零）
	l.setVerified(ip, true)
}

// 验证失败操作
func (l *Limiter) MarkVerifyFailed(ip string) {
	l.Log(ip) // 记录操作
	l.setVerified(ip, false)
}
