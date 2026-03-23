package squildx

import "strings"

func (b *insertBuilder) Build() (string, Params, error) {
	if b.err != nil {
		return "", nil, b.err
	}

	if b.table == "" {
		return "", nil, ErrNoTable
	}
	if len(b.columns) == 0 {
		return "", nil, ErrNoInsertColumns
	}

	hasValues := len(b.valueRows) > 0
	hasSelect := b.selectQuery != nil
	switch {
	case hasValues && hasSelect:
		return "", nil, ErrValuesAndSelect
	case !hasValues && !hasSelect:
		return "", nil, ErrNoInsertValues
	}

	params := make(Params)
	prefix := b.paramPrefix

	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(b.table)
	sb.WriteString(" (")
	sb.WriteString(strings.Join(b.columns, ", "))
	sb.WriteString(")")

	switch {
	case hasValues:
		sb.WriteString(" VALUES ")
		for i, row := range b.valueRows {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("(")
			sb.WriteString(row.sql)
			sb.WriteString(")")
			if err := mergeParams(params, row.params); err != nil {
				return "", nil, err
			}
		}
	case hasSelect:
		subSQL, subParams, err := b.selectQuery.Build()
		if err != nil {
			return "", nil, err
		}
		subPrefix := detectPrefix(subSQL)
		var reconcileErr error
		prefix, reconcileErr = reconcilePrefix(prefix, subPrefix)
		if reconcileErr != nil {
			return "", nil, reconcileErr
		}
		sb.WriteString(" ")
		sb.WriteString(subSQL)
		if err := mergeParams(params, subParams); err != nil {
			return "", nil, err
		}
	}

	if b.conflict != nil {
		sb.WriteString(" ON CONFLICT (")
		sb.WriteString(strings.Join(b.conflict.columns, ", "))
		sb.WriteString(")")
		switch {
		case b.conflict.doUpdate:
			sb.WriteString(" DO UPDATE SET ")
			sb.WriteString(b.conflict.set)
			if err := mergeParams(params, b.conflict.params); err != nil {
				return "", nil, err
			}
		default:
			sb.WriteString(" DO NOTHING")
		}
	}

	if len(b.returnings) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(b.returnings, ", "))
	}

	_ = prefix // prefix tracked for future reconciliation with mixed queries

	return sb.String(), params, nil
}
