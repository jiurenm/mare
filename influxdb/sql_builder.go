package influxdb

import "bytes"

type SQLBuilder interface {
	ToSQL() (sql string, err error)
	Write(p []byte) SQLBuilder
	WriteStrings(ss ...string) SQLBuilder
	WriteRunes(r ...rune) SQLBuilder
	SetError(err error) SQLBuilder
	Error() error
}

type sqlBuilder struct {
	err       error
	buf       *bytes.Buffer
	isPrepare bool
}

func newSQLBuilder(isPrepared bool) SQLBuilder {
	return &sqlBuilder{
		buf:       &bytes.Buffer{},
		isPrepare: isPrepared,
	}
}

func (sb *sqlBuilder) ToSQL() (string, error) {
	if sb.err != nil {
		return "", sb.err
	}

	return sb.buf.String(), nil
}

func (sb *sqlBuilder) Write(p []byte) SQLBuilder {
	if sb.err == nil {
		sb.buf.Write(p)
	}

	return sb
}

func (sb *sqlBuilder) WriteStrings(ss ...string) SQLBuilder {
	if sb.err == nil {
		for _, s := range ss {
			sb.buf.WriteString(s)
		}
	}

	return sb
}

func (sb *sqlBuilder) WriteRunes(rs ...rune) SQLBuilder {
	if sb.err == nil {
		for _, r := range rs {
			sb.buf.WriteRune(r)
		}
	}

	return sb
}

func (sb *sqlBuilder) Error() error {
	return sb.err
}

func (sb *sqlBuilder) SetError(err error) SQLBuilder {
	if sb.err == nil {
		sb.err = err
	}

	return sb
}
