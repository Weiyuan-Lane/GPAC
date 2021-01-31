package core

type PageRetrievalFunc func(keyArgs ...interface{}) (interface{}, error)

// func (p *PageAwareCache) Page(subjectPage interface{}, retrieveWith PageRetrievalFunc, keyArgs ...interface{}) error {

// 	// Retrieve page from intended store func
// 	pagePayload, err := retrieveWith(keyArgs...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
