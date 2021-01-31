package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"

	"github.com/weiyuan-lane/gpac/pkg/constants"
	customerrors "github.com/weiyuan-lane/gpac/pkg/errors"
)

// Using the unique namespace and the item key, create a key unique
// to this target resource
func (p *PageAwareCache) createItemCacheKeyFromStrKey(itemKey string) string {
	return fmt.Sprintf(constants.ItemKeyTemplate, p.uniqueNamespace, itemKey)
}

func (p *PageAwareCache) createItemCacheKeyFromSubKeys(subKeys ...ArgReference) string {
	cacheSubKeys := make([]string, len(subKeys))
	for i, subKey := range subKeys {
		cacheSubKeys[i] = stringifyArgReference(subKey, constants.ArgDivider)
	}

	cacheKey := strings.Join(cacheSubKeys, constants.SubkeyDivider)

	return cacheKey
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

func (p *PageAwareCache) decodeStringIntoInterfacePtr(subject interface{}, str string) error {
	buf := bytes.NewBufferString(str)
	if err := gob.NewDecoder(buf).Decode(subject); err != nil {
		return err
	}

	return nil
}

func (p *PageAwareCache) encodeInterfacePtrIntoString(subject interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(subject)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *PageAwareCache) decodeMapIntoMapPtr(subject interface{}, strMap map[string]string) error {
	if err := p.validateMapPointer(subject); err != nil {
		return err
	}

	mapVal := reflect.ValueOf(subject)
	mapElem := mapVal.Elem()
	mapType := mapElem.Type()

	placeholder := reflect.New(mapType.Elem())
	placeholderPtr := placeholder.Addr()

	for k, v := range strMap {
		buf := bytes.NewBufferString(v)
		if err := gob.NewDecoder(buf).DecodeValue(placeholderPtr); err != nil {
			return err
		}

		mapElem.SetMapIndex(
			reflect.ValueOf(k),
			placeholder,
		)
	}

	return nil
}

func (p *PageAwareCache) encodeMapPtrIntoMap(subject interface{}) (map[string]string, error) {
	if err := p.validateMapPointer(subject); err != nil {
		return nil, err
	}

	mapVal := reflect.ValueOf(subject)
	mapElem := mapVal.Elem()
	mapType := mapElem.Type()

	placeholder := reflect.New(mapType.Elem())
	placeholderPtr := placeholder.Addr()

	result := map[string]string{}
	iter := mapElem.MapRange()
	for iter.Next() {
		key := iter.Key().Interface().(string)
		placeholder.Set(iter.Value())

		buf := &bytes.Buffer{}
		err := gob.NewEncoder(buf).EncodeValue(placeholderPtr)
		if err != nil {
			return nil, err
		}

		result[key] = buf.String()
	}

	return result, nil
}

func (p *PageAwareCache) copyBetweenPointers(src, dest interface{}) error {
	if src == nil {
		return customerrors.ErrSourceValIsNil
	}
	if dest == nil {
		return customerrors.ErrDestinationValIsNil
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Ptr {
		return customerrors.ErrSourceValIsNotPtr
	}
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return customerrors.ErrDestinationValIsNotPtr
	}

	srcElem := srcVal.Elem()
	destElem := destVal.Elem()
	destElem.Set(srcElem)

	return nil
}

func (p *PageAwareCache) copyBetweenPointerMaps(srcMapPtr, destMapPtr interface{}) error {
	if srcMapPtr == nil {
		return customerrors.ErrSourceMapValIsNil
	}
	if destMapPtr == nil {
		return customerrors.ErrDestinationMapValIsNil
	}

	srcMapVal := reflect.ValueOf(srcMapPtr)
	if srcMapVal.Kind() != reflect.Ptr {
		return customerrors.ErrSourceMapValIsNotPtr
	}
	destMapVal := reflect.ValueOf(destMapPtr)
	if destMapVal.Kind() != reflect.Ptr {
		return customerrors.ErrDestinationMapValIsNotPtr
	}

	srcMapElem := srcMapVal.Elem()
	if srcMapElem.Kind() != reflect.Map {
		return customerrors.ErrSourceMapValIsNotMap
	}
	destMapElem := destMapVal.Elem()
	if destMapElem.Kind() != reflect.Map {
		return customerrors.ErrDestinationMapValIsNotMap
	}

	srcMapType := srcMapElem.Type()
	if srcMapType.Key().Kind() != reflect.String {
		return customerrors.ErrSourceMapKeyIsNotString
	}
	destMapType := srcMapElem.Type()
	if destMapType.Key().Kind() != reflect.String {
		return customerrors.ErrDestinationMapKeyIsNotString
	}

	if srcMapType.Elem() != destMapType.Elem() {
		return customerrors.ErrSourceDestinationMapValMismatch
	}

	iter := srcMapElem.MapRange()
	for iter.Next() {
		destMapElem.SetMapIndex(
			iter.Key(),
			iter.Value(),
		)
	}

	return nil
}

func (p *PageAwareCache) copyBetweenPointerLists(srcListPtr, destListPtr interface{}) error {
	if srcListPtr == nil {
		return customerrors.ErrSourceListValIsNil
	}
	if destListPtr == nil {
		return customerrors.ErrDestinationListValIsNil
	}

	srcListVal := reflect.ValueOf(srcListPtr)
	if srcListVal.Kind() != reflect.Ptr {
		return customerrors.ErrSourceListValIsNotPtr
	}
	destListVal := reflect.ValueOf(destListPtr)
	if destListVal.Kind() != reflect.Ptr {
		return customerrors.ErrDestinationListValIsNotPtr
	}

	srcListElemVal := srcListVal.Elem()
	if srcListElemVal.Kind() != reflect.Slice {
		return customerrors.ErrSourceListValIsNotSlice
	}
	destListElemVal := destListVal.Elem()
	if destListElemVal.Kind() != reflect.Slice {
		return customerrors.ErrDestinationListValIsNotSlice
	}

	targetLength := destListElemVal.Len()
	if srcListElemVal.Len() != targetLength {
		return customerrors.ErrDifferentLengthOfUnits
	}

	for i := 0; i < targetLength; i++ {
		srcIndexVal := srcListElemVal.Index(i)
		destIndexVal := destListElemVal.Index(i)

		destIndexVal.Set(srcIndexVal)
	}

	return nil
}

func (p *PageAwareCache) validateMapPointer(mapPtr interface{}) error {
	if mapPtr == nil {
		return customerrors.ErrSourceMapValIsNil
	}

	mapVal := reflect.ValueOf(mapPtr)
	if mapVal.Kind() != reflect.Ptr {
		return customerrors.ErrSourceMapValIsNotPtr
	}

	mapElem := mapVal.Elem()
	if mapElem.Kind() != reflect.Map {
		return customerrors.ErrSourceMapValIsNotMap
	}

	mapType := mapElem.Type()
	if mapType.Key().Kind() != reflect.String {
		return customerrors.ErrSourceMapKeyIsNotString
	}

	return nil
}
