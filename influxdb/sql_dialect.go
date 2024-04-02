package influxdb

import (
	"fmt"
	"strconv"
)

type SQLDialect interface {
	ToSelectSQL(sb SQLBuilder, clauses SelectClauses)
}

type sqlDialect struct {
	selectGen SelectSQLGenerator
}

func newDialect(do *SQLDialectOptions) SQLDialect {
	return &sqlDialect{
		selectGen: newSelectSQLGenerator(do),
	}
}

func (sd *sqlDialect) ToSelectSQL(sb SQLBuilder, clauses SelectClauses) {
	sd.selectGen.Generate(sb, clauses)
}

type SelectSQLGenerator interface {
	Generate(sb SQLBuilder, clauses SelectClauses)
}

type selectSQLGenerator struct {
	CommonSQLGenerator
}

func newSelectSQLGenerator(do *SQLDialectOptions) SelectSQLGenerator {
	return &selectSQLGenerator{newCommonSQLGenerator(do)}
}

func ErrNotSupportedFragment(sqlType string, f SQLFragmentType) error {
	return fmt.Errorf("unsupported %s SQL fragment %v", sqlType, f)
}

func (ssg *selectSQLGenerator) Generate(sb SQLBuilder, clauses SelectClauses) {
	for _, f := range ssg.DialectOptions().SelectSQLOrder {
		if sb.Error() != nil {
			return
		}

		switch f {
		case SelectSQLFragment:
			ssg.SelectSQL(sb, clauses)
		case FromSQLFragment:
			ssg.FromSQL(sb, clauses.From())
		case WhereSQLFragment:
			ssg.WhereSQL(sb, clauses.Where())
		case PartitionBySQLFragment:
			ssg.PartitionBySQL(sb, clauses.PartitionBy())
		case GroupBySQLFragment:
			ssg.GroupBySQL(sb, clauses.GroupBy())
		case IntervalFragment:
			ssg.IntervalSQL(sb, clauses.Interval())
		case FillSQLFragment:
			ssg.FillSQL(sb, clauses.Fill())
		case OrderSQLFragment:
			ssg.OrderSQL(sb, clauses.Order())
		case LimitSQLFragment:
			ssg.LimitSQL(sb, clauses.Limit())
		case OffsetSQLFragment:
			ssg.OffsetSQL(sb, clauses.Offset())
		case TimezoneSQLFragment:
			// ssg.TimezoneSQL(sb, clauses.Timezone())
		default:
			sb.SetError(ErrNotSupportedFragment("SELECT", f))
		}
	}
}

func (ssg *selectSQLGenerator) selectSQLCommon(sb SQLBuilder, clauses SelectClauses) {
	if cols := clauses.Select(); clauses.IsDefaultSelect() || len(cols.Columns()) == 0 {
		sb.WriteRunes(ssg.DialectOptions().StarRune)
	} else {
		ssg.ExpressionSQLGenerator().Generate(sb, cols)
	}
}

func (ssg *selectSQLGenerator) SelectSQL(sb SQLBuilder, clauses SelectClauses) {
	sb.Write(ssg.DialectOptions().SelectClause).WriteRunes(ssg.DialectOptions().SpaceRune)
	ssg.selectSQLCommon(sb, clauses)
}

func (ssg *selectSQLGenerator) PartitionBySQL(sb SQLBuilder, partitionBy ColumnListExpression) {
	if partitionBy != nil && len(partitionBy.Columns()) > 0 {
		sb.Write(ssg.DialectOptions().PartitionByFragment)
		ssg.ExpressionSQLGenerator().Generate(sb, partitionBy)
	}
}

func (ssg *selectSQLGenerator) GroupBySQL(sb SQLBuilder, groupBy ColumnListExpression) {
	if groupBy != nil && len(groupBy.Columns()) > 0 {
		sb.Write(ssg.DialectOptions().GroupByFragment)
		ssg.ExpressionSQLGenerator().Generate(sb, groupBy)
	}
}

func (ssg *selectSQLGenerator) FillSQL(sb SQLBuilder, fill interface{}) {
	if fill != nil {
		sb.Write(ssg.DialectOptions().FillFragment)
		sb.WriteRunes(ssg.DialectOptions().LeftParenRune)

		switch v := fill.(type) {
		case string:
			sb.WriteStrings(v)
		case int:
			sb.WriteStrings("VALUE, ")
			sb.WriteStrings(strconv.Itoa(v))
		case float64:
			sb.WriteStrings("VALUE, ")
			sb.WriteStrings(strconv.FormatFloat(v, 'b', 10, 64))
		default:
			ssg.ExpressionSQLGenerator().Generate(sb, fill)
		}
		sb.WriteRunes(ssg.DialectOptions().RightParenRune)
	}
}

func (ssg *selectSQLGenerator) IntervalSQL(sb SQLBuilder, interval string) {
	if interval != "" {
		sb.Write(ssg.DialectOptions().IntervalFragment)
		sb.WriteRunes(ssg.DialectOptions().LeftParenRune)

		sb.WriteStrings(interval)
		sb.WriteRunes(ssg.DialectOptions().RightParenRune)
	}
}

func (ssg *selectSQLGenerator) OffsetSQL(sb SQLBuilder, offset uint) {
	if offset > 0 {
		sb.Write(ssg.DialectOptions().OffsetFragment)
		ssg.ExpressionSQLGenerator().Generate(sb, offset)
	}
}

type CommonSQLGenerator interface {
	DialectOptions() *SQLDialectOptions
	ExpressionSQLGenerator() ExpressionSQLGenerator
	FromSQL(sb SQLBuilder, from ColumnListExpression)
	WhereSQL(sb SQLBuilder, where ExpressionList)
	OrderSQL(sb SQLBuilder, order ColumnListExpression)
	SourcesSQL(sb SQLBuilder, from ColumnListExpression)
	LimitSQL(sb SQLBuilder, limit interface{})
	TimezoneSQL(sb SQLBuilder, tz string)
}

type commonSQLGenerator struct {
	esg            ExpressionSQLGenerator
	dialectOptions *SQLDialectOptions
}

func newCommonSQLGenerator(do *SQLDialectOptions) CommonSQLGenerator {
	return &commonSQLGenerator{
		esg:            newExpressionSQLGenerator(do),
		dialectOptions: do,
	}
}

func (csg *commonSQLGenerator) DialectOptions() *SQLDialectOptions {
	return csg.dialectOptions
}

func (csg *commonSQLGenerator) ExpressionSQLGenerator() ExpressionSQLGenerator {
	return csg.esg
}

func (csg *commonSQLGenerator) FromSQL(sb SQLBuilder, from ColumnListExpression) {
	if from != nil && !from.IsEmpty() {
		sb.Write(csg.dialectOptions.FromFragment)
		csg.SourcesSQL(sb, from)
	}
}

func (csg *commonSQLGenerator) WhereSQL(sb SQLBuilder, where ExpressionList) {
	if where != nil && !where.IsEmpty() {
		sb.Write(csg.dialectOptions.WhereFragment)
		csg.esg.Generate(sb, where)
	}
}

func (csg *commonSQLGenerator) OrderSQL(sb SQLBuilder, order ColumnListExpression) {
	if order != nil && len(order.Columns()) > 0 {
		sb.Write(csg.dialectOptions.OrderByFragment)
		csg.esg.Generate(sb, order)
	}
}

func (csg *commonSQLGenerator) SourcesSQL(sb SQLBuilder, from ColumnListExpression) {
	sb.WriteRunes(csg.dialectOptions.SpaceRune)
	csg.esg.Generate(sb, from)
}

func (csg *commonSQLGenerator) LimitSQL(sb SQLBuilder, limit interface{}) {
	if limit != nil {
		sb.Write(csg.dialectOptions.LimitFragment)
		csg.esg.Generate(sb, limit)
	}
}

func (csg *commonSQLGenerator) TimezoneSQL(sb SQLBuilder, tz string) {
	if tz != "" {
		sb.Write(csg.dialectOptions.TimezoneFragment)
		sb.WriteRunes(csg.dialectOptions.LeftParenRune)
		csg.esg.Generate(sb, tz)
		sb.WriteRunes(csg.dialectOptions.RightParenRune)
	}
}
