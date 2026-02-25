package squildx

import "errors"

var (
	ErrNoColumns            = errors.New("squildx: SELECT requires at least one column")
	ErrNoFrom               = errors.New("squildx: SELECT requires a FROM clause")
	ErrDuplicateParam       = errors.New("squildx: duplicate parameter with conflicting value")
	ErrDuplicateJoin        = errors.New("squildx: duplicate join with conflicting clause")
	ErrMissingParam         = errors.New("squildx: placeholder has no matching value in params map")
	ErrExtraParam           = errors.New("squildx: params map key has no matching placeholder")
	ErrMixedPrefix          = errors.New("squildx: mixed parameter prefixes (: and @) in the same query")
	ErrHavingWithoutGroupBy = errors.New("squildx: HAVING requires a GROUP BY clause")
	ErrNotAStruct           = errors.New("squildx: SelectObject requires a struct or pointer to struct")
)
