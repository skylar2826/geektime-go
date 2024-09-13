package my_orm_mysql

// 衍生类型
type op string

func (op op) String() string {
	return string(op)
}

// 别名
//type op = string

var (
	opEq  op = "="
	opNot op = "not"
	opAnd op = "and"
	opOr  op = "or"
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
	name  string
	alias string
}

func (c Column) AS(name string) Column {
	return Column{
		name:  c.name,
		alias: name,
	}
}

func (Column) expr()       {}
func (Column) selectable() {}

func C(name string) Column {
	return Column{name: name}
}

type value struct {
	val any
}

func (value) expr() {}

// Eq C("id").Eq("5")
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(arg),
	}
}

func valueOf(val any) Expression {
	switch v := val.(type) {
	case Expression:
		return v
	default:
		return &value{val: v}
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

type Aggregate struct {
	fn    string
	arg   string
	alias string
}

func (Aggregate) selectable() {}

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
