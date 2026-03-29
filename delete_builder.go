package squildx

// DeleteBuilder provides a fluent, immutable API for constructing DELETE queries.
type DeleteBuilder interface {
	From(table string) DeleteBuilder
	Where(sql string, params ...Params) DeleteBuilder
	WhereExists(sub Builder) DeleteBuilder
	WhereNotExists(sub Builder) DeleteBuilder
	WhereIn(column string, sub Builder) DeleteBuilder
	WhereNotIn(column string, sub Builder) DeleteBuilder
	Returning(columns ...string) DeleteBuilder
	ReturningObject(obj any) DeleteBuilder
	Build() (string, Params, error)
}

type deleteBuilder struct {
	table       string
	wheres      []paramClause
	returnings  []string
	paramPrefix byte
	err         error
}

func NewDelete() DeleteBuilder {
	return &deleteBuilder{}
}

func (b *deleteBuilder) clone() *deleteBuilder {
	cp := *b
	cp.wheres = copySlice(b.wheres)
	cp.returnings = copySlice(b.returnings)
	return &cp
}

func (b *deleteBuilder) setPrefix(prefix byte) error {
	return checkSetPrefix(&b.paramPrefix, prefix)
}
