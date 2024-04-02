package influxdb

func Count(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("COUNT", col)
}

func SUM(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("SUM", col)
}

func FIRST(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("FIRST", col)
}

func LAST(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("LAST", col)
}

func ABS(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("ABS", col)
}

func DERIVATIVE(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("DERIVATIVE", col)
}

func DIFFERENCE(col interface{}) SQLFunctionExpression {
	return newIdentifierFunc("DIFFERENCE", col)
}

func Interval(interval string) SQLFunctionExpression {
	return newIdentifierFunc("time", interval)
}

func Func(name string, args ...interface{}) SQLFunctionExpression {
	return NewSQLFunctionExpression(name, args...)
}

func newIdentifierFunc(name string, col interface{}) SQLFunctionExpression {
	if s, ok := col.(string); ok {
		col = I(s)
	}

	return Func(name, col)
}

func I(ident string) IdentifierExpression {
	return ParseIdentifier(ident)
}

func C(col string) IdentifierExpression {
	return NewIdentifierExpression("", "", col)
}

func Range(start, end any) RangeVal {
	return NewRangeVal(start, end)
}
