package squildx

import "fmt"

func (b *deleteBuilder) Where(sql string, params ...Params) DeleteBuilder {
	cp := b.clone()
	p, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	parsed, prefix, err := parseParams(sql, p)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.wheres = append(cp.wheres, paramClause{sql: sql, params: parsed})
	return cp
}

func (b *deleteBuilder) WhereExists(sub Builder) DeleteBuilder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{subQuery: sub, subPrefix: "EXISTS"})
	return cp
}

func (b *deleteBuilder) WhereNotExists(sub Builder) DeleteBuilder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{subQuery: sub, subPrefix: "NOT EXISTS"})
	return cp
}

func (b *deleteBuilder) WhereIn(column string, sub Builder) DeleteBuilder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{
		subQuery:  sub,
		subPrefix: fmt.Sprintf("%s IN", column),
	})
	return cp
}

func (b *deleteBuilder) WhereNotIn(column string, sub Builder) DeleteBuilder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{
		subQuery:  sub,
		subPrefix: fmt.Sprintf("%s NOT IN", column),
	})
	return cp
}
