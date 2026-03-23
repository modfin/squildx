package squildx

import (
	"fmt"
	"strings"
)

func (b *deleteBuilder) Build() (string, Params, error) {
	if b.err != nil {
		return "", nil, b.err
	}

	if b.table == "" {
		return "", nil, ErrDeleteNoTable
	}
	if len(b.wheres) == 0 {
		return "", nil, ErrDeleteNoWhere
	}

	params := make(Params)
	prefix := b.paramPrefix

	var sb strings.Builder

	sb.WriteString("DELETE FROM ")
	sb.WriteString(b.table)

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
