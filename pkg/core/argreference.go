package core

import "fmt"

type ArgReference struct {
	key string
	arg interface{}
}

func NewArgReference(key string, arg interface{}) ArgReference {
	return ArgReference{
		key: key,
		arg: arg,
	}
}

func (ar ArgReference) Key() string {
	return ar.key
}

func (ar ArgReference) Argument() interface{} {
	return ar.arg
}

func stringifyArgReference(ar ArgReference, divider string) string {
	return fmt.Sprintf("%s%s%v", ar.arg, divider, ar.key)
}
