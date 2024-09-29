package test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type SimpleStruct struct {
	Id         uint64 `db:"id" json:"id"`
	BoolColumn bool   `db:"bool_column" json:"bool_column"`
	BoolPtr    *bool  `db:"bool_ptr" json:"bool_ptr,omitempty"`
	IntColumn  int    `db:"int_column" json:"int_column"`
	IntPtr     *int   `db:"int_ptr" json:"int_ptr,omitempty"`

	Int8Column uint8  `db:"int8_column" json:"int8_column"`
	Int8Ptr    *uint8 `db:"int8_ptr" json:"int8_ptr,omitempty"`

	Int16Column int16  `db:"int16_column" json:"int16_column"`
	Int16Ptr    *int16 `db:"int16_ptr" json:"int16_ptr,omitempty"`

	Int32Column int32  `db:"int32_column" json:"int32_column"`
	Int32Ptr    *int32 `db:"int32_ptr" json:"int32_ptr,omitempty"`

	Int64Column int64  `db:"int64_column" json:"int64_column"`
	Int64Ptr    *int64 `db:"int64_ptr" json:"int64_ptr,omitempty"`

	UintColumn uint64  `db:"uint_column" json:"uint_column"`
	UintPtr    *uint64 `db:"uint_ptr" json:"uint_ptr,omitempty"`

	Uint8Column uint8  `db:"uint8_column" json:"uint8_column"`
	Uint8Ptr    *uint8 `db:"uint8_ptr" json:"uint8_ptr,omitempty"`

	Uint16Column uint16  `db:"uint16_column" json:"uint16_column"`
	Uint16Ptr    *uint16 `db:"uint16_ptr" json:"uint16_ptr,omitempty"`

	Uint32Column uint32  `db:"uint32_column" json:"uint32_column"`
	Uint32Ptr    *uint32 `db:"uint32_ptr" json:"uint32_ptr,omitempty"` // 注意修正字段名

	Uint64Column uint64  `db:"uint64_column" json:"uint64_column"`     // 注意修正字段名
	Uint64Ptr    *uint64 `db:"uint64_ptr" json:"uint64_ptr,omitempty"` // 注意修正字段名

	Float32Column float32  `db:"float32_column" json:"float32_column"`
	Float32Ptr    *float32 `db:"float32_ptr" json:"float32_ptr,omitempty"`

	Float64Column float64  `db:"float64_column" json:"float64_column"`
	Float64Ptr    *float64 `db:"float64_ptr" json:"float64_ptr,omitempty"`

	ByteColumn uint8  `db:"byte_column" json:"byte_column"`
	BytePtr    *uint8 `db:"byte_ptr" json:"byte_ptr,omitempty"`

	ByteArray string `db:"byte_array" json:"byte_array"` // 使用string模拟TEXT

	StringColumn  string          `db:"string_column" json:"string_column"`
	NullStringPtr *sql.NullString `db:"null_string_ptr" json:"null_string_ptr,omitempty"`

	NullInt16Ptr *sql.NullInt16 `db:"null_int16_ptr" json:"null_int16_ptr,omitempty"`
	NullInt32Ptr *sql.NullInt32 `db:"null_int32_ptr" json:"null_int32_ptr,omitempty"`
	NullInt64Ptr *sql.NullInt64 `db:"null_int64_ptr" json:"null_int64_ptr,omitempty"`
	NullBoolPtr  *sql.NullBool  `db:"null_bool_ptr" json:"null_bool_ptr,omitempty"`
	//NullTimePtr    *sql.NullTime    `db:"null_time_ptr" json:"null_time_ptr,omitempty"`
	NullFloat64Ptr *sql.NullFloat64 `db:"null_float64_ptr" json:"null_float64_ptr,omitempty"`

	JsonColumn *JsonColumn `db:"json_column" json:"json_column"`
}

func NewSimpleStruct(id uint64) *SimpleStruct {
	//t := time.Unix(10, 100)
	return &SimpleStruct{
		Id:            id,
		BoolColumn:    true,
		IntColumn:     42,
		Int8Column:    1,
		Int16Column:   2,
		Int32Column:   3,
		Int64Column:   4,
		UintColumn:    5,
		Uint8Column:   6,
		Uint16Column:  7,
		Uint32Column:  8,
		Uint64Column:  9,
		Float32Column: 3.14,
		Float64Column: 2.718,
		ByteColumn:    'A',
		ByteArray:     "Hello, World!",
		StringColumn:  "Example String",
		NullStringPtr: &sql.NullString{String: "hello", Valid: true},
		NullInt16Ptr:  &sql.NullInt16{Int16: int16(16), Valid: true},
		NullInt32Ptr:  &sql.NullInt32{Int32: int32(32), Valid: true},
		//NullTimePtr:   &sql.NullTime{Time: t, Valid: true},
		JsonColumn: &JsonColumn{
			Val:   User{Name: "Tom"},
			Valid: true,
		},
	}
}

type User struct {
	Name string
}

type JsonColumn struct {
	Val   User
	Valid bool
}

func (j *JsonColumn) Scan(value any) error {
	if value == nil {
		j.Val, j.Valid = User{}, false
		return nil
	}

	j.Valid = true

	var bs []byte
	switch val := value.(type) {
	case string:
		bs = []byte(val)
	case []byte:
		bs = val
	}

	err := json.Unmarshal(bs, &j.Val)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonColumn) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	bs, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
