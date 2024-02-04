package ip_region

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/log"
)

type IPRegion struct {
	conf IPRegionConf
}

type IPRegionConf struct {
	config.IPRegionConf
	CacheEnabled bool
}

func NewIPRegion(conf IPRegionConf) *IPRegion {
	return &IPRegion{
		conf: conf,
	}
}

func (ipRegion *IPRegion) IP2Region(ip string) string {
	if strings.TrimSpace(ip) == "" {
		return ""
	}

	ip = ipScraper(ip)
	region, err := search(ip, ipRegion.conf.DBPath, ipRegion.conf.CacheEnabled)
	if err != nil {
		log.Warn("[IP2Region] ", err)
		return ""
	}

	return regionScraper(region, config.IPRegionPrecision(ipRegion.conf.Precision))
}
