package influxdb

import (
	"reflect"
)

type boolean struct {
	lhs Expression
	rhs interface{}
	op  BooleanOperation
}

func NewBooleanExpression(op BooleanOperation, lhs Expression, rhs interface{}) BooleanExpression {
	return boolean{
		lhs: lhs,
		rhs: rhs,
		op:  op,
	}
}

func (b boolean) Clone() Expression {
	return NewBooleanExpression(b.op, b.lhs.Clone(), b.rhs)
}

func (b boolean) Expression() Expression {
	return b
}

func (b boolean) RHS() interface{} {
	return b.rhs
}

func (b boolean) LHS() Expression {
	return b.lhs
}

func (b boolean) Op() BooleanOperation {
	return b.op
}

func eq(lhs Expression, rhs interface{}) BooleanExpression {
	return checkBoolExpType(EqOp, lhs, rhs, false)
}

func neq(lhs Expression, rhs interface{}) BooleanExpression {
	return checkBoolExpType(EqOp, lhs, rhs, true)
}

func gt(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(GtOp, lhs, rhs)
}

func gte(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(GteOp, lhs, rhs)
}

func lt(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(LtOp, lhs, rhs)
}

func lte(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(LteOp, lhs, rhs)
}

func checkBoolExpType(op BooleanOperation, lhs Expression, rhs interface{}, invert bool) BooleanExpression {
	if rhs == nil {
		op = IsOp
	} else {
		switch reflect.Indirect(reflect.ValueOf(rhs)).Kind() {
		case reflect.Bool:
			op = IsOp
		case reflect.Slice:
			if _, ok := rhs.([]byte); !ok {
				op = InOp
			}
		case reflect.Struct:
			switch rhs.(type) {
			case SQLExpression:
				op = InOp
			}
		default:
		}
	}

	if invert {
		op = operatorInversions[op]
	}

	return NewBooleanExpression(op, lhs, rhs)
}
