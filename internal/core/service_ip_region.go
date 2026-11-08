package core

import (
	"sync"

	"github.com/artalkjs/artalk/v2/internal/ip_region"
)

var _ Service = (*IPRegionService)(nil)

type IPRegionService struct {
	app    *App
	mu     sync.Mutex
	conf   ip_region.IPRegionConf
	client *ip_region.IPRegion
}

func NewIPRegionService(app *App) *IPRegionService {
	return &IPRegionService{app: app}
}

func (s *IPRegionService) Init() error {
	return s.Dispose()
}

func (s *IPRegionService) Dispose() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client = nil // Cached databases are memory-only; release them for GC.
	return nil
}

func (s *IPRegionService) Query(ip string) string {
	return s.getClient().IP2Region(ip)
}

func (s *IPRegionService) getClient() *ip_region.IPRegion {
	s.mu.Lock()
	defer s.mu.Unlock()
	conf := ip_region.IPRegionConf{
		IPRegionConf: s.app.Conf().IPRegion,
		CacheEnabled: s.app.Conf().Cache.Enabled,
	}
	if s.client == nil || s.conf != conf {
		s.conf = conf
		s.client = ip_region.NewIPRegion(conf)
	}
	return s.client
}
