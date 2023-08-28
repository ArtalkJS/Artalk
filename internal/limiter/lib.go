package limiter

import (
	"strconv"
	"time"
)

const (
	ActionTimeCachePrefix         = "action-time:"
	ActionCountCachePrefix        = "action-count:"
	AlwaysModeVerifiedCachePrefix = "captcha-am-verified:"
)

// 现在距离上一次操作，是否超时，需要重置
func (l *Limiter) isTimeoutSinceLastAction(ip string) bool {
	if l.conf.ResetTimeout <= 0 {
		return false // 永不会超时
	}
	lastActionTime := l.getLastTime(ip)
	timeSinceLastAction := time.Since(lastActionTime).Seconds()
	return timeSinceLastAction > float64(l.conf.ResetTimeout)
}

// 修改最后操作时间
func (l *Limiter) logTime(ip string) {
	l.store.Set(ActionTimeCachePrefix+ip, strconv.FormatInt(time.Now().Unix(), 10))
}

// 获取最后操作时间
func (l *Limiter) getLastTime(ip string) time.Time {
	var timestamp int64
	if val, found := l.store.Get(ActionTimeCachePrefix + ip); found {
		timestamp, _ = strconv.ParseInt(val.(string), 10, 64)
	}
	return time.Unix(timestamp, 0)
}

// 获取操作次数
func (l *Limiter) getCount(ip string) int {
	if val, found := l.store.Get(ActionCountCachePrefix + ip); found {
		count, _ := strconv.Atoi(val.(string))
		return count
	}
	return 0
}

// 修改操作次数
func (l *Limiter) setCount(ip string, num int) {
	l.store.Set(ActionCountCachePrefix+ip, strconv.Itoa(num))
}

// 操作次数 +1
func (l *Limiter) increaseCount(ip string) {
	l.setCount(ip, l.getCount(ip)+1)
}

// IP 是否验证过验证码 (for 总是需要验证码的选项)
func (l *Limiter) isVerified(ip string) bool {
	val, found := l.store.Get(AlwaysModeVerifiedCachePrefix + ip)
	return found && val == "1"
}

// 设置 IP 是否验证过验证码 (for 总是需要验证码的选项)
func (l *Limiter) setVerified(ip string, pass bool) {
	val := "0"
	if pass {
		val = "1"
	}
	l.store.Set(AlwaysModeVerifiedCachePrefix+ip, val)
}
