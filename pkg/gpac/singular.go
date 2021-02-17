package gpac

import (
	customerrors "github.com/weiyuan-lane/gpac/pkg/errors"
)

type SimpleRetrieveFunc func(key string) (interface{}, error)
type RetrieveFunc func(subKeys ...ArgReference) (interface{}, error)

func (p *pageAwareCache) SimpleItem(subject interface{}, retrieveWith SimpleRetrieveFunc, key string) error {
	itemCacheKey := p.createItemCacheKeyFromStrKey(key)

	cachePayload, err := p.cacheClient.Get(itemCacheKey)
	if err != nil {
		return err
	}

	if cachePayload != nil {
		if err := p.decodeStringIntoInterfacePtr(subject, *cachePayload); err != nil {
			return err
		}

		return nil
	}

	payload, err := retrieveWith(key)
	if err != nil {
		return err
	}
	if p.isNil(payload) {
		return customerrors.ErrItemNotFound
	}

	err = p.copyBetweenPointers(payload, subject)
	if err != nil {
		return err
	}

	cacheVal, err := p.encodeInterfacePtrIntoString(payload)
	if err != nil {
		return err
	}

	err = p.cacheClient.Set(itemCacheKey, cacheVal, p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}

func (p *pageAwareCache) Item(subject interface{}, retrieveWith RetrieveFunc, subKeys ...ArgReference) error {
	itemCacheKey := p.createItemCacheKeyFromSubKeys(subKeys...)

	cachePayload, err := p.cacheClient.Get(itemCacheKey)
	if err != nil {
		return err
	}

	if cachePayload != nil {
		if err := p.decodeStringIntoInterfacePtr(subject, *cachePayload); err != nil {
			return err
		}

		return nil
	}

	payload, err := retrieveWith(subKeys...)
	if err != nil {
		return err
	}
	if p.isNil(payload) {
		return customerrors.ErrItemNotFound
	}

	err = p.copyBetweenPointers(payload, subject)
	if err != nil {
		return err
	}

	cacheVal, err := p.encodeInterfacePtrIntoString(payload)
	if err != nil {
		return err
	}

	err = p.cacheClient.Set(itemCacheKey, cacheVal, p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}
