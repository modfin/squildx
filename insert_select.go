package squildx

func (b *insertBuilder) Select(sub Builder) InsertBuilder {
	cp := b.clone()
	cp.selectQuery = sub
	return cp
}
