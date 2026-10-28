package ip_region

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

type regionSearcher interface {
	Search(ip string) (region, error)
	Close()
}

func newSearcher(dbPath string, useCache bool) (regionSearcher, error) {
	if strings.EqualFold(filepath.Ext(dbPath), ".ipdb") {
		return newIPDBSearcher(dbPath)
	}
	// Preserve existing XDB paths, including files without an extension.
	return newXDBSearcher(dbPath, useCache)
}

type xdbSearcher struct {
	mu sync.Mutex // xdb.Searcher mutates its IO counter even in memory mode.
	db *xdb.Searcher
}

func newXDBSearcher(dbPath string, useCache bool) (*xdbSearcher, error) {
	var (
		content []byte
		header  *xdb.Header
		err     error
	)
	if useCache {
		content, err = xdb.LoadContentFromFile(dbPath)
		if err != nil {
			return nil, err
		}
		header, err = xdb.LoadHeaderFromBuff(content)
	} else {
		header, err = xdb.LoadHeaderFromFile(dbPath)
	}
	if err != nil {
		return nil, err
	}
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		return nil, err
	}
	var db *xdb.Searcher
	if useCache {
		db, err = xdb.NewWithBuffer(version, content)
	} else {
		db, err = xdb.NewWithFileOnly(version, dbPath)
	}
	if err != nil {
		return nil, err
	}
	return &xdbSearcher{db: db}, nil
}

func (s *xdbSearcher) Search(ip string) (region, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, err := s.db.Search(ip)
	return parseXDBRegion(raw), err
}

func (s *xdbSearcher) Close() { s.db.Close() }
