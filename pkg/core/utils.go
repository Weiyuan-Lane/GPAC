package core

import (
	"fmt"
	"reflect"

	"github.com/weiyuan-lane/gpac/pkg/constants"
	customerrors "github.com/weiyuan-lane/gpac/pkg/errors"
)

// Using the unique namespace and the item key, create a key unique
// to this target resource
func (p *PageAwareCache) createItemFullCacheKey(itemKey string) string {
	return fmt.Sprintf(constants.ItemKeyTemplate, p.uniqueNamespace, itemKey)
}

// Using the unique namespace and the page key, create a key unique
// to this target resource page
func (p *PageAwareCache) createPageFullCacheKey(pageKey string) string {
	return fmt.Sprintf(constants.PageKeyTemplate, p.uniqueNamespace, pageKey)
}

func (p *PageAwareCache) isPointer(kind interface{}) bool {
	return reflect.ValueOf(kind).Kind() == reflect.Ptr
}

func (p *PageAwareCache) isNil(kind interface{}) bool {
	return reflect.ValueOf(kind).IsNil()
}

func (p *PageAwareCache) copyBetweenPointers(src, dest interface{}) error {
	if src == nil {
		return customerrors.SourceValIsNilErr
	}
	if dest == nil {
		return customerrors.DestinationValIsNilErr
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Ptr {
		return customerrors.SourceValIsNotPtrErr
	}
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return customerrors.DestinationValIsNotPtrErr
	}

	srcElem := srcVal.Elem()
	destElem := destVal.Elem()
	destElem.Set(srcElem)

	return nil
}

func (p *PageAwareCache) copyBetweenPointerLists(srcListPtr, destListPtr interface{}) error {
	if srcListPtr == nil {
		return customerrors.SourceListValIsNilErr
	}
	if destListPtr == nil {
		return customerrors.DestinationListValIsNilErr
	}

	srcListVal := reflect.ValueOf(srcListPtr)
	if srcListVal.Kind() != reflect.Ptr {
		return customerrors.SourceListValIsNotPtrErr
	}
	destListVal := reflect.ValueOf(destListPtr)
	if destListVal.Kind() != reflect.Ptr {
		return customerrors.DestinationListValIsNotPtrErr
	}

	srcListElemVal := srcListVal.Elem()
	if srcListElemVal.Kind() != reflect.Slice {
		return customerrors.SourceListValIsNotSliceErr
	}
	destListElemVal := destListVal.Elem()
	if destListElemVal.Kind() != reflect.Slice {
		return customerrors.DestinationListValIsNotSliceErr
	}

	targetLength := destListElemVal.Len()
	if srcListElemVal.Len() != targetLength {
		return customerrors.DifferentLengthOfUnitsErr
	}

	for i := 0; i < targetLength; i++ {
		srcIndexVal := srcListElemVal.Index(i)
		destIndexVal := destListElemVal.Index(i)

		destIndexVal.Set(srcIndexVal)
	}

	return nil
}
