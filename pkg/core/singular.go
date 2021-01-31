package core

import (
	"bytes"
	"encoding/gob"

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
		buf := bytes.NewBufferString(*cachePayload)
		if err := gob.NewDecoder(buf).Decode(subject); err != nil {
			return err
		}

		return nil
	}

	payload, err := retrieveWith(itemCacheKey)
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

	buf := &bytes.Buffer{}
	err = gob.NewEncoder(buf).Encode(payload)
	if err != nil {
		return err
	}

	err = p.cacheClient.Set(itemCacheKey, buf.String(), p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}
