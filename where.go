package squildx

import "fmt"

func (b *builder) Where(sql string, params ...Params) Builder {
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
	if prefix != 0 {
		if cp.paramPrefix != 0 && cp.paramPrefix != prefix {
			cp.err = ErrMixedPrefix
			return cp
		}
		cp.paramPrefix = prefix
	}
	cp.wheres = append(cp.wheres, paramClause{sql: sql, params: parsed})
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
