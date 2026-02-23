package squildx

import "errors"

var (
	ErrNoColumns            = errors.New("squildx: SELECT requires at least one column")
	ErrNoFrom               = errors.New("squildx: SELECT requires a FROM clause")
	ErrDuplicateParam       = errors.New("squildx: duplicate parameter with conflicting value")
	ErrParamMismatch        = errors.New("squildx: number of :name placeholders does not match number of values")
	ErrHavingWithoutGroupBy = errors.New("squildx: HAVING requires a GROUP BY clause")
	ErrNotAStruct           = errors.New("squildx: SelectObject requires a struct or pointer to struct")
)
