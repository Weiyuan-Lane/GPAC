package core

import (
	"bytes"
	"encoding/gob"
)

type RetrievalFunc func(key string) (interface{}, error)

func (p *PageAwareCache) Item(subject interface{}, key string, retrieveWith RetrievalFunc) error {
	itemCacheKey := p.createItemFullCacheKey(key)

	cachePayload, err := p.cacheClient.Get(itemCacheKey)
	if err != nil {
		// return err
	}

	if cachePayload != nil {
		buf := bytes.NewBufferString(*cachePayload)
		if err := gob.NewDecoder(buf).Decode(subject); err != nil {
			// return err
		}

		return nil
	}

	payload, err := retrieveWith(itemCacheKey)
	if err != nil {
		return err
	}

	err = p.copyBetweenPointers(payload, subject)
	if err != nil {
		return err
	}

	if !p.isNil(payload) {
		buf := &bytes.Buffer{}
		err := gob.NewEncoder(buf).Encode(payload)
		if err != nil {
			return err
		}

		p.cacheClient.Set(itemCacheKey, buf.String(), p.defaultItemTTL)
	}

	return nil
}
