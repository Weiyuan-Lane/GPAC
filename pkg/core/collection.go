package core

type MultipleRetrievalFunc func(keyList []string) (map[string]interface{}, error)

func (p *PageAwareCache) Items(subjectMap interface{}, keyList []string, retrieveWith MultipleRetrievalFunc) error {
	itemCacheKeyList := make([]string, len(keyList))
	for i, key := range keyList {
		itemCacheKeyList[i] = p.createItemFullCacheKey(key)
	}

	cachePayloads, err := p.cacheClient.MultipleGet(itemCacheKeyList)
	if err != nil {
		return err
	}

	// For each cache item, find those that are not cached

	// Cache the items that were retrieved
	// return the rest
}
