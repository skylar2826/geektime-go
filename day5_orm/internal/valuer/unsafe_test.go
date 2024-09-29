package valuer

import "testing"

func TestUnSafeValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewUnsafeValue)
}
