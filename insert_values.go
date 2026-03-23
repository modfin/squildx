package squildx

import (
	"slices"
	"strings"
)

func (b *insertBuilder) Values(sql string, params ...Params) InsertBuilder {
	cp := b.clone()
	extracted, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	merged, prefix, err := parseParams(sql, extracted)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.valueRows = append(cp.valueRows, paramClause{sql: sql, params: merged})
	return cp
}

func (b *insertBuilder) ValuesObject(obj any) InsertBuilder {
	cp := b.clone()
	cols, sql, params, err := structFieldValues(obj)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(':'); err != nil {
		cp.err = err
		return cp
	}
	switch {
	case len(cp.columns) == 0:
		cp.columns = cols
	case !slices.Equal(cp.columns, cols):
		cp.err = ErrColumnMismatch
		return cp
	}
	cp.valueRows = append(cp.valueRows, paramClause{sql: sql, params: params})
	return cp
}

func structFieldValues(obj any) (columns []string, sql string, params Params, err error) {
	columns, err = structColumns(obj, "")
	if err != nil {
		return nil, "", nil, err
	}

	params, err = structValues(obj)
	if err != nil {
		return nil, "", nil, err
	}

	placeholders := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = ":" + col
	}
	sql = strings.Join(placeholders, ", ")
	return columns, sql, params, nil
}
