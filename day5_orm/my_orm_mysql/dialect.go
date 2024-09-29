package my_orm_mysql

import "errors"

type Dialect interface {
	quoter() byte
	upsert(b *Builder, upsert *Upsert) error
}

type standardDialect struct {
}

func (standardDialect) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (standardDialect) upsert(b *Builder, upsert *Upsert) error {
	//TODO implement me
	panic("implement me")
}

var (
	DialectMySql  Dialect = &mysqlDialect{}
	DialectSQLite Dialect = &sqliteDialect{}
)

//DialectPostgreSQL Dialect = &postgresDialect{}

type mysqlDialect struct {
}

func (m *mysqlDialect) quoter() byte {
	return '`'
}

func (m *mysqlDialect) upsert(b *Builder, upsert *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, a := range upsert.Assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.Name]
			if !ok {
				return errors.New("assign " + assign.Name + " is not found")
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sb.WriteString(")")

		case Assignment:
			fd, ok := b.model.FieldMap[assign.Column]
			if !ok {
				return errors.New("assign " + assign.Column + " is not found")
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.addArgs(assign.Val)
		default:
			return errors.New("invalid assignment")

		}
	}
	return nil
}

type sqliteDialect struct {
}

func (s *sqliteDialect) quoter() byte {
	return '`'
}

func (s *sqliteDialect) upsert(b *Builder, upsert *Upsert) error {
	b.sb.WriteString(" ON CONFLICT(")
	for idx, col := range upsert.conflictColumns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		err := b.buildColumn(Column{Name: col})
		if err != nil {
			return err
		}
	}
	b.sb.WriteString(") DO UPDATE SET ")
	for idx, a := range upsert.Assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.Name]
			if !ok {
				return errors.New("assign " + assign.Name + " is not found")
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=excluded.")
			b.quote(fd.ColName)

		case Assignment:
			fd, ok := b.model.FieldMap[assign.Column]
			if !ok {
				return errors.New("assign " + assign.Column + " is not found")
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.addArgs(assign.Val)
		default:
			return errors.New("invalid assignment")

		}
	}
	return nil
}
