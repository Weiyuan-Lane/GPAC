package core

import (
	"bytes"
	"encoding/gob"
	"reflect"
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

func (p *PageAwareCache) isPointer(kind interface{}) bool {
	return reflect.ValueOf(kind).Kind() == reflect.Ptr
}

func (p *PageAwareCache) isNil(kind interface{}) bool {
	return reflect.ValueOf(kind).IsNil()
}

func (p *PageAwareCache) copyBetweenPointers(dest, src interface{}) error {
	if src == nil {
		// return error
	}
	if dest == nil {
		// return error
	}

	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)
	if srcVal.Kind() != reflect.Ptr {
		// return error
	}
	if destVal.Kind() != reflect.Ptr {
		// return error
	}

	srcElem := srcVal.Elem()
	destElem := destVal.Elem()
	srcElem.Set(destElem)

	return nil
}
