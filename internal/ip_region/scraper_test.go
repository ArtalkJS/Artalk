package ip_region

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
)

func TestIpScraper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  127.0.0.1  ", "127.0.0.1"},
		{"192.168.1.1, 10.0.0.1", "192.168.1.1"},
		{"  2001:0db8:85a3:0000:0000:8a2e:0370:7334  ", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		{"", ""},
		{", 172.16.0.1", ""},
	}

	for _, test := range tests {
		result := ipScraper(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s' but got '%s'", test.input, test.expected, result)
		}
	}
}

func TestRegionScraper(t *testing.T) {
	tests := []struct {
		raw       string
		precision config.IPRegionPrecision
		expected  string
	}{
		{"0|0|0|0|0", config.IPRegionCountry, ""},
		{"中国|0|上海市|上海市|电信", config.IPRegionProvince, "上海"},
		{"中国|0|上海市|上海市|电信", config.IPRegionCity, "上海"},
		{"中国|0|重庆市|重庆市|电信", config.IPRegionCity, "重庆"},
		{"中国|0|浙江省|杭州市|电信", config.IPRegionProvince, "浙江"},
		{"中国|0|北京市|北京市|电信", config.IPRegionCity, "北京"},
		{"中国|0|四川省|泸州市|电信", config.IPRegionCity, "四川泸州"},
		{"中国|0|广东省|深圳市|电信", config.IPRegionCity, "广东深圳"},
		{"中国|0|湖南省|岳阳市|电信", config.IPRegionCity, "湖南岳阳"},
		{"中国|0|陕西省|西安市|电信", config.IPRegionCity, "陕西西安"},
		{"中国|0|河北省|石家庄市|电信", config.IPRegionCity, "河北石家庄"},
		{"中国|0|江苏省|南京市|电信", config.IPRegionCity, "江苏南京"},
		{"美国|0|0|0|微软云", config.IPRegionCity, "美国"},
		// Legacy databases keep an unused area column before the province.
		{"中国|0|福建省|福州市|电信", "", "福建"},
		{"中国|0|福建省|福州市|电信", config.IPRegionCountry, "中国"},
		{"中国|0|福建省|福州市|电信", config.IPRegionProvince, "福建"},
		{"中国|0|福建省|福州市|电信", config.IPRegionCity, "福建福州"},
		{"印度尼西亚|0|大雅加达|雅加达|XL", config.IPRegionProvince, "大雅加达"},
		// Current databases omit area and append an ISO country code instead.
		{"中国|福建省|福州市|中国电信|CN", "", "福建"},
		{"中国|福建省|福州市|中国电信|CN", config.IPRegionCountry, "中国"},
		{"中国|福建省|福州市|中国电信|CN", config.IPRegionProvince, "福建"},
		{"中国|福建省|福州市|中国电信|CN", config.IPRegionCity, "福建福州"},
		{"中国|重庆市|重庆市|中国电信|CN", config.IPRegionCity, "重庆"},
		{"中国|福建省|0|中国电信|CN", config.IPRegionCity, "福建"},
		{"中国|0|0|阿里|CN", config.IPRegionProvince, "中国"},
		{"中国|0|0|阿里|CN", config.IPRegionCity, "中国"},
		{"United States|California|Los Angeles|0|US", config.IPRegionProvince, "California"},
		{"Seychelles|0|0|0|SC", config.IPRegionCity, "Seychelles"},
		{"Kosovo|Pristina|Pristina|0|XK", config.IPRegionProvince, "Pristina"},
		// Missing provinces must not expose a city at province precision.
		{"中国|0|福州市|中国电信|CN", config.IPRegionProvince, "中国"},
		{"中国|0|福州市|中国电信|CN", config.IPRegionCity, "福州"},
		{"中国|0|0|福州市|电信", config.IPRegionProvince, "中国"},
		{"中国|0|0|福州市|电信", config.IPRegionCity, "福州"},
		{"中国|||中国电信|CN", config.IPRegionProvince, "中国"},
		{"中国|||中国电信|CN", config.IPRegionCity, "中国"},
		{"中国|福建省|福州市|中国电信|CN|extra", config.IPRegionProvince, "福建"},
		{"", config.IPRegionProvince, ""},
		{"中国|福建省|福州市", config.IPRegionProvince, ""},
		{"||||", config.IPRegionProvince, ""},
		{"中国|福建省|福州市|中国电信|CN", "invalid", ""},
	}

	for _, test := range tests {
		result := parseXDBRegion(test.raw).format(test.precision)
		if result != test.expected {
			t.Errorf("For input '%s' and precision '%s', expected '%s' but got '%s'", test.raw, test.precision, test.expected, result)
		}
	}
}
