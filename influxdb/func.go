package influxdb

type SqlFunctionExpression struct {
	name string
	args []interface{}
}

func NewSQLFunctionExpression(name string, args ...interface{}) SQLFunctionExpression {
	return SqlFunctionExpression{
		name: name,
		args: args,
	}
}

func (sfe SqlFunctionExpression) Clone() Expression {
	return SqlFunctionExpression{
		name: sfe.name,
		args: sfe.args,
	}
}

func (sfe SqlFunctionExpression) Expression() Expression {
	return sfe
}

func (sfe SqlFunctionExpression) Name() string {
	return sfe.name
}

func (sfe SqlFunctionExpression) Args() []interface{} {
	return sfe.args
}

func (sfe SqlFunctionExpression) As(val interface{}) AliasedExpression {
	return NewAliasExpression(sfe, val)
}

func (sfe SqlFunctionExpression) Add(val interface{}) ComputerExpression {
	return NewComputerExpression(Plus, sfe, val)
}

func (sfe SqlFunctionExpression) Sub(val interface{}) ComputerExpression {
	return NewComputerExpression(Minus, sfe, val)
}

func (sfe SqlFunctionExpression) Mul(val interface{}) ComputerExpression {
	return NewComputerExpression(Multi, sfe, val)
}
