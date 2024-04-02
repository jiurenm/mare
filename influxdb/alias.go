package influxdb

import "fmt"

type aliasExpression struct {
	aliased Expression
	alias   IdentifierExpression
}

func NewAliasExpression(exp Expression, alias interface{}) AliasedExpression {
	switch v := alias.(type) {
	case string:
		return aliasExpression{aliased: exp, alias: ParseIdentifier(v)}
	case IdentifierExpression:
		return aliasExpression{aliased: exp, alias: v}
	default:
		panic(fmt.Sprintf("Cannot create alias from %+v", v))
	}
}

func (ae aliasExpression) Clone() Expression {
	return NewAliasExpression(ae.aliased, ae.alias.Clone())
}

func (ae aliasExpression) Expression() Expression {
	return ae
}

func (ae aliasExpression) Aliased() Expression {
	return ae.aliased
}

func (ae aliasExpression) GetAs() IdentifierExpression {
	return ae.alias
}
