package squildx

func (b *updateBuilder) Returning(columns ...string) UpdateBuilder {
	cp := b.clone()
	cp.returnings = append(cp.returnings, columns...)
	return cp
}

func (b *updateBuilder) ReturningObject(obj any) UpdateBuilder {
	cp := b.clone()
	cols, err := structColumns(obj, "")
	if err != nil {
		cp.err = err
		return cp
	}
	cp.returnings = append(cp.returnings, cols...)
	return cp
}
