package squildx

func (b *builder) OrderBy(exprs ...string) Builder {
	cp := b.clone()
	cp.orderBys = append(cp.orderBys, exprs...)
	return cp
}
