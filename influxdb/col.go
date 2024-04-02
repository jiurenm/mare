package influxdb

type columnList struct {
	columns []Expression
}

func newColumnListExpression(vals ...any) ColumnListExpression {
	var cols []Expression

	for _, val := range vals {
		switch t := val.(type) {
		case string:
			cols = append(cols, ParseIdentifier(t))
		case ColumnListExpression:
			cols = append(cols, t.Columns()...)
		case Expression:
			cols = append(cols, t)
		}
	}

	return columnList{columns: cols}
}

func newOrderedColumnList(vals ...OrderedExpression) ColumnListExpression {
	exps := make([]interface{}, 0, len(vals))
	for _, col := range vals {
		exps = append(exps, col.Expression())
	}

	return newColumnListExpression(exps...)
}

func (cl columnList) Clone() Expression {
	newExps := make([]Expression, 0, len(cl.columns))
	for _, exp := range cl.columns {
		newExps = append(newExps, exp.Clone())
	}

	return columnList{columns: nil}
}

func (cl columnList) Expression() Expression {
	return cl
}

func (cl columnList) IsEmpty() bool {
	return len(cl.columns) == 0
}

func (cl columnList) Columns() []Expression {
	return cl.columns
}

func (cl columnList) Append(cols ...Expression) ColumnListExpression {
	ret := columnList{}
	exps := append(ret.columns, cl.columns...)
	exps = append(exps, cols...)
	ret.columns = exps

	return ret
}
