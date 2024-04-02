package influxdb

type SelectClauses interface {
	IsDefaultSelect() bool

	Select() ColumnListExpression
	SetSelect(cl ColumnListExpression) SelectClauses

	Where() ExpressionList
	WhereAppend(expressions ...Expression) SelectClauses

	Order() ColumnListExpression
	SetOrder(oes ...OrderedExpression) SelectClauses

	GroupBy() ColumnListExpression
	GroupByAppend(cl ColumnListExpression) SelectClauses
	SetGroupBy(cl ColumnListExpression) SelectClauses

	PartitionBy() ColumnListExpression
	SetPartitionBy(cl ColumnListExpression) SelectClauses

	Interval() string
	SetInterval(interval string) SelectClauses

	Fill() interface{}
	SetFill(fill interface{}) SelectClauses

	Limit() interface{}
	ClearLimit() SelectClauses
	SetLimit(limit interface{}) SelectClauses

	Offset() uint
	SetOffset(offset uint) SelectClauses

	Distinct() ColumnListExpression
	SetDistinct(cle ColumnListExpression) SelectClauses

	From() ColumnListExpression
	SetFrom(cl ColumnListExpression) SelectClauses

	Timezone() string
	SetTimezone(tz string) SelectClauses

	Clone() *selectClauses
	Clear()
}

type selectClauses struct {
	selectColumns ColumnListExpression
	distinct      ColumnListExpression
	from          ColumnListExpression
	where         ExpressionList
	order         ColumnListExpression
	partitionBy   ColumnListExpression
	groupBy       ColumnListExpression
	fill          interface{}
	limit         interface{}
	interval      string
	timezone      string
	offset        uint
}

func newSelectClauses() SelectClauses {
	return &selectClauses{
		selectColumns: newColumnListExpression(Star()),
	}
}

func (sc *selectClauses) IsDefaultSelect() bool {
	ret := false

	if sc.selectColumns != nil {
		selects := sc.selectColumns.Columns()
		if len(selects) == 1 {
			if l, ok := selects[0].(LiteralExpression); ok && l.Literal() == "*" {
				ret = true
			}
		}
	}

	return ret
}

func (sc *selectClauses) Select() ColumnListExpression {
	return sc.selectColumns
}

func (sc *selectClauses) SetSelect(cl ColumnListExpression) SelectClauses {
	sc.selectColumns = cl

	return sc
}

func (sc *selectClauses) Where() ExpressionList {
	return sc.where
}

func (sc *selectClauses) WhereAppend(expressions ...Expression) SelectClauses {
	if len(expressions) == 0 {
		return sc
	}

	if sc.where == nil {
		sc.where = NewExpressionList(AndType, expressions...)
	} else {
		sc.where = sc.where.Append(expressions...)
	}

	return sc
}

func (sc *selectClauses) Order() ColumnListExpression {
	return sc.order
}

func (sc *selectClauses) SetOrder(oes ...OrderedExpression) SelectClauses {
	sc.order = newOrderedColumnList(oes...)

	return sc
}

func (sc *selectClauses) GroupBy() ColumnListExpression {
	return sc.groupBy
}

func (sc *selectClauses) GroupByAppend(cl ColumnListExpression) SelectClauses {
	if sc.groupBy == nil {
		return sc.SetGroupBy(cl)
	}

	sc.groupBy = sc.groupBy.Append(cl.Columns()...)

	return sc
}

func (sc *selectClauses) SetGroupBy(cl ColumnListExpression) SelectClauses {
	sc.groupBy = cl

	return sc
}

func (sc *selectClauses) PartitionBy() ColumnListExpression {
	return sc.partitionBy
}

func (sc *selectClauses) SetPartitionBy(cl ColumnListExpression) SelectClauses {
	sc.partitionBy = cl

	return sc
}

func (sc *selectClauses) Interval() string {
	return sc.interval
}

func (sc *selectClauses) SetInterval(interval string) SelectClauses {
	sc.interval = interval

	return sc
}

func (sc *selectClauses) Fill() interface{} {
	return sc.fill
}

func (sc *selectClauses) SetFill(fill interface{}) SelectClauses {
	sc.fill = fill

	return sc
}

func (sc *selectClauses) Limit() interface{} {
	return sc.limit
}

func (sc *selectClauses) ClearLimit() SelectClauses {
	sc.limit = nil

	return sc
}

func (sc *selectClauses) SetLimit(limit interface{}) SelectClauses {
	sc.limit = limit

	return sc
}

func (sc *selectClauses) Offset() uint {
	return sc.offset
}

func (sc *selectClauses) SetOffset(offset uint) SelectClauses {
	sc.offset = offset

	return sc
}

func (sc *selectClauses) Distinct() ColumnListExpression {
	return sc.distinct
}

func (sc *selectClauses) SetDistinct(cle ColumnListExpression) SelectClauses {
	sc.distinct = cle

	return sc
}

func (sc *selectClauses) From() ColumnListExpression {
	return sc.from
}

func (sc *selectClauses) SetFrom(cl ColumnListExpression) SelectClauses {
	sc.from = cl

	return sc
}

func (sc *selectClauses) Timezone() string {
	return sc.timezone
}

func (sc *selectClauses) SetTimezone(tz string) SelectClauses {
	sc.timezone = tz

	return sc
}

func (sc *selectClauses) Clone() *selectClauses {
	return &selectClauses{
		selectColumns: sc.selectColumns,
		distinct:      sc.distinct,
		from:          sc.from,
		where:         sc.where,
		interval:      sc.interval,
		order:         sc.order,
		groupBy:       sc.groupBy,
		fill:          sc.fill,
		limit:         sc.limit,
		offset:        sc.offset,
		timezone:      sc.timezone,
	}
}

func (sc *selectClauses) Clear() {
	sc.ClearLimit()
	sc.SetSelect(newColumnListExpression(Star())).SetDistinct(nil)
	sc.from = nil
	sc.where = nil
	sc.interval = ""
	sc.partitionBy = nil
	sc.order = nil
	sc.groupBy = nil
	sc.fill = nil
	sc.offset = 0
	sc.timezone = ""
}
