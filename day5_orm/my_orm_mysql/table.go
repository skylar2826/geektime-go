package my_orm_mysql

type tableReference interface {
	table() string
}

type Table struct {
	entity any
	alias  string
}

func (t Table) table() string {
	return ""
}

func (t Table) C(name string) Column {
	return Column{
		Name:  name,
		table: t,
	}
}

func (t Table) join(right tableReference) *joinBuilder {
	return &joinBuilder{
		left:  t,
		right: right,
		typ:   "Join",
	}
}

func (t Table) leftJoin(right tableReference) *joinBuilder {
	return &joinBuilder{
		left:  t,
		right: right,
		typ:   "Left Join",
	}
}

func (t Table) rightJoin(right tableReference) *joinBuilder {
	return &joinBuilder{
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

func (j Join) join(right tableReference) *joinBuilder {
	return &joinBuilder{
		left:  j,
		right: right,
		typ:   "Join",
	}
}

func (j Join) leftJoin(right tableReference) *joinBuilder {
	return &joinBuilder{
		left:  j,
		right: right,
		typ:   "Left Join",
	}
}

func (j Join) rightJoin(right tableReference) *joinBuilder {
	return &joinBuilder{
		left:  j,
		right: right,
		typ:   "Right Join",
	}
}

type joinBuilder struct {
	left  tableReference
	right tableReference
	typ   string
}

func (j *joinBuilder) On(pres ...Predicate) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		on:    pres,
	}
}

// Using t1.Join(t2).Using("UserId)
func (j *joinBuilder) Using(colNames ...string) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		using: colNames,
	}
}
