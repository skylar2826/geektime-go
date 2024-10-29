package my_orm_mysql

type tableReference interface {
	//table() string
}

type Table struct {
	entity any
	alias  string
}

//
//func (t Table) table() string {
//	return ""
//}

func (t Table) C(name string) Column {
	return Column{
		Name:  name,
		table: t,
	}
}

func (t Table) join(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "Join",
	}
}

func (t Table) leftJoin(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "Left Join",
	}
}

func (t Table) rightJoin(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "Right Join",
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

type Join struct {
	left  tableReference
	right tableReference
	typ   string
	on    []Predicate
	using []string
}

func (j Join) table() string {
	//TODO implement me
	panic("implement me")
}

func (j Join) join(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "Join",
	}
}

func (j Join) leftJoin(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "Left Join",
	}
}

func (j Join) rightJoin(right tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "Right Join",
	}
}

type JoinBuilder struct {
	left  tableReference
	right tableReference
	typ   string
}

func (j *JoinBuilder) On(pres ...Predicate) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		on:    pres,
	}
}

// Using t1.Join(t2).Using("UserId)
func (j *JoinBuilder) Using(colNames ...string) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		using: colNames,
	}
}

type SubQuery struct {
	s       QueryBuilder
	columns []Selectable
	alias   string
	table   tableReference
}

func (SubQuery) expr() {}

func (s SubQuery) tableAlias() string {
	return s.alias
}

func (s SubQuery) Join(target tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "Join",
	}
}

func (s SubQuery) LeftJoin(target tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "Left Join",
	}
}

func (s SubQuery) RightJoin(target tableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "Right Join",
	}
}

func (s SubQuery) C(name string) Column {
	return Column{
		Name:  name,
		table: s.table,
	}
}

type SubQueryExpr struct {
	s    SubQuery
	pred string
}

func (s SubQueryExpr) expr() {}

func Any(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "Any",
	}
}

func All(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "All",
	}
}

func Some(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "Some",
	}
}
