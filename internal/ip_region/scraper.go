package ip_region

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
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

	var (
		country  = sep[0]
		_        = sep[1] // Area
		province = strings.TrimSuffix(strings.TrimSuffix(sep[2], "省"), "市")
		city     = strings.TrimSuffix(sep[3], "市")
		// isp      = sep[4]
	)

	if country == "0" {
		return ""
	}

	switch precision {
	case config.IPRegionCountry:
		return country

	case config.IPRegionProvince:
		if province == "0" {
			return country
		}
		return province

	case config.IPRegionCity:
		if province == "0" && city == "0" {
			return country
		}
		if city == province { // e.g. 重庆重庆
			return province
		}
		return strings.TrimSuffix(province+""+city, "0")

	default:
		return ""
	}
}
