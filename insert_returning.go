package squildx

func (b *insertBuilder) Returning(columns ...string) InsertBuilder {
	cp := b.clone()
	cp.returnings = append(cp.returnings, columns...)
	return cp
}

func (b *insertBuilder) ReturningObject(obj any) InsertBuilder {
	cp := b.clone()
	cols, err := structColumns(obj, "")
	if err != nil {
		cp.err = err
		return cp
	}
	cp.returnings = append(cp.returnings, cols...)
	return cp
}
