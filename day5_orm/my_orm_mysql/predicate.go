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
	name string
}

func (Column) expr() {}

func C(name string) Column {
	return Column{name}
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
		right: value{val: arg},
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
