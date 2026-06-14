package ip_region

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	inMemorySearcher *xdb.Searcher
	inMemoryLoadOnce sync.Once
	inMemoryLoadErr  error
)

func newSearcherWithCache(dbPath string) (*xdb.Searcher, error) {
	inMemoryLoadOnce.Do(func() {
		// 1、从 dbPath 加载整个 xdb 到内存
		cBuff, err := xdb.LoadContentFromFile(dbPath)
		if err != nil {
			inMemoryLoadErr = fmt.Errorf("failed to load content from %s: %w", strconv.Quote(dbPath), err)
			return
		}

		// 2、从内容缓冲区中读取 header 以确定 IP 版本
		header, err := xdb.LoadHeaderFromBuff(cBuff)
		if err != nil {
			inMemoryLoadErr = fmt.Errorf("failed to load header from content: %w", err)
			return
		}
		version, err := xdb.VersionFromHeader(header)
		if err != nil {
			inMemoryLoadErr = fmt.Errorf("failed to detect IP version from header: %w", err)
			return
		}

		// 3、用全局的 cBuff 创建完全基于内存的查询对象。
		searcher, err := xdb.NewWithBuffer(version, cBuff)
		if err != nil {
			inMemoryLoadErr = fmt.Errorf("failed to create searcher with content: %w", err)
			return
		}

		inMemorySearcher = searcher
	})

	return inMemorySearcher, inMemoryLoadErr
}

func newSearcherWithFileOnly(dbPath string) (*xdb.Searcher, error) {
	header, err := xdb.LoadHeaderFromFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load header from %s: %w", strconv.Quote(dbPath), err)
	}
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to detect IP version from %s: %w", strconv.Quote(dbPath), err)
	}
	searcher, err := xdb.NewWithFileOnly(version, dbPath)
	if err != nil {
		return searcher, fmt.Errorf("failed to load content from %s: %w", strconv.Quote(dbPath), err)
	}

	return searcher, nil
}

func getSearcher(dbPath string, useCache bool) (*xdb.Searcher, error) {
	// Using a cache would consume more memory space but accelerate IO.
	if useCache {
		return newSearcherWithCache(dbPath)
	} else {
		return newSearcherWithFileOnly(dbPath)
	}
}

func search(ip string, dbPath string, useCache bool) (string, error) {
	searcher, err := getSearcher(dbPath, useCache)
	if err != nil {
		return "", err
	}

	if !useCache { // 未使用缓存查询后释放以节省内存
		defer searcher.Close()
	}

	return searcher.Search(ip)
}
