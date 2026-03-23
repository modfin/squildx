package squildx

func (b *updateBuilder) Table(table string) UpdateBuilder {
	cp := b.clone()
	cp.table = table
	return cp
}
