package ip_region

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"golang.org/x/text/language"
)

func ipScraper(ip string) string {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ""
	}

	// 多 IP 仅选第一个
	ipSep := strings.Split(ip, ",")
	if len(ipSep) > 1 {
		ip = strings.TrimSpace(ipSep[0])
	}

	return ip
}

func regionScraper(raw string, precision config.IPRegionPrecision) string {
	if precision == "" {
		precision = config.IPRegionProvince
	}

	sep := strings.Split(raw, "|")
	if len(sep) < 5 {
		return ""
	}

	country := sep[0]
	if country == "" || country == "0" {
		return ""
	}

	// Legacy: country|area|province|city|isp.
	// Current: country|province|city|isp|ISO country code.
	// Both have five fields. Validate the code rather than treating any
	// two-letter ISP (e.g. "XL" in legacy data) as a country code.
	province, city := sep[2], sep[3]
	if len(sep[4]) == 2 {
		if code, err := language.ParseRegion(sep[4]); err == nil && code.IsCountry() {
			province, city = sep[1], sep[2]
		}
	}
	if province == "0" {
		province = ""
	}
	if city == "0" {
		city = ""
	}
	province = strings.TrimSuffix(strings.TrimSuffix(province, "省"), "市")
	city = strings.TrimSuffix(city, "市")

	switch precision {
	case config.IPRegionCountry:
		return country

	case config.IPRegionProvince:
		if province == "" {
			return country
		}
		return province

	case config.IPRegionCity:
		if province == "" && city == "" {
			return country
		}
		if city == province { // e.g. 重庆重庆
			return province
		}
		return province + city

	default:
		return ""
	}
}
