package influxdb

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"
)

type QueryBuilder struct {
	dialect SQLDialect
	clauses SelectClauses
	err     error
}

type UnionBuilder struct {
	query1 *QueryBuilder
	query2 *QueryBuilder
	tz     string
}

func newQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		dialect: newDialect(defaultDialectOptions),
		clauses: newSelectClauses(),
	}
}

var queryBuilderPool = &sync.Pool{
	New: func() interface{} {
		return newQueryBuilder()
	},
}

func From(table string) *QueryBuilder {
	return queryBuilderPool.Get().(*QueryBuilder).From(table)
}

func UnionAll(query1, query2 *QueryBuilder) *UnionBuilder {
	return &UnionBuilder{
		query1: query1,
		query2: query2,
	}
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.clauses.SetFrom(newColumnListExpression(table))

	return qb
}

func (qb *QueryBuilder) Select(selects ...interface{}) *QueryBuilder {
	if len(selects) == 0 {
		return qb.ClearSelect()
	}

	qb.clauses.SetSelect(newColumnListExpression(selects...))

	return qb
}

func (qb *QueryBuilder) ClearSelect() *QueryBuilder {
	qb.clauses.SetSelect(newColumnListExpression(Star())).SetDistinct(nil)

	return qb
}

func (qb *QueryBuilder) Where(expressions ...Expression) *QueryBuilder {
	qb.clauses.WhereAppend(expressions...)

	return qb
}

func (qb *QueryBuilder) Order(order ...OrderedExpression) *QueryBuilder {
	qb.clauses.SetOrder(order...)

	return qb
}

func (qb *QueryBuilder) PartitionBy(part ...interface{}) *QueryBuilder {
	qb.clauses.SetPartitionBy(newColumnListExpression(part...))

	return qb
}

func (qb *QueryBuilder) GroupBy(groupBy ...interface{}) *QueryBuilder {
	qb.clauses.SetGroupBy(newColumnListExpression(groupBy...))

	return qb
}

func (qb *QueryBuilder) GroupByAppend(groupBy ...interface{}) *QueryBuilder {
	qb.clauses.GroupByAppend(newColumnListExpression(groupBy...))

	return qb
}

func (qb *QueryBuilder) Interval(interval string) *QueryBuilder {
	qb.clauses.SetInterval(interval)

	return qb
}

func (qb *QueryBuilder) Fill(fill interface{}) *QueryBuilder {
	qb.clauses.SetFill(fill)

	return qb
}

func (qb *QueryBuilder) Timezone(tz string) *QueryBuilder {
	qb.clauses.SetTimezone(tz)

	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	if limit > 0 {
		qb.clauses.SetLimit(limit)

		return qb
	}

	qb.clauses.ClearLimit()

	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.clauses.SetOffset(uint(offset))

	return qb
}

func (qb *QueryBuilder) Clone() *QueryBuilder {
	return &QueryBuilder{
		dialect: newDialect(defaultDialectOptions),
		clauses: qb.clauses.Clone(),
		err:     qb.err,
	}
}

func needOffset(interval string) bool {
	if strings.HasSuffix(interval, "d") {
		return true
	}

	if strings.HasSuffix(interval, "n") {
		return true
	}

	if strings.HasSuffix(interval, "y") {
		return true
	}

	return false
}

func (qb *QueryBuilder) ToSQL() (string, string, error) {
	tz := qb.clauses.Timezone()

	// d:日 n:月 y:年
	interval := qb.clauses.Interval()

	if tz != "" && interval != "" && needOffset(interval) {
		loc, _ := time.LoadLocation(tz)
		_, offset := time.Now().In(loc).Zone()

		offset = 8 - (offset / 3600)

		if offset > 0 {
			qb.clauses.SetInterval(interval + "," + strconv.Itoa(offset) + "h")
		} else if offset < 0 {
			offset = 24 + offset
		}
	}

	sql, err := qb.selectSQLBuilder().ToSQL()

	qb.clear()
	queryBuilderPool.Put(qb)

	return sql, tz, err
}

func (qb *QueryBuilder) Query(ctx context.Context, conn *InfluxDB, dest interface{}, format ...FormatType) error {
	sql, _, err := qb.ToSQL()
	if err != nil {
		return err
	}

	return conn.Query(ctx, sql, dest, format...)
}

func (qb *QueryBuilder) QueryTaos(ctx context.Context, conn *InfluxDB, dest interface{}) error {
	sql, tz, err := qb.ToSQL()
	if err != nil {
		return err
	}

	return conn.Query2(ctx, sql, dest, tz)
}

func (qb *QueryBuilder) clear() {
	qb.clauses.Clear()
}

func (qb *QueryBuilder) selectSQLBuilder() SQLBuilder {
	buf := newSQLBuilder(true)
	if qb.err != nil {
		return buf.SetError(qb.err)
	}

	qb.dialect.ToSelectSQL(buf, qb.clauses)

	return buf
}

func (ub *UnionBuilder) Timezone(tz string) *UnionBuilder {
	ub.tz = tz

	return ub
}

func (ub *UnionBuilder) Query(ctx context.Context, conn *InfluxDB, dest interface{}) error {
	sql1, _, err := ub.query1.ToSQL()
	if err != nil {
		return err
	}

	sql2, _, err := ub.query2.ToSQL()
	if err != nil {
		return err
	}

	sql := sql1 + " UNION ALL " + sql2

	return conn.Query2(ctx, sql, dest, ub.tz)
}
