package errors

import (
	"errors"
)

var (
	ErrItemNotFound = errors.New("GPAC: item was not found")

	ErrSourceValIsNil         = errors.New("GPAC: source value is nil")
	ErrDestinationValIsNil    = errors.New("GPAC: destination value is nil")
	ErrSourceValIsNotPtr      = errors.New("GPAC: source value is not a pointer")
	ErrDestinationValIsNotPtr = errors.New("GPAC: destination value is not a pointer")

	ErrSourceListValIsNil           = errors.New("GPAC: source list value is nil")
	ErrDestinationListValIsNil      = errors.New("GPAC: destination list value is nil")
	ErrSourceListValIsNotPtr        = errors.New("GPAC: source list value is not a pointer")
	ErrDestinationListValIsNotPtr   = errors.New("GPAC: destination list value is not a pointer")
	ErrSourceListValIsNotSlice      = errors.New("GPAC: source list value is not a slice")
	ErrDestinationListValIsNotSlice = errors.New("GPAC: destination list value is not a slice")

	ErrSourceMapValIsNil               = errors.New("GPAC: source map value is nil")
	ErrDestinationMapValIsNil          = errors.New("GPAC: destination map value is nil")
	ErrSourceMapValIsNotPtr            = errors.New("GPAC: source map value is not a pointer")
	ErrDestinationMapValIsNotPtr       = errors.New("GPAC: destination map value is not a pointer")
	ErrSourceMapValIsNotMap            = errors.New("GPAC: source map value is not a map")
	ErrDestinationMapValIsNotMap       = errors.New("GPAC: destination map value is not a map")
	ErrSourceMapKeyIsNotInt            = errors.New("GPAC: source map key is not a int")
	ErrDestinationMapKeyIsNotInt       = errors.New("GPAC: destination map key is not a int")
	ErrSourceDestinationMapValMismatch = errors.New("GPAC: source map val type does not equate to destination map val type")

	ErrDifferentLengthOfUnits = errors.New("GPAC: source and destination slices do not correspond to the same length")
)
