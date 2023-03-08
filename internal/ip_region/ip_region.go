package ip_region

import (
	"strconv"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

var searcher *xdb.Searcher
var precisionConf = Province

type Precision string

const (
	Province Precision = "province"
	City     Precision = "city"
	Country  Precision = "country"
)

func Init(dbPath string) {
	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		logrus.Fatal("failed to load content from "+strconv.Quote(dbPath)+", ", err)
		return
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err = xdb.NewWithBuffer(cBuff)
	if err != nil {
		logrus.Fatal("failed to create searcher with content: ", err)
		return
	}
}

func SetPrecision(p Precision) {
	precisionConf = p
}

func IP2Region(ip string) string {
	if searcher == nil || strings.TrimSpace(ip) == "" {
		return ""
	}
	if precisionConf == "" {
		precisionConf = Province
	}

	// 多 IP 仅选第一个
	ipSep := strings.Split(ip, ",")
	if len(ipSep) > 1 {
		ip = strings.TrimSpace(ipSep[0])
	}

	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return scraper(region, precisionConf)
}

func scraper(raw string, precision Precision) (r string) {
	sep := strings.Split(raw, "|")
	if len(sep) < 5 || sep[0] == "0" {
		return
	}

	var (
		country  = sep[0]
		province = strings.TrimSuffix(sep[2], "省")
		city     = strings.TrimSuffix(sep[3], "市")
		// isp      = sep[4]
	)

	if precision == Country || province == "0" {
		return country
	}

	switch precision {
	case Province:
		return province
	case City:
		if city == province { // e.g. 重庆重庆
			return province
		}
		return strings.TrimSuffix(province+""+city, "0")
	}
	return
}
