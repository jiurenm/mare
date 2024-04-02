package influxdb

type orderedExpression struct {
	sortExpression Expression
	direction      SortDirection
}

func NewOrderedExpression(exp Expression, direction SortDirection) OrderedExpression {
	return orderedExpression{
		sortExpression: exp,
		direction:      direction,
	}
}

func (oe orderedExpression) Clone() Expression {
	return NewOrderedExpression(oe.sortExpression, oe.direction)
}

func (oe orderedExpression) Expression() Expression {
	return oe
}

func (oe orderedExpression) SortExpression() Expression {
	return oe.sortExpression
}

func (oe orderedExpression) IsAsc() bool {
	return oe.direction == AscDir
}

func asc(exp Expression) OrderedExpression {
	return NewOrderedExpression(exp, AscDir)
}

func desc(exp Expression) OrderedExpression {
	return NewOrderedExpression(exp, DescSortDir)
}
