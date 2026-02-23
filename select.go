package squildx

func (b *builder) Select(columns ...string) Builder {
	cp := b.clone()
	cp.columns = columns
	return cp
}
