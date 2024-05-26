package cache_test

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/stretchr/testify/assert"
)

func newTestCache(t *testing.T) *cache.Cache {
	cache, err := cache.New(config.CacheConf{
		Enabled: true,
		Type:    config.CacheTypeBuiltin,
	})
	if err != nil {
		t.Fatal(err)
	}
	return cache
}

func TestNew(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	assert.NotNil(t, cache)
}
