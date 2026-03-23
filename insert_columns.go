package squildx

func (b *insertBuilder) Columns(columns ...string) InsertBuilder {
	cp := b.clone()
	cp.columns = append(cp.columns, columns...)
	return cp
}

func (b *insertBuilder) ColumnsObject(obj any) InsertBuilder {
	cp := b.clone()
	cols, err := structColumns(obj, "")
	if err != nil {
		cp.err = err
		return cp
	}
	cp.columns = append(cp.columns, cols...)
	return cp
}
