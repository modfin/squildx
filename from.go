package squildx

func (b *builder) From(table string) Builder {
	cp := b.clone()
	cp.from = table
	return cp
}
