package squildx

import "fmt"

func (b *builder) Where(sql string, values ...any) Builder {
	cp := b.clone()
	params, err := parseParams(sql, values)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.wheres = append(cp.wheres, paramClause{sql: sql, params: params})
	return cp
}

func (b *builder) WhereExists(sub Builder) Builder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{subQuery: sub, subPrefix: "EXISTS"})
	return cp
}

func (b *builder) WhereNotExists(sub Builder) Builder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{subQuery: sub, subPrefix: "NOT EXISTS"})
	return cp
}

func (b *builder) WhereIn(column string, sub Builder) Builder {
	return addWhereInSubquery(b, column, "IN", sub)
}

func (b *builder) WhereNotIn(column string, sub Builder) Builder {
	return addWhereInSubquery(b, column, "NOT IN", sub)
}

func addWhereInSubquery(b *builder, column, keyword string, sub Builder) *builder {
	cp := b.clone()
	cp.wheres = append(cp.wheres, paramClause{
		subQuery:  sub,
		subPrefix: fmt.Sprintf("%s %s", column, keyword),
	})
	return cp
}
