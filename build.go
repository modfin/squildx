package squildx

import (
	"fmt"
	"strconv"
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

	if len(b.wheres) > 0 {
		ands := make([]string, len(b.wheres))
		for i, w := range b.wheres {
			if w.subQuery != nil {
				subSQL, subParams, err := w.subQuery.Build()
				if err != nil {
					return "", nil, err
				}
				ands[i] = fmt.Sprintf("%s (%s)", w.subPrefix, subSQL)
				if err := mergeParams(params, subParams); err != nil {
					return "", nil, err
				}
			} else {
				ands[i] = w.sql
				if err := mergeParams(params, w.params); err != nil {
					return "", nil, err
				}
			}
		}
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(ands, " AND "))
	}

	if len(b.groupBys) > 0 {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(strings.Join(b.groupBys, ", "))
	}

	if len(b.havings) > 0 {
		if len(b.groupBys) == 0 {
			return "", nil, ErrHavingWithoutGroupBy
		}
		ands := make([]string, len(b.havings))
		for i, h := range b.havings {
			ands[i] = h.sql
			if err := mergeParams(params, h.params); err != nil {
				return "", nil, err
			}
		}
		sb.WriteString(" HAVING ")
		sb.WriteString(strings.Join(ands, " AND "))
	}

	if len(b.orderBys) > 0 {
		exprs := make([]string, len(b.orderBys))
		for i, o := range b.orderBys {
			exprs[i] = o.sql
			if err := mergeParams(params, o.params); err != nil {
				return "", nil, err
			}
		}
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(exprs, ", "))
	}

	if b.limit != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.FormatUint(*b.limit, 10))
	}

	if b.offset != nil {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.FormatUint(*b.offset, 10))
	}

	return sb.String(), params, nil
}
