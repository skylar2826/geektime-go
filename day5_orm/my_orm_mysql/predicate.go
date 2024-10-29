package my_orm_mysql

// 衍生类型
type op string

func (op op) String() string {
	return string(op)
}

// 别名
//type op = string

var (
	opEq    op = "="
	opNot   op = "not"
	opAnd   op = "and"
	opOr    op = "or"
	opLt    op = "<"
	opIn    op = "in"
	opExist op = "exist"
)

// Expression 是标记接口，代表表达式;
// 实现一个无用的标记接口，以赋予这个结构体特定含义
type Expression interface {
	expr()
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) expr() {}

type Column struct {
	Name  string
	alias string
	table tableReference // join查询时，每个table需要单独维护自己的model
}

func (c Column) AS(name string) Column {
	return Column{
		Name:  c.Name,
		alias: name,
		table: c.table,
	}
}

func (Column) expr()       {}
func (Column) selectable() {}
func (Column) assign()     {}

func C(name string) Column {
	return Column{Name: name}
}

func (c Column) InQuery(sub SubQuery) Predicate {
	return Predicate{
		left:  c,
		op:    opIn,
		right: sub,
	}
}

// Eq C("id").Eq("5")
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(arg),
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valueOf(arg),
	}
}

func valueOf(val any) Expression {
	switch v := val.(type) {
	case Expression:
		return v
	default:
		return value{val: v}
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

// And C("id").Eq("5").And(C("name").Eq("lili"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

type value struct {
	val any
}

func (value) expr() {}

type Aggregate struct {
	fn    string
	arg   string
	alias string
}

func (Aggregate) selectable() {}
func (Aggregate) expr()       {}

func Avg(col string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: col,
	}
}

func (a Aggregate) AS(name string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: name,
	}
}

func (a Aggregate) Eq(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opEq,
		right: valueOf(arg),
	}
}

func (a Aggregate) Lt(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opLt,
		right: valueOf(arg),
	}
}

type RawExpr struct {
	expression string
	args       []any
}

func (RawExpr) expr()       {}
func (RawExpr) selectable() {}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		expression: expr,
		args:       args,
	}
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

type OrderBy struct {
	col   Column
	order string
}

func ASC(column Column) OrderBy {
	return OrderBy{
		col:   column,
		order: "ASC",
	}
}

func DESC(column Column) OrderBy {
	return OrderBy{
		col:   column,
		order: "DESC",
	}
}

func Exist(sub SubQuery) Predicate {
	return Predicate{
		op:    opExist,
		right: sub,
	}
}
