package squildx

func (b *builder) Offset(n uint64) Builder {
	cp := b.clone()
	cp.offset = &n
	return cp
}
