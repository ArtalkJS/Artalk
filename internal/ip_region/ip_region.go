package ip_region

import (
	"strconv"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

var searcher *xdb.Searcher

func Init(dbPath string) {
	if dbPath == "" {
		dbPath = "./data/ip2region.xdb"
	}

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

type Precision string

const (
	Province Precision = "Province"
	City     Precision = "City"
	Country  Precision = "Country"
)

func IP2Region(ip string, precision Precision) string {
	if searcher == nil {
		return ""
	}
	if precision == "" {
		precision = Province
	}

	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return scraper(region, precision)
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
		return strings.TrimSuffix(province+city, "0")
	}
	return
}
