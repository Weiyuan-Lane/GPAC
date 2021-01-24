package core

import (
	"fmt"
	"reflect"

	"github.com/weiyuan-lane/gpac/pkg/constants"
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

func (p *PageAwareCache) copyBetweenPointers(dest, src interface{}) error {
	if src == nil {
		// return error
	}
	if dest == nil {
		// return error
	}

	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)
	if srcVal.Kind() != reflect.Ptr {
		// return error
	}
	if destVal.Kind() != reflect.Ptr {
		// return error
	}

	srcElem := srcVal.Elem()
	destElem := destVal.Elem()
	srcElem.Set(destElem)

	return nil
}
