package squildx

func (b *builder) Having(sql string, params ...map[string]any) Builder {
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
	if prefix != 0 {
		if cp.paramPrefix != 0 && cp.paramPrefix != prefix {
			cp.err = ErrMixedPrefix
			return cp
		}
		cp.paramPrefix = prefix
	}
	cp.havings = append(cp.havings, paramClause{sql: sql, params: parsed})
	return cp
}
