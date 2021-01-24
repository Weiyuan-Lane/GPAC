package errors

import (
	"errors"
)

var (
	ItemNotFoundErr = errors.New("GPAC: item was not found")

	SourceValIsNilErr         = errors.New("GPAC: source value is nil")
	DestinationValIsNilErr    = errors.New("GPAC: destination value is nil")
	SourceValIsNotPtrErr      = errors.New("GPAC: source value is not a pointer")
	DestinationValIsNotPtrErr = errors.New("GPAC: destination value is not a pointer")

	SourceListValIsNilErr           = errors.New("GPAC: source list value is nil")
	DestinationListValIsNilErr      = errors.New("GPAC: destination list value is nil")
	SourceListValIsNotPtrErr        = errors.New("GPAC: source list value is not a pointer")
	DestinationListValIsNotPtrErr   = errors.New("GPAC: destination list value is not a pointer")
	SourceListValIsNotSliceErr      = errors.New("GPAC: source list value is not a slice")
	DestinationListValIsNotSliceErr = errors.New("GPAC: destination list value is not a slice")

	DifferentLengthOfUnitsErr = = errors.New("GPAC: source and destination slices do not correspond to the same length")
)
