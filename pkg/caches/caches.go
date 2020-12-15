package caches

type CacheClient interface {
	Get(key string) (*string, error)
	Set(key, val string, ttl int) error
	Delete(key string) error
}
