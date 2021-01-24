package core

import (
	"github.com/weiyuan-lane/gpac/pkg/caches"
	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/constants"
)

type KeyFromItemFunc func(item interface{}) string

type PageAwareCache struct {
	cacheClient     caches.CacheClient
	uniqueNamespace string
	defaultItemTTL  int
	defaultPageTTL  int
}

type GPACOption func(*PageAwareCache)

func WithCacheClient(cacheClient caches.CacheClient) GPACOption {
	return func(p *PageAwareCache) {
		p.cacheClient = cacheClient
	}
}

func WithDefaultItemTTL(defaultItemTTL int) GPACOption {
	return func(p *PageAwareCache) {
		p.defaultItemTTL = defaultItemTTL
	}
}

func WithDefaultPageTTL(defaultPageTTL int) GPACOption {
	return func(p *PageAwareCache) {
		p.defaultPageTTL = defaultPageTTL
	}
}

func WithUniqueNamespace(uniqueNamespace string) GPACOption {
	return func(p *PageAwareCache) {
		p.uniqueNamespace = uniqueNamespace
	}
}

func NewGPAC(options ...GPACOption) *PageAwareCache {
	cache := &PageAwareCache{
		cacheClient:    localmap.New(),
		defaultItemTTL: constants.DefaultItemTTL,
		defaultPageTTL: constants.DefaultPageTTL,
	}

	for _, option := range options {
		option(cache)
	}

	return cache
}
