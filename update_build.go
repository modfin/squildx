package squildx

import (
	"fmt"
	"strings"
)

func (b *updateBuilder) Build() (string, Params, error) {
	if b.err != nil {
		return "", nil, b.err
	}

	if b.table == "" {
		return "", nil, ErrUpdateNoTable
	}
	if len(b.sets) == 0 {
		return "", nil, ErrUpdateNoSet
	}
	if len(b.wheres) == 0 {
		return "", nil, ErrUpdateNoWhere
	}

	params := make(Params)
	prefix := b.paramPrefix

	var sb strings.Builder

	sb.WriteString("UPDATE ")
	sb.WriteString(b.table)

	setClauses := make([]string, len(b.sets))
	for i, s := range b.sets {
		setPrefix := detectPrefix(s.sql)
		var err error
		prefix, err = reconcilePrefix(prefix, setPrefix)
		if err != nil {
			return "", nil, err
		}
		setClauses[i] = s.sql
		if err := mergeParams(params, s.params); err != nil {
			return "", nil, err
		}
	}
	sb.WriteString(" SET ")
	sb.WriteString(strings.Join(setClauses, ", "))

	ands := make([]string, len(b.wheres))
	for i, w := range b.wheres {
		if w.subQuery != nil {
			subSQL, subParams, err := w.subQuery.Build()
			if err != nil {
				return "", nil, err
			}
			subPrefix := detectPrefix(subSQL)
			prefix, err = reconcilePrefix(prefix, subPrefix)
			if err != nil {
				return "", nil, err
			}
			ands[i] = fmt.Sprintf("%s (%s)", w.subPrefix, subSQL)
			if err := mergeParams(params, subParams); err != nil {
				return "", nil, err
			}
			continue
		}
		ands[i] = w.sql
		if err := mergeParams(params, w.params); err != nil {
			return "", nil, err
		}
	}
	sb.WriteString(" WHERE ")
	sb.WriteString(strings.Join(ands, " AND "))

	if len(b.returnings) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(b.returnings, ", "))
	}

	return sb.String(), params, nil
}
