package influxdb

type ranged struct {
	lhs Expression
	rhs RangeVal
	op  RangeOperation
}

func NewRangeExpression(op RangeOperation, lhs Expression, rhs RangeVal) RangeExpression {
	return ranged{
		lhs: lhs,
		rhs: rhs,
		op:  op,
	}
}

func (r ranged) Clone() Expression {
	return NewRangeExpression(r.op, r.lhs, r.rhs)
}

func (r ranged) Expression() Expression {
	return r
}

func (r ranged) RHS() RangeVal {
	return r.rhs
}

func (r ranged) LHS() Expression {
	return r.lhs
}

func (r ranged) Op() RangeOperation {
	return r.op
}

func between(lhs Expression, rhs RangeVal) RangeExpression {
	return NewRangeExpression(BetweenOp, lhs, rhs)
}

func notBetween(lhs Expression, rhs RangeVal) RangeExpression {
	return NewRangeExpression(NotBetweenOp, lhs, rhs)
}

type rangeVal struct {
	start interface{}
	end   interface{}
}

func NewRangeVal(start, end interface{}) RangeVal {
	return rangeVal{
		start: start,
		end:   end,
	}
}

func (rv rangeVal) Start() interface{} {
	return rv.start
}

func (rv rangeVal) End() interface{} {
	return rv.end
}
