package gpacmock

import (
	"reflect"

	customerrors "github.com/weiyuan-lane/gpac/pkg/errors"
	"github.com/weiyuan-lane/gpac/pkg/gpac"
)

type GPACMock struct{}

func NewGPACMock() gpac.PageAwareCache {
	return &GPACMock{}
}

func (g *GPACMock) SimpleItem(subject interface{}, retrieveWith gpac.SimpleRetrieveFunc, key string) error {
	val, err := retrieveWith(key)
	if err != nil {
		return err
	}

	if reflect.ValueOf(val).IsNil() {
		return customerrors.ErrItemNotFound
	}

	g.copyBetweenPointers(val, subject)
	return nil
}

func (g *GPACMock) Item(subject interface{}, retrieveWith gpac.RetrieveFunc, subKeys ...gpac.ArgReference) error {
	val, err := retrieveWith(subKeys...)
	if err != nil {
		return err
	}

	if reflect.ValueOf(val).IsNil() {
		return customerrors.ErrItemNotFound
	}

	g.copyBetweenPointers(val, subject)
	return nil
}

func (g *GPACMock) Items(subjectMap interface{}, retrieveWith gpac.SimpleRetrieveMultipleFunc, keyList []string) error {
	payload, err := retrieveWith(keyList)
	if err != nil {
		return err
	}

	payloadPtr := g.makePointerTo(payload)

	// Move results into subject map
	g.copyBetweenPointerMaps(payloadPtr, subjectMap)
	return nil
}

func (g *GPACMock) Page(pageSubject interface{}, retrieveWith gpac.PageRetrievalFunc, retrieveItemsFrom gpac.PageToItemsFunc, retrieveKeyFrom gpac.ItemToKeyFunc, subKeys ...gpac.ArgReference) error {
	val, err := retrieveWith(subKeys...)
	if err != nil {
		return err
	}

	if reflect.ValueOf(val).IsNil() {
		return customerrors.ErrItemNotFound
	}

	g.copyBetweenPointers(val, pageSubject)
	return nil
}

func (g *GPACMock) copyBetweenPointers(src, dest interface{}) {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)
	srcElem := srcVal.Elem()
	destElem := destVal.Elem()
	destElem.Set(srcElem)
}

func (g *GPACMock) copyBetweenPointerMaps(srcMapPtr, destMapPtr interface{}) error {
	srcMapVal := reflect.ValueOf(srcMapPtr)
	destMapVal := reflect.ValueOf(destMapPtr)
	srcMapElem := srcMapVal.Elem()
	destMapElem := destMapVal.Elem()

	iter := srcMapElem.MapRange()
	for iter.Next() {
		destMapElem.SetMapIndex(
			iter.Key(),
			iter.Value(),
		)
	}

	return nil
}

func (p *GPACMock) makePointerTo(subject interface{}) interface{} {
	subjectPtr := reflect.New(reflect.TypeOf(subject))
	reflect.Indirect(subjectPtr).Set(reflect.ValueOf(subject))

	return subjectPtr.Interface()
}
