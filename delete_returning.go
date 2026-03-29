package squildx

func (b *deleteBuilder) Returning(columns ...string) DeleteBuilder {
	cp := b.clone()
	cp.returnings = append(cp.returnings, columns...)
	return cp
}

func (b *deleteBuilder) ReturningObject(obj any) DeleteBuilder {
	cp := b.clone()
	cols, err := structColumns(obj, "")
	if err != nil {
		cp.err = err
		return cp
	}
	cp.returnings = append(cp.returnings, cols...)
	return cp
}
