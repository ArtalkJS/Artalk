package core

import "github.com/artalkjs/artalk/v2/internal/ip_region"

var _ Service = (*IPRegionService)(nil)

type IPRegionService struct {
	app *App
}

func NewIPRegionService(app *App) *IPRegionService {
	return &IPRegionService{app: app}
}

func (s *IPRegionService) Init() error {
	return nil
}

func (s *IPRegionService) Dispose() error {
	return nil
}

func (s *IPRegionService) Query(ip string) string {
	client := ip_region.NewIPRegion(ip_region.IPRegionConf{
		IPRegionConf: s.app.Conf().IPRegion,
		CacheEnabled: s.app.Conf().Cache.Enabled,
	})

	return client.IP2Region(ip)
}
