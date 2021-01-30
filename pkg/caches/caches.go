package caches

type CacheClient interface {
	Get(key string) (*string, error)
	MultipleGet(keys []string) (map[string]string, error)
	Set(key, val string, ttl int) error
	MultipleSet(valMap map[string]string, ttl int) error
	Delete(key string) error
}
