package core

import (
	customerrors "github.com/weiyuan-lane/gpac/pkg/errors"
)

type RetrievalFunc func(key string) (interface{}, error)

func (p *PageAwareCache) Item(subject interface{}, key string, retrieveWith RetrievalFunc) error {
	itemCacheKey := p.createItemFullCacheKey(key)

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

	err = p.copyBetweenPointers(subject, payload)
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
