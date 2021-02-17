package gpac

import (
	"reflect"
)

type PageRetrievalFunc func(keyArgs ...ArgReference) (interface{}, error)
type ItemToKeyFunc func(subject interface{}) (string, error)
type PageToItemsFunc func(pageSubject interface{}) ([]interface{}, error)

func (p *PageAwareCache) Page(
	pageSubject interface{},
	retrieveWith PageRetrievalFunc,
	retrieveItemsFrom PageToItemsFunc,
	retrieveKeyFrom ItemToKeyFunc,
	subKeys ...ArgReference,
) error {
	pageCacheKey := p.createPageCacheKeyFromSubKeys(subKeys...)

	// Get page item from cache
	cachePagePayload, err := p.cacheClient.Get(pageCacheKey)
	if err != nil {
		return err
	}

	if cachePagePayload != nil {
		if err := p.decodeStringIntoInterfacePtr(pageSubject, *cachePagePayload); err != nil {
			return err
		}

		return nil
	}

	// Retrieve page from intended store func
	pagePayload, err := retrieveWith(subKeys...)
	if err != nil {
		return err
	}

	err = p.copyBetweenPointers(&pagePayload, pageSubject)
	if err != nil {
		return err
	}

	pageCacheVal, err := p.encodeInterfacePtrIntoString(pagePayload)
	if err != nil {
		return err
	}

	err = p.cacheClient.Set(pageCacheKey, pageCacheVal, p.defaultPageTTL)
	if err != nil {
		return err
	}

	items, err := retrieveItemsFrom(pagePayload)
	if err != nil {
		return err
	}

	err = p.cachePageItems(items, retrieveItemsFrom, retrieveKeyFrom)
	if err != nil {
		return err
	}

	return nil
}

func (p *PageAwareCache) cachePageItems(
	pagePayload interface{},
	retrieveItemsFrom PageToItemsFunc,
	retrieveKeyFrom ItemToKeyFunc,
) error {

	items, err := retrieveItemsFrom(pagePayload)
	if err != nil {
		return err
	}

	cacheInputMap := map[string]string{}
	for i, item := range items {
		key, err := retrieveKeyFrom(item)
		if err != nil {
			return err
		}

		cacheKey := p.createItemCacheKeyFromStrKey(key)
		cacheItemPtr := reflect.ValueOf(items[i]).Addr().Interface()
		cacheVal, err := p.encodeInterfacePtrIntoString(cacheItemPtr)
		if err != nil {
			return err
		}

		cacheInputMap[cacheKey] = cacheVal
	}

	err = p.cacheClient.MultipleSet(cacheInputMap, p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}
