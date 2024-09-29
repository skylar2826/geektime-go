package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

var a sql.NullString

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch data := src.(type) {
	case string:
		bs = []byte(data)
	case []byte:
		bs = data
	case nil:
		return nil
	default:
		return errors.New("sql: cannot scan value of type T")
	}

	err := json.Unmarshal(bs, &j.Val)
	if err != nil {
		return err
	}
	j.Valid = true
	return nil
}
