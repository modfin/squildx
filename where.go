package squildx

func (b *builder) Where(sql string, values ...any) Builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.wheres = append(cp.wheres, paramClause{sql: sql, params: params})
	return cp
}
