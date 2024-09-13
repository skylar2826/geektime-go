package valuer

import (
	"database/sql"
	rft "geektime-go/day5_orm/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func testSetColumns(t *testing.T, creator Creator) {
	testCases := []struct {
		name       string
		entity     any
		wantEntity *TestModel
		wantErr    error
		rows       func() *sqlmock.Rows
	}{
		{
			name:   "set columns",
			entity: &TestModel{},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "John",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow(1, "John", 18, "Jerry")
				return rows
			},
		},
		{
			name:   "order",
			entity: &TestModel{},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "John",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"age", "last_name", "id", "first_name"})
				rows.AddRow(18, "Jerry", 1, "John")
				return rows
			},
		},
		{
			name:   "part column",
			entity: &TestModel{},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"age", "last_name", "id"})
				rows.AddRow(18, "Jerry", 1)
				return rows
			},
		},
	}

	r := rft.NewRegister()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRows := tc.rows()
			mock.ExpectQuery("Select XX").WillReturnRows(mockRows)
			var rows *sql.Rows
			rows, err = mockDB.Query("Select XX")
			require.NoError(t, err)
			rows.Next()

			var model *rft.Model
			model, err = r.Get(tc.entity)
			require.NoError(t, err)
			val := creator(model, tc.entity)
			err = val.SetColumns(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

func TestReflectValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewReflectValue)
}
