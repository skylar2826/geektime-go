package my_orm_mysql

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRawQuery(t *testing.T) {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory", DBWithDialect(DialectMySql))
	require.NoError(t, err)
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "has args",
			builder: RawQuery[TestModel](db, "select * from `test_model` where `age` = ?;", 22),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where `age` = ?;",
				Args: []any{22},
			},
		},

		{
			name:    "no args",
			builder: RawQuery[TestModel](db, "select * from `test_model`;"),
			wantQuery: &Query{
				SQL: "select * from `test_model`;",
			},
		},
		{
			name:    "insert column values",
			builder: RawQuery[TestUser](db, "INSERT INTO `test_user`(`first_name`,`age`) VALUES (?,?) ON CONFLICT(`id`) DO UPDATE SET `first_name`=excluded.`first_name`,`age`=excluded.`age`;", "Tom", int8(18)),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_user`(`first_name`,`age`) VALUES (?,?) ON CONFLICT(`id`) DO UPDATE SET `first_name`=excluded.`first_name`,`age`=excluded.`age`;",
				Args: []any{"Tom", int8(18)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, b)
		})
	}
}
