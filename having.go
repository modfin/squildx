package squildx

func (b *builder) Having(sql string, params ...Params) Builder {
	cp := b.clone()
	p, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	parsed, prefix, err := parseParams(sql, p)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.havings = append(cp.havings, paramClause{sql: sql, params: parsed})
	return cp
}
