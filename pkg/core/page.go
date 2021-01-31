package core

type SimplePageRetrievalFunc func(keyArgs ...ArgReference) (interface{}, error)
type ItemToKeyFunc func(subject interface{}) (string, error)
type PageToItemsFunc func(page interface{}) ([]interface{}, error)

func (p *PageAwareCache) SimplePage(
	subjectPage interface{},
	retrieveWith SimplePageRetrievalFunc,
	retrieveKeyFrom ItemToKeyFunc,
	retrieveItemsFrom PageToItemsFunc,
	subKeys ...ArgReference,
) error {
	pageCacheKey := p.createPageCacheKeyFromSubKeys(subKeys...)

	// Get page from cache
	cachePagePayload, err := p.cacheClient.Get(pageCacheKey)
	if err != nil {
		return err
	}

	if cachePagePayload != nil {

		return nil
	}

	// Retrieve page from intended store func
	// pagePayload, err := retrieveWith(subKeys...)
	// if err != nil {
	// 	return err
	// }

	return nil
}
