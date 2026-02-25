package squildx

import (
	"reflect"
	"unicode"
)

func copySlice[T any](s []T) []T {
	if s == nil {
		return nil
	}
	cp := make([]T, len(s))
	copy(cp, s)
	return cp
}

func valueEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func paramsEqual(a, b Params) bool {
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		vb, ok := b[k]
		if !ok || !valueEqual(va, vb) {
			return false
		}
	}
	return true
}

func buildersEqual(a, b Builder) bool {
	sqlA, paramsA, errA := a.Build()
	sqlB, paramsB, errB := b.Build()
	if errA != nil || errB != nil {
		return false
	}
	return sqlA == sqlB && paramsEqual(paramsA, paramsB)
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) {
					// camelCase boundary: aB -> a_b
					result = append(result, '_')
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					// acronym boundary: HTTPCode -> http_code (insert _ before last capital of acronym)
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
