package valuer

import (
	"database/sql"
	rft "geektime-go/day5_orm/model"
)

//$ go test -bench=BenchmarkSetColumns -benchtime=10000x -benchmem
//                                                              cpu占用时间              内存分配            总分配次数
//BenchmarkSetColumns/reflect-8              10000              4489 ns/op             280 B/op         14 allocs/op
//BenchmarkSetColumns/unsafe-8               10000              2491 ns/op             208 B/op          6 allocs/op

// Creator 工厂方法
type Creator func(m *rft.Model, entity any) Valuer

type Valuer interface {
	Field(name string) (any, error)
	SetColumns(rows *sql.Rows) error
}
