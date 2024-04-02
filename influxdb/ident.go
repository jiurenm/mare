package influxdb

import "strings"

type IdentifierExpression interface {
	Expression
	Comparable
	Orderable
	Rangeable
	Computable

	Table(table string) IdentifierExpression
	GetTable() string
	Schema(schema string) IdentifierExpression
	GetSchema() string
	Col(col interface{}) IdentifierExpression
	GetCol() interface{}

	IsEmpty() bool
}

type identifier struct {
	col    interface{}
	schema string
	table  string
}

const (
	tableAndColumnParts                 = 2
	schemaTableAndColumnIdentifierParts = 3
)

func ParseIdentifier(ident string) IdentifierExpression {
	parts := strings.Split(ident, ".")
	switch len(parts) {
	case tableAndColumnParts:
		return NewIdentifierExpression("", parts[0], parts[1])
	case schemaTableAndColumnIdentifierParts:
		return NewIdentifierExpression(parts[0], parts[1], parts[2])
	}

	return NewIdentifierExpression("", "", ident)
}

func NewIdentifierExpression(schema, table string, col interface{}) IdentifierExpression {
	return identifier{}.Schema(schema).Table(table).Col(col)
}

func (i identifier) clone() identifier {
	return identifier{schema: i.schema, table: i.table, col: i.col}
}

func (i identifier) Clone() Expression {
	return i.clone()
}

func (i identifier) Expression() Expression { return i }

func (i identifier) Table(table string) IdentifierExpression {
	i.table = table

	return i
}

func (i identifier) GetTable() string {
	return i.table
}

func (i identifier) Schema(schema string) IdentifierExpression {
	i.schema = schema

	return i
}

func (i identifier) GetSchema() string {
	return i.schema
}

func (i identifier) Col(col interface{}) IdentifierExpression {
	if col == "*" {
		i.col = Star()
	} else {
		i.col = col
	}

	return i
}

func (i identifier) GetCol() interface{} {
	return i.col
}

func (i identifier) IsEmpty() bool {
	isEmpty := i.schema == "" && i.table == ""
	if isEmpty {
		switch t := i.col.(type) {
		case nil:
			return true
		case string:
			return t == ""
		default:
			return false
		}
	}

	return isEmpty
}

func (i identifier) Eq(val interface{}) BooleanExpression {
	return eq(i, val)
}

func (i identifier) Neq(val interface{}) BooleanExpression {
	return neq(i, val)
}

func (i identifier) Gt(val interface{}) BooleanExpression {
	return gt(i, val)
}

func (i identifier) Gte(val interface{}) BooleanExpression {
	return gte(i, val)
}

func (i identifier) Lt(val interface{}) BooleanExpression {
	return lt(i, val)
}

func (i identifier) Lte(val interface{}) BooleanExpression {
	return lte(i, val)
}

func (i identifier) Asc() OrderedExpression {
	return asc(i)
}

func (i identifier) Desc() OrderedExpression {
	return desc(i)
}

func (i identifier) Between(val RangeVal) RangeExpression {
	return between(i, val)
}

func (i identifier) NotBetween(val RangeVal) RangeExpression {
	return notBetween(i, val)
}

func (i identifier) Add(val interface{}) ComputerExpression {
	return add(i, val)
}

func (i identifier) Sub(val interface{}) ComputerExpression {
	return sub(i, val)
}

func (i identifier) Mul(val interface{}) ComputerExpression {
	return mul(i, val)
}
