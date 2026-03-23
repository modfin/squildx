package squildx

func (b *insertBuilder) Into(table string) InsertBuilder {
	cp := b.clone()
	cp.table = table
	return cp
}
