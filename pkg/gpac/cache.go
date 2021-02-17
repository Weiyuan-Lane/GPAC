package gpac

import (
	"github.com/weiyuan-lane/gpac/pkg/caches"
	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/constants"
)

type PageAwareCache interface {
	SimpleItem(subject interface{}, retrieveWith SimpleRetrieveFunc, key string) error
	Item(subject interface{}, retrieveWith RetrieveFunc, subKeys ...ArgReference) error
	Items(subjectMap interface{}, retrieveWith SimpleRetrieveMultipleFunc, keyList []string) error
	Page(pageSubject interface{}, retrieveWith PageRetrievalFunc, retrieveItemsFrom PageToItemsFunc, retrieveKeyFrom ItemToKeyFunc, subKeys ...ArgReference) error
}

type pageAwareCache struct {
	cacheClient     caches.CacheClient
	uniqueNamespace string
	defaultItemTTL  int
	defaultPageTTL  int
}

type GPACOption func(*pageAwareCache)

func WithCacheClient(cacheClient caches.CacheClient) GPACOption {
	return func(p *pageAwareCache) {
		p.cacheClient = cacheClient
	}
}

func WithDefaultItemTTL(defaultItemTTL int) GPACOption {
	return func(p *pageAwareCache) {
		p.defaultItemTTL = defaultItemTTL
	}
}

func WithDefaultPageTTL(defaultPageTTL int) GPACOption {
	return func(p *pageAwareCache) {
		p.defaultPageTTL = defaultPageTTL
	}
}

func WithUniqueNamespace(uniqueNamespace string) GPACOption {
	return func(p *pageAwareCache) {
		p.uniqueNamespace = uniqueNamespace
	}
}

func NewGPAC(options ...GPACOption) PageAwareCache {
	cache := &pageAwareCache{
		cacheClient:    localmap.New(),
		defaultItemTTL: constants.DefaultItemTTL,
		defaultPageTTL: constants.DefaultPageTTL,
	}

	for _, option := range options {
		option(cache)
	}

	return cache
}
