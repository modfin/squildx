package squildx

func (b *builder) Distinct() Builder {
	cp := b.clone()
	cp.distinct = true
	return cp
}
