package core

// Return should either be map[interface]{} or []interface{}
type MultipleRetrievalFunc func(keyList []string) (interface{}, error)

func (p *PageAwareCache) Items(subjectList interface{}, keyList []string, retrieveWith MultipleRetrievalFunc) error {
	itemCacheKeyList := make([]string, len(keyList))
	for i, key := range keyList {
		itemCacheKeyList[i] = p.createItemFullCacheKey(key)
	}

	for _, cacheKey := range itemCacheKeyList {

	}

	cachePayload, err := p.cacheClient.Get(cacheKey)
	if err != nil {
		return err
	}
}
