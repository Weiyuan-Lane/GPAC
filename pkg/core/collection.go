package core

type MultipleRetrievalFunc func(keyList []string) (map[string]interface{}, error)

func (p *PageAwareCache) Items(subjectMap interface{}, keyList []string, retrieveWith MultipleRetrievalFunc) error {
	// Create a list of keys that we want to reference from cache
	cacheKeyList := make([]string, len(keyList))
	for i, key := range keyList {
		cacheKeyList[i] = p.createItemFullCacheKey(key)
	}

	// Hit cache to get results
	cachePayloadMap, err := p.cacheClient.MultipleGet(cacheKeyList)
	if err != nil {
		return err
	}

	// Divide the results into those found and those not found
	cacheOriginalKeyMap := map[string]string{}
	missingKeys := make([]string, 0, len(cacheKeyList))
	for i, cacheKey := range cacheKeyList {
		originalKey := keyList[i]

		if val, ok := cachePayloadMap[cacheKey]; ok {
			cacheOriginalKeyMap[originalKey] = val
			continue
		}

		missingKeys = append(missingKeys, originalKey)
	}

	// If at least one item was found in cache
	if len(cacheOriginalKeyMap) > 0 {
		// Decode from cache values into subject map
		err = p.decodeMapIntoMapPtr(subjectMap, cacheOriginalKeyMap)
		if err != nil {
			return err
		}
	}

	// All items are found in cache, can terminate early
	if len(missingKeys) == 0 {
		return nil
	}

	// Get remaining results from data store
	payload, err := retrieveWith(missingKeys)
	if err != nil {
		return err
	}

	// Move results into subject map
	err = p.copyBetweenPointerMaps(subjectMap, &payload)
	if err != nil {
		return err
	}

	// Convert retrieved payload to a cacheable payload
	encodedMap, err := p.encodeMapPtrIntoMap(&payload)
	if err != nil {
		return err
	}
	cacheInputMap := map[string]string{}
	for key, val := range encodedMap {
		cacheKey := p.createItemFullCacheKey(key)
		cacheInputMap[cacheKey] = val
	}

	// Set all items to cache
	err = p.cacheClient.MultipleSet(cacheInputMap, p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}
