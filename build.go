package squildx

import (
	"fmt"
	"strings"
)

func (b *builder) Build() (string, map[string]any, error) {
	if b.err != nil {
		return "", nil, b.err
	}

	if len(b.columns) == 0 {
		return "", nil, ErrNoColumns
	}
	if b.from == "" {
		return "", nil, ErrNoFrom
	}

	params := make(map[string]any)

	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(b.columns, ", "))

	sb.WriteString(" FROM ")
	sb.WriteString(b.from)

	for _, j := range b.joins {
		sb.WriteString(" ")
		sb.WriteString(string(j.joinType))
		sb.WriteString(" ")
		sb.WriteString(j.clause.sql)
		if err := mergeParams(params, j.clause.params); err != nil {
			return "", nil, err
		}
	}

	var whereParts []string

	if len(b.wheres) > 0 {
		ands := make([]string, len(b.wheres))
		for i, w := range b.wheres {
			ands[i] = w.sql
			if err := mergeParams(params, w.params); err != nil {
				return "", nil, err
			}
		}
		whereParts = append(whereParts, strings.Join(ands, " AND "))
	}

	if len(whereParts) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(whereParts, " AND "))
	}

	if len(b.orderBys) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orderBys, ", "))
	}

	if b.limit != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", *b.limit))
	}

	if b.offset != nil {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", *b.offset))
	}

	return sb.String(), params, nil
}
