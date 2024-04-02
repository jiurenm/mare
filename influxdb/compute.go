package influxdb

type computer struct {
	lhs Expression
	rhs interface{}
	op  Operator
}

func NewComputerExpression(op Operator, lhs Expression, rhs interface{}) ComputerExpression {
	return computer{
		lhs: lhs,
		rhs: rhs,
		op:  op,
	}
}

func (c computer) Clone() Expression {
	return nil
}

func (c computer) Expression() Expression {
	return c
}

func (c computer) Op() Operator {
	return c.op
}

func (c computer) RHS() interface{} {
	return c.rhs
}

func (c computer) LHS() Expression {
	return c.lhs
}

func (c computer) As(val interface{}) AliasedExpression {
	return NewAliasExpression(c, val)
}

func (c computer) Add(val interface{}) ComputerExpression {
	return NewComputerExpression(Plus, c, val)
}

func (c computer) Sub(val interface{}) ComputerExpression {
	return NewComputerExpression(Minus, c, val)
}

func (c computer) Mul(val interface{}) ComputerExpression {
	return NewComputerExpression(Multi, c, val)
}

func add(lhs Expression, rhs interface{}) ComputerExpression {
	return NewComputerExpression(Plus, lhs, rhs)
}

func sub(lhs Expression, rhs interface{}) ComputerExpression {
	return NewComputerExpression(Minus, lhs, rhs)
}

func mul(lhs Expression, rhs interface{}) ComputerExpression {
	return NewComputerExpression(Multi, lhs, rhs)
}
