package valuer

import (
	"database/sql"
	"database/sql/driver"
	rft "geektime-go/day5_orm/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

// 基准测试
func BenchmarkSetColumns(b *testing.B) {
	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer mockDB.Close()

		mockRows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
		row := []driver.Value{1, "John", "Doe", 18}
		for n := 0; n < b.N; n++ {
			mockRows.AddRow(row...)
		}

		mock.ExpectQuery("Select xx").WillReturnRows(mockRows)
		var rows *sql.Rows
		rows, err = mockDB.Query("Select xx")
		require.NoError(b, err)

		r := rft.NewRegister()
		var m *rft.Model
		m, err = r.Get(&TestModel{})
		require.NoError(b, err)

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			rows.Next()
			val := creator(m, &TestModel{})
			err = val.SetColumns(rows)
		}
	}
	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})
	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}
