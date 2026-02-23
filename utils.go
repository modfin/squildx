package squildx

func copySlice[T any](s []T) []T {
	if s == nil {
		return nil
	}
	cp := make([]T, len(s))
	copy(cp, s)
	return cp
}
