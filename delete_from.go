package squildx

func (b *deleteBuilder) From(table string) DeleteBuilder {
	cp := b.clone()
	cp.table = table
	return cp
}
