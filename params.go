package squildx

import (
	"fmt"
	"regexp"
)

var paramRegex = regexp.MustCompile(`[:@][a-zA-Z_][a-zA-Z0-9_]*`)

// extractParams merges the variadic Params slices into a single map.
// Duplicate keys with different values produce ErrDuplicateParam.
func extractParams(maps []Params) (Params, error) {
	switch len(maps) {
	case 0:
		return nil, nil
	case 1:
		return maps[0], nil
	}
	merged := make(Params)
	for _, m := range maps {
		if err := mergeParams(merged, m); err != nil {
			return nil, err
		}
	}
	return merged, nil
}

// parseParams extracts named placeholders from sql, validates them against the
// provided params map, and returns the validated params plus the detected prefix
// byte (':' or '@', or 0 if no placeholders were found).
//
// It skips doubled-prefix sequences (:: and @@) so that PostgreSQL type casts
// (value::integer) and session variables (@@var) are not treated as parameters.
func parseParams(sql string, params Params) (Params, byte, error) {
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

	// Defensive copy so the builder never holds a reference to the caller's map.
	copied := make(Params, len(params))
	for k, v := range params {
		copied[k] = v
	}
	return copied, prefix, nil
}

func mergeParams(dst, src Params) error {
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
