package squildx

import (
	"fmt"
	"regexp"
)

var paramRegex = regexp.MustCompile(`[:@][a-zA-Z_][a-zA-Z0-9_]*`)

// extractParams validates the variadic maps slice and returns a single map.
// 0 maps → nil, nil; 1 map → that map, nil; 2+ maps → nil, ErrMultipleParamMaps.
func extractParams(maps []map[string]any) (map[string]any, error) {
	switch len(maps) {
	case 0:
		return nil, nil
	case 1:
		return maps[0], nil
	default:
		return nil, ErrMultipleParamMaps
	}
}

// parseParams extracts named placeholders from sql, validates them against the
// provided params map, and returns the validated params plus the detected prefix
// byte (':' or '@', or 0 if no placeholders were found).
//
// It skips doubled-prefix sequences (:: and @@) so that PostgreSQL type casts
// (value::integer) and session variables (@@var) are not treated as parameters.
func parseParams(sql string, params map[string]any) (map[string]any, byte, error) {
	indices := paramRegex.FindAllStringIndex(sql, -1)

	var prefix byte
	placeholders := make(map[string]struct{})
	for _, idx := range indices {
		// Skip doubled-prefix: e.g. :: in "value::integer" or @@ in "@@session_var"
		if idx[0] > 0 && sql[idx[0]-1] == sql[idx[0]] {
			continue
		}

		p := sql[idx[0]]
		if prefix == 0 {
			prefix = p
		}
		if p != prefix {
			return nil, 0, ErrMixedPrefix
		}

		name := sql[idx[0]+1 : idx[1]]
		placeholders[name] = struct{}{}
	}

	if len(placeholders) == 0 && len(params) == 0 {
		return nil, 0, nil
	}

	// Validate: every placeholder has a matching map key
	for name := range placeholders {
		if _, ok := params[name]; !ok {
			return nil, 0, fmt.Errorf("%w: %q", ErrMissingParam, name)
		}
	}

	// Validate: every map key has a matching placeholder
	for key := range params {
		if _, ok := placeholders[key]; !ok {
			return nil, 0, fmt.Errorf("%w: %q", ErrExtraParam, key)
		}
	}

	return params, prefix, nil
}

func mergeParams(dst, src map[string]any) error {
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			if !valueEqual(existing, v) {
				return fmt.Errorf("%w: %q", ErrDuplicateParam, k)
			}
		}
		dst[k] = v
	}
	return nil
}
