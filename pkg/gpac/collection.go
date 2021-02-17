package gpac

type SimpleRetrieveMultipleFunc func(keyList []string) (interface{}, error)

func (p *pageAwareCache) Items(subjectMap interface{}, retrieveWith SimpleRetrieveMultipleFunc, keyList []string) error {
	// Create a list of keys that we want to reference from cache
	cacheKeyList := make([]string, len(keyList))
	keyIndexMap := map[string]int{}
	for i, key := range keyList {
		keyIndexMap[key] = i
		cacheKeyList[i] = p.createItemCacheKeyFromStrKey(key)
	}

	// Hit cache to get results
	cachePayloadMap, err := p.cacheClient.MultipleGet(cacheKeyList)
	if err != nil {
		return err
	}

	// Divide the results into those found and those not found
	cacheOriginalKeyMap := map[int]string{}
	missingKeys := make([]string, 0, len(cacheKeyList))
	for i, cacheKey := range cacheKeyList {
		originalKey := keyList[i]

		if val, ok := cachePayloadMap[cacheKey]; ok {
			index := keyIndexMap[originalKey]
			cacheOriginalKeyMap[index] = val
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

	payloadPtr := p.makePointerTo(payload)

	// Move results into subject map
	err = p.copyBetweenPointerMaps(payloadPtr, subjectMap)
	if err != nil {
		return err
	}

	// Convert retrieved payload to a cacheable payload
	encodedMap, err := p.encodeMapPtrIntoMap(payloadPtr)
	if err != nil {
		return err
	}
	cacheInputMap := map[string]string{}
	for i, val := range encodedMap {
		key := cacheKeyList[i]
		cacheInputMap[key] = val
	}

	// Set all items to cache
	err = p.cacheClient.MultipleSet(cacheInputMap, p.defaultItemTTL)
	if err != nil {
		return err
	}

	return nil
}
