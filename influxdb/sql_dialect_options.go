package influxdb

type SQLFragmentType int

type SQLDialectOptions struct {
	// EscapedRunes is a map of a rune and the corresponding escape sequence in bytes. Used when escaping text
	// types.
	// (Default= map[rune][]byte{
	// 		'\'': []byte("''"),
	// 	})
	EscapedRunes map[rune][]byte

	ComputeOperatorLookup map[Operator][]byte
	// A map used to look up BooleanOperations and their SQL equivalents
	// (Default= map[exp.BooleanOperation][]byte{
	// 		exp.EqOp:             []byte("="),
	// 		exp.NeqOp:            []byte("!="),
	// 		exp.GtOp:             []byte(">"),
	// 		exp.GteOp:            []byte(">="),
	// 		exp.LtOp:             []byte("<"),
	// 		exp.LteOp:            []byte("<="),
	// 		exp.InOp:             []byte("IN"),
	// 		exp.NotInOp:          []byte("NOT IN"),
	// 		exp.IsOp:             []byte("IS"),
	// 		exp.IsNotOp:          []byte("IS NOT"),
	// 		exp.LikeOp:           []byte("LIKE"),
	// 		exp.NotLikeOp:        []byte("NOT LIKE"),
	// 		exp.ILikeOp:          []byte("ILIKE"),
	// 		exp.NotILikeOp:       []byte("NOT ILIKE"),
	// 		exp.RegexpLikeOp:     []byte("~"),
	// 		exp.RegexpNotLikeOp:  []byte("!~"),
	// 		exp.RegexpILikeOp:    []byte("~*"),
	// 		exp.RegexpNotILikeOp: []byte("!~*"),
	// })
	BooleanOperatorLookup map[BooleanOperation][]byte
	// Empty string (DEFAULT="")
	EmptyString string
	// The OR keyword used when joining ExpressionLists (DEFAULT=[]byte(" OR "))
	OrFragment []byte
	// The AND keyword used when joining ExpressionLists (DEFAULT=[]byte(" AND "))
	AndFragment []byte
	// The SQL LIMIT BY clause fragment(DEFAULT=[]byte(" LIMIT "))
	LimitFragment []byte
	// The SQL OFFSET BY clause fragment(DEFAULT=[]byte(" OFFSET "))
	OffsetFragment []byte
	// The SQL TZ BY clause fragment(DEFAULT=[]byte(" TZ"))
	TimezoneFragment []byte
	// The SQL AS fragment when aliasing an Expression(DEFAULT=[]byte(" AS "))
	AsFragment []byte
	// The ASC fragment when specifying column order (DEFAULT=[]byte(" ASC"))
	AscFragment []byte
	// The NULL literal to use when interpolating nulls values (DEFAULT=[]byte("NULL"))
	Null []byte
	// The DESC fragment when specifying column order (DEFAULT=[]byte(" DESC"))
	DescFragment []byte
	// The SQL PARTITION BY clause fragment(DEFAULT=[]byte(" PARTITION BY "))
	PartitionByFragment []byte
	// The SQL GROUP BY clause fragment(DEFAULT=[]byte(" GROUP BY "))
	GroupByFragment []byte
	// The SELECT fragment to use when generating sql. (DEFAULT=[]byte("SELECT"))
	SelectClause []byte

	// The order of SQL fragments when creating a SELECT statement
	// (Default=[]SQLFragmentType{
	// 		CommonTableSQLFragment,
	// 		SelectSQLFragment,
	// 		FromSQLFragment,
	// 		JoinSQLFragment,
	// 		WhereSQLFragment,
	// 		GroupBySQLFragment,
	// 		HavingSQLFragment,
	// 		CompoundsSQLFragment,
	// 		OrderSQLFragment,
	// 		LimitSQLFragment,
	// 		OffsetSQLFragment,
	// 		ForSQLFragment,
	// 	})
	SelectSQLOrder []SQLFragmentType
	// The SQL FROM clause fragment (DEFAULT=[]byte(" FROM"))
	FromFragment []byte
	// The SQL ORDER BY clause fragment(DEFAULT=[]byte(" ORDER BY "))
	OrderByFragment []byte
	// The placeholder fragment to use when generating a non interpolated statement (DEFAULT=[]byte"?")
	PlaceHolderFragment []byte
	// The SQL FILL clause fragment(DEFAULT=[]byte(" FILL"))
	FillFragment []byte
	// The SQL INTERVAL clause fragment(DEFAULT=[]byte(" INTERVAL"))
	IntervalFragment []byte
	// The SQL WHERE clause fragment (DEFAULT=[]byte(" WHERE "))
	WhereFragment []byte
	// The operator to use when setting values in an update statement (DEFAULT='=')
	SetOperatorRune rune
	// Left paren rune (DEFAULT='(')
	LeftParenRune rune
	// Right paren rune (DEFAULT=')')
	RightParenRune rune
	// Star rune (DEFAULT='*')
	StarRune rune
	// Period rune (DEFAULT='.')
	PeriodRune rune
	// Space rune (DEFAULT=' ')
	SpaceRune rune
	// Comma rune (DEFAULT=',')
	CommaRune rune
	// The quote rune to use when quoting identifiers(DEFAULT='"')
	QuoteRune rune
	// The quote rune to use when quoting string literals (DEFAULT='\'')
	StringQuote rune
}

const (
	SelectSQLFragment = iota
	FromSQLFragment
	WhereSQLFragment
	PartitionBySQLFragment
	GroupBySQLFragment
	IntervalFragment
	FillSQLFragment
	OrderSQLFragment
	LimitSQLFragment
	OffsetSQLFragment
	TimezoneSQLFragment
)

func DefaultDialectOptions() *SQLDialectOptions {
	return &SQLDialectOptions{
		SelectClause:        []byte("SELECT"),
		FromFragment:        []byte(" FROM"),
		WhereFragment:       []byte(" WHERE "),
		PartitionByFragment: []byte(" PARTITION BY "),
		GroupByFragment:     []byte(" GROUP BY "),
		IntervalFragment:    []byte(" INTERVAL"),
		FillFragment:        []byte(" FILL"),
		OrderByFragment:     []byte(" ORDER BY "),
		LimitFragment:       []byte(" LIMIT "),
		OffsetFragment:      []byte(" OFFSET "),
		TimezoneFragment:    []byte(" TZ"),
		AsFragment:          []byte(" AS "),
		AscFragment:         []byte(" ASC"),
		Null:                []byte("NULL"),
		DescFragment:        []byte(" DESC"),
		AndFragment:         []byte(" AND "),
		OrFragment:          []byte(" OR "),
		StringQuote:         '\'',
		SetOperatorRune:     '=',
		QuoteRune:           '"',
		PlaceHolderFragment: []byte("?"),
		EmptyString:         "",
		CommaRune:           ',',
		SpaceRune:           ' ',
		LeftParenRune:       '(',
		RightParenRune:      ')',
		StarRune:            '*',
		PeriodRune:          '.',
		BooleanOperatorLookup: map[BooleanOperation][]byte{
			EqOp:             []byte("="),
			NeqOp:            []byte("!="),
			GtOp:             []byte(">"),
			GteOp:            []byte(">="),
			LtOp:             []byte("<"),
			LteOp:            []byte("<="),
			RegexpLikeOp:     []byte("~"),
			RegexpNotLikeOp:  []byte("!~"),
			RegexpILikeOp:    []byte("~*"),
			RegexpNotILikeOp: []byte("!~*"),
			InOp:             []byte("IN"),
		},
		ComputeOperatorLookup: map[Operator][]byte{
			Plus:  []byte("+"),
			Minus: []byte("-"),
			Multi: []byte("*"),
		},
		EscapedRunes: map[rune][]byte{
			'\'': []byte("''"),
		},
		SelectSQLOrder: []SQLFragmentType{
			SelectSQLFragment,
			FromSQLFragment,
			WhereSQLFragment,
			PartitionBySQLFragment,
			GroupBySQLFragment,
			IntervalFragment,
			FillSQLFragment,
			OrderSQLFragment,
			LimitSQLFragment,
			OffsetSQLFragment,
			TimezoneSQLFragment,
		},
	}
}

var defaultDialectOptions = DefaultDialectOptions()
