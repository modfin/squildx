package squildx

func (b *builder) GroupBy(exprs ...string) Builder {
	cp := b.clone()
	cp.groupBys = append(cp.groupBys, exprs...)
	return cp
}
