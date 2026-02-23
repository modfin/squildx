package squildx

func (b *builder) OrderBy(expr string, values ...any) Builder {
	cp := b.clone()
	params, err := parseParams(expr, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.orderBys = append(cp.orderBys, paramClause{sql: expr, params: params})
	return cp
}
