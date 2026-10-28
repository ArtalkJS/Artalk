package ip_region

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// A one-record XDB fixture, generated locally so tests never download a database.
func xdbFixture(ip, raw string, legacy bool) []byte {
	address := net.ParseIP(ip).To4()
	version := xdb.IPv4
	if address == nil {
		address = net.ParseIP(ip).To16()
		version = xdb.IPv6
	}
	index := xdb.HeaderInfoLength + xdb.VectorIndexRows*xdb.VectorIndexCols*xdb.VectorIndexSize
	data := index + version.SegmentIndexSize
	body := make([]byte, data+len(raw))
	structure := uint16(xdb.Structure30)
	if legacy {
		structure = xdb.Structure20
	}
	binary.LittleEndian.PutUint16(body, structure)
	binary.LittleEndian.PutUint16(body[2:], uint16(xdb.VectorIndexPolicy))
	binary.LittleEndian.PutUint32(body[8:], uint32(index))
	binary.LittleEndian.PutUint32(body[12:], uint32(index))
	binary.LittleEndian.PutUint16(body[16:], uint16(version.Id))
	binary.LittleEndian.PutUint16(body[18:], 4)
	vector := xdb.HeaderInfoLength + (int(address[0])*256+int(address[1]))*8
	binary.LittleEndian.PutUint32(body[vector:], uint32(index))
	binary.LittleEndian.PutUint32(body[vector+4:], uint32(index))
	encoded := append([]byte(nil), address...)
	if version == xdb.IPv4 {
		encoded[0], encoded[1], encoded[2], encoded[3] = address[3], address[2], address[1], address[0]
	}
	copy(body[index:], encoded)
	copy(body[index+len(encoded):], encoded)
	binary.LittleEndian.PutUint16(body[index+2*len(encoded):], uint16(len(raw)))
	binary.LittleEndian.PutUint32(body[index+2*len(encoded)+2:], uint32(data))
	copy(body[data:], raw)
	return body
}

// A single leaf covers both address families. Field order intentionally differs
// from qqwry's published order to verify lookup by metadata field name.
func ipdbFixture(t *testing.T, fields, values []string, languages map[string]int) []byte {
	t.Helper()
	record := strings.Join(values, "\t")
	data := make([]byte, 11+len(record))
	binary.BigEndian.PutUint32(data, 2)
	binary.BigEndian.PutUint32(data[4:], 2)
	binary.BigEndian.PutUint16(data[9:], uint16(len(record)))
	copy(data[11:], record)
	meta, err := json.Marshal(map[string]any{
		"ip_version": 3, "node_count": 1, "total_size": len(data), "fields": fields, "languages": languages,
	})
	if err != nil {
		t.Fatal(err)
	}
	body := make([]byte, 4)
	binary.BigEndian.PutUint32(body, uint32(len(meta)))
	return append(append(body, meta...), data...)
}

func standardIPDB(t *testing.T) []byte {
	return ipdbFixture(t, []string{"city_name", "country_name", "isp_domain", "region_name", "country_code"},
		[]string{"福州市", "中国", "中国电信", "福建省", "CN"}, map[string]int{"CN": 0})
}

func writeDatabase(t *testing.T, path string, body []byte) {
	t.Helper()
	if err := os.WriteFile(path, body, 0600); err != nil {
		t.Fatal(err)
	}
}

func testClient(path, precision string, cached bool) *IPRegion {
	return NewIPRegion(IPRegionConf{IPRegionConf: config.IPRegionConf{DBPath: path, Precision: precision}, CacheEnabled: cached})
}

func TestDatabaseFormats(t *testing.T) {
	for _, tc := range []struct {
		name string
		body []byte
		ips  []string
	}{
		{"legacy-without-extension", xdbFixture("1.0.1.1", "中国|0|福建省|福州市|电信", true), []string{"1.0.1.1"}},
		{"current.xdb", xdbFixture("1.0.1.1", "中国|福建省|福州市|中国电信|CN", false), []string{"1.0.1.1"}},
		{"ipv6.xdb", xdbFixture("240e::1", "中国|福建省|福州市|中国电信|CN", false), []string{"240e::1"}},
		{"standard.IPDB", standardIPDB(t), []string{"1.0.1.1", "240e::1", "::ffff:1.0.1.1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tc.name)
			writeDatabase(t, path, tc.body)
			for _, cached := range []bool{false, true} {
				for precision, want := range map[string]string{"": "福建", "country": "中国", "province": "福建", "city": "福建福州", "invalid": ""} {
					t.Run(fmt.Sprintf("cache=%t/%s", cached, precision), func(t *testing.T) {
						client := testClient(path, precision, cached)
						for _, ip := range tc.ips {
							if got := client.IP2Region(" " + ip + ", 127.0.0.1 "); got != want {
								t.Errorf("%s: got %q want %q", ip, got, want)
							}
						}
					})
				}
			}
		})
	}
}

func TestIPDBPrecisionFallbacks(t *testing.T) {
	for _, tc := range []struct{ name, province, city, wantProvince, wantCity string }{
		{"municipality", "重庆市", "重庆市", "重庆", "重庆"},
		{"missing province", "", "福州市", "中国", "福州"},
		{"zero province", "0", "福州市", "中国", "福州"},
		{"missing city", "福建省", "", "福建", "福建"},
		{"missing both", "", "", "中国", "中国"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.ipdb")
			writeDatabase(t, path, ipdbFixture(t, []string{"country_name", "region_name", "city_name"}, []string{"中国", tc.province, tc.city}, map[string]int{"CN": 0}))
			for precision, want := range map[string]string{"country": "中国", "province": tc.wantProvince, "city": tc.wantCity} {
				if got := testClient(path, precision, true).IP2Region("240e::1"); got != want {
					t.Errorf("%s: got %q want %q", precision, got, want)
				}
			}
		})
	}
}

func TestIPDBLanguagesAndFields(t *testing.T) {
	for _, tc := range []struct {
		name           string
		fields, values []string
		languages      map[string]int
		want           string
		wantError      bool
	}{
		{"english", []string{"country_name", "region_name", "city_name"}, []string{"United States", "California", "Los Angeles"}, map[string]int{"EN": 0}, "California", false},
		{"prefer Chinese", []string{"country_name", "region_name", "city_name"}, []string{"China", "Fujian", "Fuzhou", "中国", "福建省", "福州市"}, map[string]int{"EN": 0, "CN": 3}, "福建", false},
		{"unsupported language", []string{"country_name", "region_name", "city_name"}, []string{"中国", "福建省", "福州市"}, map[string]int{"XX": 0}, "", true},
		{"missing province field", []string{"country_name", "city_name"}, []string{"中国", "福州市"}, map[string]int{"CN": 0}, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.ipdb")
			writeDatabase(t, path, ipdbFixture(t, tc.fields, tc.values, tc.languages))
			client := testClient(path, "province", true)
			_, err := client.search("1.0.1.1")
			if (err != nil) != tc.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := client.IP2Region("1.0.1.1"); got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestDatabaseCacheIsolationAndConcurrency(t *testing.T) {
	xdbPath := filepath.Join(t.TempDir(), "test.xdb")
	writeDatabase(t, xdbPath, xdbFixture("1.0.1.1", "中国|浙江省|杭州市|电信|CN", false))
	ipdbPath := filepath.Join(t.TempDir(), "test.ipdb")
	writeDatabase(t, ipdbPath, standardIPDB(t))
	for _, cached := range []bool{false, true} {
		clients := []*IPRegion{testClient(xdbPath, "province", cached), testClient(ipdbPath, "province", cached)}
		var wg sync.WaitGroup
		for range 16 {
			wg.Go(func() {
				for range 10 {
					for i, want := range []string{"浙江", "福建"} {
						if got := clients[i].IP2Region("1.0.1.1"); got != want {
							t.Errorf("cache=%t: got %q want %q", cached, got, want)
						}
					}
				}
			})
		}
		wg.Wait()
	}
}

func TestDatabaseCacheAndLoadRetry(t *testing.T) {
	for _, extension := range []string{".xdb", ".ipdb"} {
		for _, cached := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cache=%t", extension, cached), func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "test"+extension)
				client := testClient(path, "province", cached)
				if _, err := client.search("1.0.1.1"); err == nil {
					t.Fatal("missing file should fail")
				}
				body := standardIPDB(t)
				if extension == ".xdb" {
					body = xdbFixture("1.0.1.1", "中国|福建省|福州市|电信|CN", false)
				}
				writeDatabase(t, path, body)
				if got := client.IP2Region("1.0.1.1"); got != "福建" {
					t.Fatalf("load after repair: %q", got)
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				want := ""
				if cached {
					want = "福建"
				}
				if got := client.IP2Region("1.0.1.1"); got != want {
					t.Errorf("cache behavior: got %q want %q", got, want)
				}
			})
		}
	}
}

func TestInvalidIPAndAddressFamily(t *testing.T) {
	for _, extension := range []string{".xdb", ".ipdb"} {
		client := testClient(filepath.Join(t.TempDir(), "missing"+extension), "province", true)
		for _, ip := range []string{"", " ", "not-an-ip", "1.2.3.999", ", 1.0.1.1"} {
			if got := client.IP2Region(ip); got != "" {
				t.Errorf("invalid IP returned %q", got)
			}
		}
	}
	for _, ip := range []string{"1.0.1.1", "240e::1"} {
		path := filepath.Join(t.TempDir(), "test.xdb")
		writeDatabase(t, path, xdbFixture(ip, "中国|福建省|福州市|电信|CN", false))
		other := "240e::1"
		if ip == other {
			other = "1.0.1.1"
		}
		for _, cached := range []bool{false, true} {
			if got := testClient(path, "province", cached).IP2Region(other); got != "" {
				t.Errorf("wrong address family returned %q", got)
			}
		}
	}
}

func TestCorruptDatabaseDoesNotPanic(t *testing.T) {
	corruptIPDB := standardIPDB(t)
	offset := 4 + int(binary.BigEndian.Uint32(corruptIPDB))
	// Point the leaf at the last byte: the upstream decoder would read past it.
	pointer := uint32(len(corruptIPDB) - offset - 8)
	binary.BigEndian.PutUint32(corruptIPDB[offset:], pointer)
	binary.BigEndian.PutUint32(corruptIPDB[offset+4:], pointer)
	corruptXDB := xdbFixture("1.0.1.1", "中国|福建省|福州市|电信|CN", false)
	index := int(binary.LittleEndian.Uint32(corruptXDB[8:]))
	binary.LittleEndian.PutUint32(corruptXDB[index+10:], uint32(len(corruptXDB)+100))
	for _, tc := range []struct {
		name string
		body []byte
	}{
		{"truncated.xdb", make([]byte, 17)},
		{"truncated.ipdb", []byte{0, 0, 0, 100}},
		{"offset.xdb", corruptXDB},
		{"offset.ipdb", corruptIPDB},
	} {
		for _, cached := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cache=%t", tc.name, cached), func(t *testing.T) {
				path := filepath.Join(t.TempDir(), tc.name)
				writeDatabase(t, path, tc.body)
				client := testClient(path, "province", cached)
				for range 2 {
					if _, err := client.search("1.0.1.1"); err == nil {
						t.Fatal("expected database error")
					}
					if got := client.IP2Region("1.0.1.1"); got != "" {
						t.Errorf("got %q from corrupt database", got)
					}
				}
			})
		}
	}
}
