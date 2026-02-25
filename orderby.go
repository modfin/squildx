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
	if prefix != 0 {
		if cp.paramPrefix != 0 && cp.paramPrefix != prefix {
			cp.err = ErrMixedPrefix
			return cp
		}
		cp.paramPrefix = prefix
	}
	cp.orderBys = append(cp.orderBys, paramClause{sql: expr, params: parsed})
	return cp
}
