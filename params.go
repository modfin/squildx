package squildx

import (
	"fmt"
	"regexp"
)

var paramRegex = regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)

func parseParams(sql string, values []any) (map[string]any, error) {
	matches := paramRegex.FindAllString(sql, -1)

	if len(matches) != len(values) {
		return nil, fmt.Errorf("%w: got %d placeholder(s) but %d value(s)", ErrParamMismatch, len(matches), len(values))
	}

	params := make(map[string]any, len(matches))
	for i, match := range matches {
		name := match[1:]
		params[name] = values[i]
	}
	return params, nil
}

func mergeParams(dst, src map[string]any) error {
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			if existing != v {
				return fmt.Errorf("%w: %q", ErrDuplicateParam, k)
			}
		}
		dst[k] = v
	}
	return nil
}
