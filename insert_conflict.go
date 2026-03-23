package squildx

func (b *insertBuilder) OnConflictDoNothing(columns ...string) InsertBuilder {
	cp := b.clone()
	cp.conflict = &conflictClause{
		columns: columns,
	}
	return cp
}

func (b *insertBuilder) OnConflictDoUpdate(columns []string, set string, params ...Params) InsertBuilder {
	cp := b.clone()
	extracted, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	merged, prefix, err := parseParams(set, extracted)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.conflict = &conflictClause{
		columns:  columns,
		doUpdate: true,
		set:      set,
		params:   merged,
	}
	return cp
}
