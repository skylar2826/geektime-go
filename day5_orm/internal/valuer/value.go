package valuer

import "database/sql"

type Valuer interface {
	SetColumns(rows sql.Rows) error
}
