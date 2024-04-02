package influxdb

type expressionList struct {
	expressions []Expression
	operator    ExpressionListType
}

func NewExpressionList(operator ExpressionListType, expressions ...Expression) ExpressionList {
	el := expressionList{operator: operator}
	exps := make([]Expression, 0, len(el.expressions))

	for _, e := range expressions {
		switch t := e.(type) {
		case ExpressionList:
			if !t.IsEmpty() {
				exps = append(exps, e)
			}
		default:
			exps = append(exps, e)
		}
	}

	el.expressions = exps

	return el
}

func (el expressionList) Clone() Expression {
	newExps := make([]Expression, 0, len(el.expressions))
	for _, exp := range el.expressions {
		newExps = append(newExps, exp.Clone())
	}

	return expressionList{
		operator:    el.operator,
		expressions: newExps,
	}
}

func (el expressionList) Expression() Expression {
	return el
}

func (el expressionList) IsEmpty() bool {
	return len(el.expressions) == 0
}

func (el expressionList) Type() ExpressionListType {
	return el.operator
}

func (el expressionList) Expressions() []Expression {
	return el.expressions
}

func (el expressionList) Append(expressions ...Expression) ExpressionList {
	exps := make([]Expression, len(el.expressions)+len(expressions))
	copy(exps[:len(el.expressions)], el.expressions)
	copy(exps[len(el.expressions):], expressions)

	return NewExpressionList(el.operator, exps...)
}

func And(expressions ...Expression) ExpressionList {
	return NewExpressionList(AndType, expressions...)
}

func Or(expressions ...Expression) ExpressionList {
	return NewExpressionList(OrType, expressions...)
}
