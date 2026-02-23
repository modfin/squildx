package squildx

type joinType string

const (
	innerJoin joinType = "INNER JOIN"
	leftJoin  joinType = "LEFT JOIN"
	rightJoin joinType = "RIGHT JOIN"
	fullJoin  joinType = "FULL JOIN"
	crossJoin joinType = "CROSS JOIN"
)

type joinClause struct {
	joinType joinType
	clause   whereClause
}

func (b *builder) addJoin(jt joinType, sql string, values []any) *builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.joins = append(cp.joins, joinClause{
		joinType: jt,
		clause:   whereClause{sql: sql, params: params},
	})
	return cp
}

func (b *builder) InnerJoin(sql string, values ...any) Builder {
	return b.addJoin(innerJoin, sql, values)
}

func (b *builder) LeftJoin(sql string, values ...any) Builder {
	return b.addJoin(leftJoin, sql, values)
}

func (b *builder) RightJoin(sql string, values ...any) Builder {
	return b.addJoin(rightJoin, sql, values)
}

func (b *builder) FullJoin(sql string, values ...any) Builder {
	return b.addJoin(fullJoin, sql, values)
}

func (b *builder) CrossJoin(sql string, values ...any) Builder {
	return b.addJoin(crossJoin, sql, values)
}
