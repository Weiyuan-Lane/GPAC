package gpac

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
	return fmt.Sprintf(constants.SimpleItemKeyTemplate, p.uniqueNamespace, itemKey)
}

func (p *PageAwareCache) createItemCacheKeyFromSubKeys(subKeys ...ArgReference) string {
	return p.createCacheKeyFromSubKeysAndTemplate(
		constants.ItemKeyTemplate,
		subKeys...,
	)
}

// Using the unique namespace and the page key, create a key unique
// to this target resource page
func (p *PageAwareCache) createPageCacheKeyFromStrKey(pageKey string) string {
	return fmt.Sprintf(constants.SimplePageKeyTemplate, p.uniqueNamespace, pageKey)
}

func (p *PageAwareCache) createPageCacheKeyFromSubKeys(subKeys ...ArgReference) string {
	return p.createCacheKeyFromSubKeysAndTemplate(
		constants.PageKeyTemplate,
		subKeys...,
	)
}

func (p *PageAwareCache) createCacheKeyFromSubKeysAndTemplate(template string, subKeys ...ArgReference) string {
	cacheSubKeys := make([]string, len(subKeys))
	cacheSubNamespaces := make([]string, len(subKeys))
	for i, subKey := range subKeys {
		cacheSubKeys[i] = stringifyArgReference(subKey, constants.ArgDivider)
		cacheSubNamespaces[i] = subKey.Key()
	}

	cacheKeyComponent := strings.Join(cacheSubKeys, constants.SubkeyDivider)
	cacheNamespaceComponent := strings.Join(cacheSubNamespaces, constants.ArgDivider)

	return fmt.Sprintf(
		template,
		p.uniqueNamespace,
		cacheNamespaceComponent,
		cacheKeyComponent,
	)
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

func (p *PageAwareCache) decodeMapIntoMapPtr(subject interface{}, strMap map[int]string) error {
	if err := p.validateMapPointer(subject); err != nil {
		return err
	}

	mapVal := reflect.ValueOf(subject)
	mapElem := mapVal.Elem()
	mapType := mapElem.Type()

	placeholderPtr := reflect.New(mapType.Elem())
	placeholder := reflect.Indirect(placeholderPtr)

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

func (p *PageAwareCache) encodeMapPtrIntoMap(subject interface{}) (map[int]string, error) {
	if err := p.validateMapPointer(subject); err != nil {
		return nil, err
	}

	mapVal := reflect.ValueOf(subject)
	mapElem := mapVal.Elem()
	mapType := mapElem.Type()

	placeholderPtr := reflect.New(mapType.Elem())
	placeholder := reflect.Indirect(placeholderPtr)

	result := map[int]string{}
	iter := mapElem.MapRange()
	for iter.Next() {
		key := iter.Key().Interface().(int)
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
	if srcMapType.Key().Kind() != reflect.Int {
		return customerrors.ErrSourceMapKeyIsNotInt
	}
	destMapType := srcMapElem.Type()
	if destMapType.Key().Kind() != reflect.Int {
		return customerrors.ErrDestinationMapKeyIsNotInt
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
	if mapType.Key().Kind() != reflect.Int {
		return customerrors.ErrSourceMapKeyIsNotInt
	}

	return nil
}

func (p *PageAwareCache) makePointerTo(subject interface{}) interface{} {
	subjectPtr := reflect.New(reflect.TypeOf(subject))
	reflect.Indirect(subjectPtr).Set(reflect.ValueOf(subject))

	return subjectPtr.Interface()
}
