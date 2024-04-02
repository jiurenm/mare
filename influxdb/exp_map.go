package influxdb

import (
	"fmt"
	"sort"
	"strings"
)

type (
	Ex   map[string]interface{}
	ExOr map[string]interface{}
	Op   map[string]interface{}
)

func (e Ex) Expression() Expression {
	return e
}

func (e Ex) Clone() Expression {
	ret := Ex{}
	for key, val := range e {
		ret[key] = val
	}

	return ret
}

func (e Ex) IsEmpty() bool {
	return len(e) == 0
}

func (e Ex) ToExpressions() (ExpressionList, error) {
	return mapToExpressionList(e, AndType)
}

func (eo ExOr) Expression() Expression {
	return eo
}

func (eo ExOr) Clone() Expression {
	ret := ExOr{}
	for key, val := range eo {
		ret[key] = val
	}

	return ret
}

func (eo ExOr) IsEmpty() bool {
	return len(eo) == 0
}

func (eo ExOr) ToExpressions() (ExpressionList, error) {
	return mapToExpressionList(eo, OrType)
}

func getExMapKeys(ex map[string]interface{}) []string {
	keys := make([]string, 0, len(ex))
	for key := range ex {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func createOredExpressionFromMap(lhs IdentifierExpression, op Op) ([]Expression, error) {
	opKeys := getExMapKeys(op)
	ors := make([]Expression, 0, len(opKeys))

	for _, opKey := range opKeys {
		if exp, err := createExpressionFromOp(lhs, opKey, op); err != nil {
			return nil, err
		} else if exp != nil {
			ors = append(ors, exp)
		}
	}

	return ors, nil
}

func createExpressionFromOp(lhs IdentifierExpression, opKey string, op Op) (exp Expression, err error) {
	switch strings.ToLower(opKey) {
	case EqOp.String():
		exp = lhs.Eq(op[opKey])
	case NeqOp.String():
		exp = lhs.Neq(op[opKey])
	case GtOp.String():
		exp = lhs.Gt(op[opKey])
	case GteOp.String():
		exp = lhs.Gte(op[opKey])
	case LtOp.String():
		exp = lhs.Lt(op[opKey])
	case LteOp.String():
		exp = lhs.Lte(op[opKey])
	default:
		err = fmt.Errorf("unsupported expression type %s", opKey)
	}

	return exp, err
}

func mapToExpressionList(ex map[string]interface{}, eType ExpressionListType) (ExpressionList, error) {
	keys := getExMapKeys(ex)
	ret := make([]Expression, 0, len(keys))

	for _, key := range keys {
		lhs := ParseIdentifier(key)
		rhs := ex[key]

		var exp Expression

		if op, ok := rhs.(Op); ok {
			ors, err := createOredExpressionFromMap(lhs, op)
			if err != nil {
				return nil, err
			}

			exp = NewExpressionList(OrType, ors...)
		} else {
			exp = lhs.Eq(rhs)
		}

		ret = append(ret, exp)
	}

	if eType == OrType {
		return NewExpressionList(OrType, ret...), nil
	}

	return NewExpressionList(AndType, ret...), nil
}
