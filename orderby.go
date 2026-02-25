package squildx

func (b *builder) OrderBy(expr string, params ...Params) Builder {
	cp := b.clone()
	p, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	parsed, prefix, err := parseParams(expr, p)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.orderBys = append(cp.orderBys, paramClause{sql: expr, params: parsed})
	return cp
}
