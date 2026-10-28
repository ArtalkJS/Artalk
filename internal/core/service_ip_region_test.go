package core

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
)

func TestIPRegionServiceConfigurationChanges(t *testing.T) {
	// Each tiny IPDB has a single leaf for every address. Using real readers here
	// exercises cache invalidation through the service, not only config equality.
	writeIPDB := func(path, province, city string) {
		t.Helper()
		record := "中国\t" + province + "\t" + city
		data := make([]byte, 11+len(record))
		binary.BigEndian.PutUint32(data, 2)
		binary.BigEndian.PutUint32(data[4:], 2)
		binary.BigEndian.PutUint16(data[9:], uint16(len(record)))
		copy(data[11:], record)
		meta, err := json.Marshal(map[string]any{"ip_version": 3, "node_count": 1, "total_size": len(data),
			"fields": []string{"country_name", "region_name", "city_name"}, "languages": map[string]int{"CN": 0}})
		if err != nil {
			t.Fatal(err)
		}
		body := make([]byte, 4)
		binary.BigEndian.PutUint32(body, uint32(len(meta)))
		if err := os.WriteFile(path, append(append(body, meta...), data...), 0600); err != nil {
			t.Fatal(err)
		}
	}
	first := filepath.Join(t.TempDir(), "first.ipdb")
	writeIPDB(first, "福建省", "福州市")
	second := filepath.Join(t.TempDir(), "second.ipdb")
	writeIPDB(second, "浙江省", "杭州市")
	conf := &config.Config{}
	conf.IPRegion = config.IPRegionConf{Enabled: true, DBPath: first, Precision: "province"}
	conf.Cache.Enabled = true
	app := NewApp(conf)
	service := NewIPRegionService(app)
	check := func(want string) {
		t.Helper()
		if got := service.Query("240e::1"); got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	}
	check("福建")
	updated := *app.Conf()
	updated.IPRegion.DBPath = second
	app.SetConf(&updated)
	check("浙江")
	updated = *app.Conf()
	updated.IPRegion.Precision = "city"
	app.SetConf(&updated)
	check("浙江杭州")
	updated = *app.Conf()
	updated.Cache.Enabled = false
	app.SetConf(&updated)
	writeIPDB(second, "广东省", "广州市")
	check("广东广州")
	writeIPDB(second, "四川省", "成都市")
	check("四川成都")
	updated = *app.Conf()
	updated.Cache.Enabled = true
	app.SetConf(&updated)
	check("四川成都")
	writeIPDB(second, "福建省", "福州市")
	check("四川成都") // Cached until the service is reinitialized.
	if err := service.Init(); err != nil {
		t.Fatal(err)
	}
	check("福建福州")

	// Disposing/reinitializing only drops memory-owned readers, so in-flight
	// requests remain safe while later requests acquire a new reader.
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 10 {
				if got := service.Query("1.0.1.1"); got != "福建福州" {
					t.Errorf("concurrent lookup: %q", got)
				}
			}
		})
	}
	for range 10 {
		if err := service.Dispose(); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
}
