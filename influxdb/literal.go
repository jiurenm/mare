package influxdb

type LiteralExpression interface {
	Expression

	Literal() string
	Args() []any
}

type literal struct {
	literal string
	args    []any
}

func newLiteralExpression(sql string, args ...any) LiteralExpression {
	return literal{
		literal: sql,
		args:    args,
	}
}

func (l literal) Clone() Expression {
	return newLiteralExpression(l.literal, l.args...)
}

func (l literal) Expression() Expression {
	return l
}

func (l literal) Literal() string {
	return l.literal
}

func (l literal) Args() []any {
	return l.args
}

func Star() LiteralExpression {
	return newLiteralExpression("*")
}

func Now(s ...string) LiteralExpression {
	if len(s) > 0 {
		return newLiteralExpression("NOW" + s[0])
	}

	return newLiteralExpression("NOW")
}
