package valuer

import (
	"database/sql"
	"fmt"
	rft "geektime-go/day5_orm/reflect"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	val   any
	model *rft.Model
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(m *rft.Model, val any) Valuer {
	return &unsafeValue{val: val, model: m}
}

func (u *unsafeValue) SetColumns(rows *sql.Rows) error {

	address := reflect.ValueOf(u.val).UnsafePointer()
	cs, err := rows.Columns()
	var vals []any
	if err != nil {
		return err
	}

	for _, c := range cs {
		fd, ok := u.model.ColumnMap[c]
		if !ok {
			return fmt.Errorf("column %s not found", c)
		}
		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
		val := reflect.NewAt(fd.Typ, fdAddress).Interface()
		vals = append(vals, val)
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}
	return nil
}
