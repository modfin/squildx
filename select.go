package squildx

func (b *builder) Select(columns ...string) Builder {
	cp := b.clone()
	cp.columns = append(cp.columns, columns...)
	return cp
}

func (b *builder) RemoveSelect(columns ...string) Builder {
	cp := b.clone()
	remove := make(map[string]struct{}, len(columns))
	for _, c := range columns {
		remove[c] = struct{}{}
	}
	filtered := cp.columns[:0]
	for _, c := range cp.columns {
		if _, ok := remove[c]; !ok {
			filtered = append(filtered, c)
		}
	}
	cp.columns = filtered
	return cp
}
