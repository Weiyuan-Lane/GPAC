package localmap

import (
	"time"
)

type localMapEntry struct {
	data        string
	expiredTime int64
}

type LocalMap struct {
	storage map[string]localMapEntry
}

func New() *LocalMap {
	return &LocalMap{
		storage: map[string]localMapEntry{},
	}
}

func (l *LocalMap) Get(key string) (*string, error) {
	currTime := currTimeSeconds()

	if val, ok := l.storage[key]; ok {
		if currTime > val.expiredTime {
			l.Delete(key)
			return nil, nil
		}

		return &val.data, nil
	}

	return nil, nil
}

func (l *LocalMap) MultipleGet(keys []string) (map[string]string, error) {
	result := map[string]string{}
	for _, key := range keys {
		strPtr, err := l.Get(key)
		if err != nil {
			return map[string]string{}, err
		}

		if strPtr != nil {
			result[key] = *strPtr
		}
	}

	return result, nil
}

func (l *LocalMap) Set(key, val string, ttl int) error {
	currTime := currTimeSeconds()

	l.storage[key] = localMapEntry{
		data:        val,
		expiredTime: currTime + int64(ttl),
	}

	return nil
}

func (l *LocalMap) MultipleSet(valMap map[string]string, ttl int) error {
	for k, v := range valMap {
		err := l.Set(k, v, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalMap) Delete(key string) error {
	if _, ok := l.storage[key]; ok {
		delete(l.storage, key)

		return nil
	}

	return nil
}

func currTimeSeconds() int64 {
	return time.Now().UTC().Unix()
}
