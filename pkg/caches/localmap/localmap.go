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

func (l *LocalMap) Set(key, val string, ttl int) error {
	currTime := currTimeSeconds()

	l.storage[key] = localMapEntry{
		data:        val,
		expiredTime: currTime + int64(ttl),
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
