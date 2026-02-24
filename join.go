package squildx

import (
	"fmt"
	"reflect"
)

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

func (b *builder) addJoin(jt joinType, sql string, values []any) *builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	for _, j := range cp.joins {
		if j.joinType != jt || j.clause.sql != sql {
			continue
		}
		if reflect.DeepEqual(j.clause.params, params) {
			return cp
		}
		cp.err = fmt.Errorf("%w: join %s %s", ErrDuplicateParam, jt, sql)
		return cp
	}
	cp.joins = append(cp.joins, joinClause{
		joinType: jt,
		clause:   paramClause{sql: sql, params: params},
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

func (b *builder) addJoinLateral(jt joinType, sub Builder, alias string, on string, values []any) *builder {
	cp := b.clone()
	params, err := parseParams(on, values)
	if err != nil {
		cp.err = err
		return cp
	}
	for _, j := range cp.joins {
		if j.joinType != jt || j.alias != alias {
			continue
		}
		if j.subQuery == sub && j.clause.sql == on && reflect.DeepEqual(j.clause.params, params) {
			return cp
		}
		cp.err = fmt.Errorf("%w: lateral join %s %s", ErrDuplicateParam, jt, alias)
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

func (b *builder) InnerJoinLateral(sub Builder, alias string, on string, values ...any) Builder {
	return b.addJoinLateral(innerJoinLateral, sub, alias, on, values)
}

func (b *builder) LeftJoinLateral(sub Builder, alias string, on string, values ...any) Builder {
	return b.addJoinLateral(leftJoinLateral, sub, alias, on, values)
}

// CrossJoinLateral has no ON clause, so empty sql and nil values are passed through to
// addJoinLateral where parseParams handles them as a no-op.
func (b *builder) CrossJoinLateral(sub Builder, alias string) Builder {
	return b.addJoinLateral(crossJoinLateral, sub, alias, "", nil)
}
