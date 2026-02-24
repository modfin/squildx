package squildx

func (b *builder) Distinct() Builder {
	cp := b.clone()
	cp.distinct = true
	return cp
}

func (b *builder) DistinctOn(columns ...string) Builder {
	cp := b.clone()
	cp.distinctOn = append(cp.distinctOn, columns...)
	return cp
}
