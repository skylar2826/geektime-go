package my_orm_mysql

import (
	"database/sql"
	"geektime-go/day5_orm/internal"
	rft "geektime-go/day5_orm/reflect"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestUser struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestInsert_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := rft.OpenDB(mockDB)
	testCases := []struct {
		name    string
		i       *Insert[TestUser]
		wantErr error
		wantRes *Query
	}{
		{
			name:    "invalid row",
			i:       NewInsert[TestUser](db).Values(),
			wantErr: internal.ErrorInsertZeroRow,
		},
		{
			name: "insert row",
			i: NewInsert[TestUser](db).Values(&TestUser{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "xi"},
			}),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_user`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?);",
				Args: []any{int64(1), "Tom", int8(18), &sql.NullString{Valid: true, String: "xi"}},
			},
		},
		{
			name: "insert part row",
			i: NewInsert[TestUser](db).Columns("FirstName", "Age").Values(&TestUser{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "xi"},
			}),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_user`(`first_name`,`age`) VALUES (?,?);",
				Args: []any{"Tom", int8(18)},
			},
		},
		{
			name: "insert rows",
			i: NewInsert[TestUser](db).Values(&TestUser{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "xi"},
			}, &TestUser{
				Id:        12,
				FirstName: "Tom1",
				Age:       19,
				LastName:  &sql.NullString{Valid: true, String: "x"},
			}),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_user`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{int64(1), "Tom", int8(18), &sql.NullString{Valid: true, String: "xi"}, int64(12), "Tom1", int8(19), &sql.NullString{Valid: true, String: "x"}},
			},
		},
		{
			name: "insert assignment",
			i: NewInsert[TestUser](db).Columns("FirstName", "Age").Values(&TestUser{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "xi"},
			}).OnDuplicateKey().Update(Assign("FirstName", "lili"), Assign("Age", 88)),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_user`(`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `first_name`=?,`age`=?;",
				Args: []any{"Tom", int8(18), "lili", 88},
			},
		},
		{
			name: "insert column values",
			i: NewInsert[TestUser](db).Columns("FirstName", "Age").Values(&TestUser{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "xi"},
			}).OnDuplicateKey().Update(C("FirstName"), C("Age")),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_user`(`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`),`age`=VALUES(`age`);",
				Args: []any{"Tom", int8(18)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, q, tc.wantRes)
		})
	}
}
