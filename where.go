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
	return addWhereSubquery(b, "EXISTS", sub)
}

func (b *builder) WhereNotExists(sub Builder) Builder {
	return addWhereSubquery(b, "NOT EXISTS", sub)
}

func (b *builder) WhereIn(column string, sub Builder) Builder {
	return addWhereInSubquery(b, column, "IN", sub)
}

func (b *builder) WhereNotIn(column string, sub Builder) Builder {
	return addWhereInSubquery(b, column, "NOT IN", sub)
}

func addWhereSubquery(b *builder, keyword string, sub Builder) *builder {
	cp := b.clone()
	subSQL, subParams, err := sub.Build()
	if err != nil {
		cp.err = err
		return cp
	}
	clause := fmt.Sprintf("%s (%s)", keyword, subSQL)
	cp.wheres = append(cp.wheres, paramClause{sql: clause, params: subParams})
	return cp
}

func addWhereInSubquery(b *builder, column string, keyword string, sub Builder) *builder {
	cp := b.clone()
	subSQL, subParams, err := sub.Build()
	if err != nil {
		cp.err = err
		return cp
	}
	clause := fmt.Sprintf("%s %s (%s)", column, keyword, subSQL)
	cp.wheres = append(cp.wheres, paramClause{sql: clause, params: subParams})
	return cp
}
