package squildx

import "fmt"

type joinType string

const (
	innerJoin        joinType = "INNER JOIN"
	leftJoin         joinType = "LEFT JOIN"
	rightJoin        joinType = "RIGHT JOIN"
	fullJoin         joinType = "FULL JOIN"
	crossJoin        joinType = "CROSS JOIN"
	innerJoinLateral joinType = "INNER JOIN LATERAL"
	leftJoinLateral  joinType = "LEFT JOIN LATERAL"
	crossJoinLateral joinType = "CROSS JOIN LATERAL"
)

type joinClause struct {
	joinType joinType
	clause   paramClause
	subQuery Builder
	alias    string
}

func (b *builder) addJoin(jt joinType, sql string, maps []map[string]any) *builder {
	cp := b.clone()
	p, err := extractParams(maps)
	if err != nil {
		cp.err = err
		return cp
	}
	params, prefix, err := parseParams(sql, p)
	if err != nil {
		cp.err = err
		return cp
	}
	if prefix != 0 {
		if cp.paramPrefix != 0 && cp.paramPrefix != prefix {
			cp.err = ErrMixedPrefix
			return cp
		}
		cp.paramPrefix = prefix
	}
	for _, j := range cp.joins {
		if j.joinType != jt || j.clause.sql != sql {
			continue
		}
		if paramsEqual(j.clause.params, params) {
			return cp
		}
		cp.err = fmt.Errorf("%w: %s %s", ErrDuplicateJoin, jt, sql)
		return cp
	}
	cp.joins = append(cp.joins, joinClause{
		joinType: jt,
		clause:   paramClause{sql: sql, params: params},
	})
	return cp
}

func (b *builder) InnerJoin(sql string, params ...map[string]any) Builder {
	return b.addJoin(innerJoin, sql, params)
}

func (b *builder) LeftJoin(sql string, params ...map[string]any) Builder {
	return b.addJoin(leftJoin, sql, params)
}

func (b *builder) RightJoin(sql string, params ...map[string]any) Builder {
	return b.addJoin(rightJoin, sql, params)
}

func (b *builder) FullJoin(sql string, params ...map[string]any) Builder {
	return b.addJoin(fullJoin, sql, params)
}

func (b *builder) CrossJoin(sql string, params ...map[string]any) Builder {
	return b.addJoin(crossJoin, sql, params)
}

func (b *builder) addJoinLateral(jt joinType, sub Builder, alias string, on string, maps []map[string]any) *builder {
	cp := b.clone()
	p, err := extractParams(maps)
	if err != nil {
		cp.err = err
		return cp
	}
	params, prefix, err := parseParams(on, p)
	if err != nil {
		cp.err = err
		return cp
	}
	if prefix != 0 {
		if cp.paramPrefix != 0 && cp.paramPrefix != prefix {
			cp.err = ErrMixedPrefix
			return cp
		}
		cp.paramPrefix = prefix
	}
	for _, j := range cp.joins {
		if j.joinType != jt || j.alias != alias {
			continue
		}
		if j.clause.sql == on && paramsEqual(j.clause.params, params) && buildersEqual(j.subQuery, sub) {
			return cp
		}
		cp.err = fmt.Errorf("%w: %s LATERAL %s", ErrDuplicateJoin, jt, alias)
		return cp
	}
	cp.joins = append(cp.joins, joinClause{
		joinType: jt,
		clause:   paramClause{sql: on, params: params},
		subQuery: sub,
		alias:    alias,
	})
	return cp
}

func (b *builder) InnerJoinLateral(sub Builder, alias string, on string, params ...map[string]any) Builder {
	return b.addJoinLateral(innerJoinLateral, sub, alias, on, params)
}

func (b *builder) LeftJoinLateral(sub Builder, alias string, on string, params ...map[string]any) Builder {
	return b.addJoinLateral(leftJoinLateral, sub, alias, on, params)
}

// CrossJoinLateral has no ON clause, so empty sql and nil maps are passed through to
// addJoinLateral where parseParams handles them as a no-op.
func (b *builder) CrossJoinLateral(sub Builder, alias string) Builder {
	return b.addJoinLateral(crossJoinLateral, sub, alias, "", nil)
}
