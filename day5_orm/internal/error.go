package internal

import "errors"

var (
	ErrorEntityNotStruct = errors.New("entity is not a struct")
	ErrorValueIsZero     = errors.New("value is zero")
	ErrorFieldIsEmpty    = errors.New("field is empty")
	ErrorFieldCantSet    = errors.New("field can't set")
	ErrorEntityIsNil     = errors.New("entity is nil")
	ErrorTagFormat       = errors.New("tag format error")
	ErrorNoRows          = errors.New("no rows")
	ErrorInsertZeroRow   = errors.New("insert zero row")
	ErrorFieldUnknown    = errors.New("field unknown")
)
