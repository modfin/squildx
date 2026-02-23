package squildx

type whereClause struct {
	sql    string
	params map[string]any
}

func (b *builder) Where(sql string, values ...any) Builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.wheres = append(cp.wheres, whereClause{sql: sql, params: params})
	return cp
}
