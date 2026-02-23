package squildx

func (b *builder) Limit(n uint64) Builder {
	cp := b.clone()
	cp.limit = &n
	return cp
}
