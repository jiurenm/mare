package influxdb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

var (
	ErrEmptyIdentifier = errors.New(
		`a empty identifier was encountered, please specify a "schema", "table" or "column"`,
	)
	ErrUnexpectedNamedWindow = errors.New(`unexpected named window function`)
	ErrEmptyCaseWhens        = errors.New(`when conditions not found for case statement`)
)

func errUnsupportedExpressionType(e Expression) error {
	return fmt.Errorf("unsupported expression type %T", e)
}

func errUnsupportedIdentifierExpression(t interface{}) error {
	return fmt.Errorf("unexpected col type must be string or LiteralExpression received %T", t)
}

func errUnsupportedBooleanExpressionOperator(op BooleanOperation) error {
	return fmt.Errorf("boolean operator '%+v' not supported", op)
}

func errUnsupportedComputerExpressionOperator(op Operator) error {
	return fmt.Errorf("operator '%+v' not supported", op)
}

func errUnsupportedRangeExpressionOperator(op RangeOperation) error {
	return fmt.Errorf("range operator %+v not supported", op)
}

const (
	// =.
	EqOp BooleanOperation = iota
	// != or <>.
	NeqOp
	// IS.
	IsOp
	// IS NOT.
	IsNotOp
	// >.
	GtOp
	// >=.
	GteOp
	// <.
	LtOp
	// <=.
	LteOp
	// IN.
	InOp
	// NOT IN.
	NotInOp
	// LIKE, LIKE BINARY...
	LikeOp
	// NOT LIKE, NOT LIKE BINARY...
	NotLikeOp
	// ~, REGEXP BINARY.
	RegexpLikeOp
	// !~, NOT REGEXP BINARY.
	RegexpNotLikeOp
	// ~*, REGEXP.
	RegexpILikeOp
	// !~*, NOT REGEXP.
	RegexpNotILikeOp

	AndType ExpressionListType = iota
	OrType

	// BETWEEN.
	BetweenOp RangeOperation = iota
	// NOT BETWEEN.
	NotBetweenOp

	// ASC.
	AscDir SortDirection = iota
	// DESC.
	DescSortDir

	Plus Operator = iota
	Minus
	Multi
)

var operatorInversions = map[BooleanOperation]BooleanOperation{
	IsOp:             IsNotOp,
	EqOp:             NeqOp,
	GtOp:             LteOp,
	GteOp:            LtOp,
	LtOp:             GteOp,
	LteOp:            GtOp,
	InOp:             NotInOp,
	LikeOp:           NotLikeOp,
	RegexpLikeOp:     RegexpNotLikeOp,
	RegexpILikeOp:    RegexpNotILikeOp,
	IsNotOp:          IsOp,
	NeqOp:            EqOp,
	NotInOp:          InOp,
	NotLikeOp:        LikeOp,
	RegexpNotLikeOp:  RegexpLikeOp,
	RegexpNotILikeOp: RegexpILikeOp,
}

type Expression interface {
	Clone() Expression
	Expression() Expression
}

type ColumnListExpression interface {
	Expression
	Columns() []Expression
	IsEmpty() bool
	Append(...Expression) ColumnListExpression
}

type ExpressionSQLGenerator interface {
	Generate(sb SQLBuilder, val any)
}

type expressionSQLGenerator struct {
	dialectOptions *SQLDialectOptions
}

func newExpressionSQLGenerator(do *SQLDialectOptions) ExpressionSQLGenerator {
	return &expressionSQLGenerator{dialectOptions: do}
}

func (esg *expressionSQLGenerator) Generate(sb SQLBuilder, val any) {
	if sb.Error() != nil {
		return
	}

	if val == nil {
		esg.literalNil(sb)

		return
	}

	switch v := val.(type) {
	case Expression:
		esg.expressionSQL(sb, v)
	case int:
		esg.literalInt(sb, int64(v))
	case int32:
		esg.literalInt(sb, int64(v))
	case int64:
		esg.literalInt(sb, v)
	case float32:
		esg.literalFloat(sb, float64(v))
	case float64:
		esg.literalFloat(sb, v)
	case string:
		esg.literalString(sb, v)
	default:
		esg.reflectSQL(sb, v)
	}
}

func (esg *expressionSQLGenerator) expressionSQL(sb SQLBuilder, expression Expression) {
	switch e := expression.(type) {
	case ColumnListExpression:
		esg.columnListSQL(sb, e)
	case ExpressionList:
		esg.expressionListSQL(sb, e)
	case LiteralExpression:
		esg.literalExpressionSQL(sb, e)
	case IdentifierExpression:
		esg.identifierExpressionSQL(sb, e)
	case SqlFunctionExpression:
		esg.sqlFunctionExpressionSQL(sb, e)
	case AliasedExpression:
		esg.aliasedExpressionSQL(sb, e)
	case BooleanExpression:
		esg.booleanExpressionSQL(sb, e)
	case ComputerExpression:
		esg.computerExpressionSQL(sb, e)
	case RangeExpression:
		esg.rangeExpressionSQL(sb, e)
	case OrderedExpression:
		esg.orderedExpressionSQL(sb, e)
	case Ex:
		esg.expressionMapSQL(sb, e)
	case ExOr:
		esg.expressionOrMapSQL(sb, e)
	default:
		sb.SetError(errUnsupportedExpressionType(e))
	}
}

func (esg *expressionSQLGenerator) columnListSQL(sb SQLBuilder, columnList ColumnListExpression) {
	cols := columnList.Columns()
	colLen := len(cols)

	for i, col := range cols {
		esg.Generate(sb, col)

		if i < colLen-1 {
			sb.WriteRunes(esg.dialectOptions.CommaRune, esg.dialectOptions.SpaceRune)
		}
	}
}

func (esg *expressionSQLGenerator) expressionListSQL(sb SQLBuilder, expressionList ExpressionList) {
	if expressionList.IsEmpty() {
		return
	}

	var op []byte

	if expressionList.Type() == AndType {
		op = esg.dialectOptions.AndFragment
	} else {
		op = esg.dialectOptions.OrFragment
	}

	exps := expressionList.Expressions()
	expLen := len(exps) - 1

	if expLen > 0 {
		sb.WriteRunes(esg.dialectOptions.LeftParenRune)
	} else {
		esg.Generate(sb, exps[0])

		return
	}

	for i, e := range exps {
		esg.Generate(sb, e)

		if i < expLen {
			sb.Write(op)
		}
	}

	sb.WriteRunes(esg.dialectOptions.RightParenRune)
}

func (esg *expressionSQLGenerator) literalExpressionSQL(sb SQLBuilder, literal LiteralExpression) {
	l := literal.Literal()
	args := literal.Args()

	if argsLen := len(args); argsLen > 0 {
		currIndex := 0
		for _, char := range l {
			if char == '?' && currIndex < argsLen {
				esg.Generate(sb, args[currIndex])
				currIndex++
			} else {
				sb.WriteRunes(char)
			}
		}

		return
	}

	sb.WriteStrings(l)
}

func (esg *expressionSQLGenerator) identifierExpressionSQL(sb SQLBuilder, ident IdentifierExpression) {
	if ident.IsEmpty() {
		sb.SetError(ErrEmptyIdentifier)

		return
	}

	schema, table, col := ident.GetSchema(), ident.GetTable(), ident.GetCol()
	if schema != esg.dialectOptions.EmptyString {
		sb.WriteStrings(schema)
	}

	if table != esg.dialectOptions.EmptyString {
		if schema != esg.dialectOptions.EmptyString {
			sb.WriteRunes(esg.dialectOptions.PeriodRune)
		}

		sb.WriteStrings(table)
	}

	switch t := col.(type) {
	case nil:
	case string:
		if col != esg.dialectOptions.EmptyString {
			if table != esg.dialectOptions.EmptyString || schema != esg.dialectOptions.EmptyString {
				sb.WriteRunes(esg.dialectOptions.PeriodRune)
			}

			sb.WriteStrings(t)
		}
	case LiteralExpression:
		if table != esg.dialectOptions.EmptyString || schema != esg.dialectOptions.EmptyString {
			sb.WriteRunes(esg.dialectOptions.PeriodRune)
		}

		esg.Generate(sb, t)
	default:
		sb.SetError(errUnsupportedIdentifierExpression(col))
	}
}

func (esg *expressionSQLGenerator) sqlFunctionExpressionSQL(sb SQLBuilder, sqlFunc SQLFunctionExpression) {
	sb.WriteStrings(sqlFunc.Name())
	esg.Generate(sb, sqlFunc.Args())
}

func (esg *expressionSQLGenerator) aliasedExpressionSQL(sb SQLBuilder, aliased AliasedExpression) {
	esg.Generate(sb, aliased.Aliased())
	sb.Write(esg.dialectOptions.AsFragment)
	esg.Generate(sb, aliased.GetAs())
}

func (esg *expressionSQLGenerator) booleanExpressionSQL(sb SQLBuilder, operator BooleanExpression) {
	sb.WriteRunes(esg.dialectOptions.LeftParenRune)
	esg.Generate(sb, operator.LHS())
	sb.WriteRunes(esg.dialectOptions.SpaceRune)

	operatorOp := operator.Op()
	if val, ok := esg.dialectOptions.BooleanOperatorLookup[operatorOp]; ok {
		sb.Write(val)
	} else {
		sb.SetError(errUnsupportedBooleanExpressionOperator(operatorOp))

		return
	}

	rhs := operator.RHS()

	sb.WriteRunes(esg.dialectOptions.SpaceRune)
	esg.Generate(sb, rhs)
	sb.WriteRunes(esg.dialectOptions.RightParenRune)
}

func (esg *expressionSQLGenerator) computerExpressionSQL(sb SQLBuilder, operator ComputerExpression) {
	esg.Generate(sb, operator.LHS())
	op := operator.Op()

	if val, ok := esg.dialectOptions.ComputeOperatorLookup[op]; ok {
		sb.Write(val)
	} else {
		sb.SetError(errUnsupportedComputerExpressionOperator(op))

		return
	}

	esg.Generate(sb, operator.RHS())
}

func (esg *expressionSQLGenerator) rangeExpressionSQL(sb SQLBuilder, operator RangeExpression) {
	sb.WriteRunes(esg.dialectOptions.LeftParenRune)

	lhs, rhs := operator.LHS(), operator.RHS()

	if operator.Op() == BetweenOp {
		esg.Generate(sb, lhs)
		sb.Write(esg.dialectOptions.BooleanOperatorLookup[GteOp])
		esg.Generate(sb, rhs.Start())
		sb.Write(esg.dialectOptions.AndFragment)
		esg.Generate(sb, lhs)
		sb.Write(esg.dialectOptions.BooleanOperatorLookup[LteOp])
		esg.Generate(sb, rhs.End())
		sb.WriteRunes(esg.dialectOptions.RightParenRune)
	} else if operator.Op() == NotBetweenOp {
		esg.Generate(sb, lhs)
		sb.Write(esg.dialectOptions.BooleanOperatorLookup[LtOp])
		esg.Generate(sb, rhs.Start())
		sb.Write(esg.dialectOptions.AndFragment)
		esg.Generate(sb, lhs)
		sb.Write(esg.dialectOptions.BooleanOperatorLookup[GtOp])
		esg.Generate(sb, rhs.End())
		sb.WriteRunes(esg.dialectOptions.RightParenRune)
	} else {
		sb.SetError(errUnsupportedRangeExpressionOperator(operator.Op()))

		return
	}
}

func (esg *expressionSQLGenerator) orderedExpressionSQL(sb SQLBuilder, order OrderedExpression) {
	esg.Generate(sb, order.SortExpression())

	if order.IsAsc() {
		sb.Write(esg.dialectOptions.AscFragment)
	} else {
		sb.Write(esg.dialectOptions.DescFragment)
	}
}

func (esg *expressionSQLGenerator) expressionMapSQL(sb SQLBuilder, ex Ex) {
	list, err := ex.ToExpressions()
	if err != nil {
		sb.SetError(err)

		return
	}

	esg.Generate(sb, list)
}

func (esg *expressionSQLGenerator) expressionOrMapSQL(sb SQLBuilder, ex ExOr) {
	list, err := ex.ToExpressions()
	if err != nil {
		sb.SetError(err)

		return
	}

	esg.Generate(sb, list)
}

func (esg *expressionSQLGenerator) literalNil(sb SQLBuilder) {
	sb.Write(esg.dialectOptions.Null)
}

func (esg *expressionSQLGenerator) literalInt(sb SQLBuilder, i int64) {
	sb.WriteStrings(strconv.FormatInt(i, 10))
}

func (esg *expressionSQLGenerator) literalFloat(sb SQLBuilder, f float64) {
	sb.WriteStrings(strconv.FormatFloat(f, 'f', -1, 64))
}

func (esg *expressionSQLGenerator) literalString(sb SQLBuilder, s string) {
	sb.WriteRunes(esg.dialectOptions.StringQuote)

	for _, char := range s {
		if e, ok := esg.dialectOptions.EscapedRunes[char]; ok {
			sb.Write(e)
		} else {
			sb.WriteRunes(char)
		}
	}

	sb.WriteRunes(esg.dialectOptions.StringQuote)
}

func (esg *expressionSQLGenerator) literalBytes(sb SQLBuilder, bs []byte) {
	sb.WriteRunes(esg.dialectOptions.StringQuote)

	i := 0

	for len(bs) > 0 {
		char, l := utf8.DecodeRune(bs)
		if e, ok := esg.dialectOptions.EscapedRunes[char]; ok {
			sb.Write(e)
		} else {
			sb.WriteRunes(char)
		}
		i++

		bs = bs[l:]
	}
	sb.WriteRunes(esg.dialectOptions.StringQuote)
}

func (esg *expressionSQLGenerator) sliceValueSQL(sb SQLBuilder, slice reflect.Value) {
	sb.WriteRunes(esg.dialectOptions.LeftParenRune)

	for i, l := 0, slice.Len(); i < l; i++ {
		esg.Generate(sb, slice.Index(i).Interface())

		if i < l-1 {
			sb.WriteRunes(esg.dialectOptions.CommaRune, esg.dialectOptions.SpaceRune)
		}
	}
	sb.WriteRunes(esg.dialectOptions.RightParenRune)
}

func (esg *expressionSQLGenerator) reflectSQL(sb SQLBuilder, val interface{}) {
	v := reflect.Indirect(reflect.ValueOf(val))
	valKind := v.Kind()

	switch {
	case IsInvalid(valKind):
		esg.literalNil(sb)
	case IsSlice(valKind):
		switch t := val.(type) {
		case []byte:
			esg.literalBytes(sb, t)
		default:
			esg.sliceValueSQL(sb, v)
		}
	case IsInt(valKind):
		esg.Generate(sb, v.Int())
	case IsUint(valKind):
		esg.Generate(sb, int64(v.Uint()))
	case IsFloat(valKind):
		esg.Generate(sb, v.Float())
	case IsString(valKind):
		esg.Generate(sb, v.String())
	default:
		sb.SetError(fmt.Errorf("encode error: Unable to encode value %+v", v))
	}
}

type SQLFunctionExpression interface {
	Expression
	Aliaseable
	Computable

	Name() string
	Args() []interface{}
}

type ExpressionListType int

type ExpressionList interface {
	Expression
	Type() ExpressionListType
	Expressions() []Expression
	Append(...Expression) ExpressionList
	IsEmpty() bool
}

type (
	BooleanOperation  int
	BooleanExpression interface {
		Expression
		Op() BooleanOperation
		LHS() Expression
		RHS() interface{}
	}
	Comparable interface {
		Eq(interface{}) BooleanExpression
		Neq(interface{}) BooleanExpression
		Gt(interface{}) BooleanExpression
		Gte(interface{}) BooleanExpression
		Lt(interface{}) BooleanExpression
		Lte(interface{}) BooleanExpression
	}
)

func (bo BooleanOperation) String() string {
	switch bo {
	case EqOp:
		return "eq"
	case NeqOp:
		return "neq"
	case IsOp:
		return "is"
	case IsNotOp:
		return "isnot"
	case GtOp:
		return "gt"
	case GteOp:
		return "gte"
	case LtOp:
		return "lt"
	case LteOp:
		return "lte"
	case RegexpLikeOp:
		return "regexplike"
	case RegexpNotLikeOp:
		return "regexpnotlike"
	case RegexpILikeOp:
		return "regexpilike"
	case RegexpNotILikeOp:
		return "regexpnotilike"
	}

	return fmt.Sprintf("%d", bo)
}

type SQLExpression interface {
	Expression
	ToSQL() (string, []interface{}, error)
}

type (
	RangeOperation  int
	RangeExpression interface {
		Expression
		Op() RangeOperation
		LHS() Expression
		RHS() RangeVal
	}
	RangeVal interface {
		Start() interface{}
		End() interface{}
	}
	Rangeable interface {
		Between(RangeVal) RangeExpression
		NotBetween(RangeVal) RangeExpression
	}
)

type (
	Aliaseable interface {
		As(interface{}) AliasedExpression
	}
	AliasedExpression interface {
		Expression
		Aliased() Expression
		GetAs() IdentifierExpression
	}
)

type (
	SortDirection     int
	OrderedExpression interface {
		Expression
		SortExpression() Expression
		IsAsc() bool
	}
	Orderable interface {
		Asc() OrderedExpression
		Desc() OrderedExpression
	}
)

type (
	Operator           int
	ComputerExpression interface {
		Expression
		Aliaseable
		Computable
		Op() Operator
		LHS() Expression
		RHS() interface{}
	}
	Computable interface {
		Add(val interface{}) ComputerExpression
		Sub(val interface{}) ComputerExpression
		Mul(val interface{}) ComputerExpression
	}
)
