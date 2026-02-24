package squildx

import (
	"fmt"
	"reflect"
	"regexp"
)

var paramRegex = regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)

func parseParams(sql string, values []any) (map[string]any, error) {
	matches := paramRegex.FindAllString(sql, -1)

	seen := make(map[string]struct{}, len(matches))
	unique := make([]string, 0, len(matches))
	for _, match := range matches {
		name := match[1:]
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			unique = append(unique, name)
		}
	}

	if len(unique) != len(values) {
		return nil, fmt.Errorf("%w: got %d placeholder(s) but %d value(s)", ErrParamMismatch, len(unique), len(values))
	}

	params := make(map[string]any, len(unique))
	for i, name := range unique {
		params[name] = values[i]
	}
	return params, nil
}

func mergeParams(dst, src map[string]any) error {
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			if !reflect.DeepEqual(existing, v) {
				return fmt.Errorf("%w: %q", ErrDuplicateParam, k)
			}
		}
		dst[k] = v
	}
	return nil
}
