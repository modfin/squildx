package squildx

func (b *builder) Having(sql string, values ...any) Builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.havings = append(cp.havings, paramClause{sql: sql, params: params})
	return cp
}
