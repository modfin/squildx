package squildx

// UpdateBuilder provides a fluent, immutable API for constructing UPDATE queries.
type UpdateBuilder interface {
	Table(table string) UpdateBuilder
	Set(sql string, params ...Params) UpdateBuilder
	SetObject(obj any) UpdateBuilder
	Where(sql string, params ...Params) UpdateBuilder
	WhereExists(sub Builder) UpdateBuilder
	WhereNotExists(sub Builder) UpdateBuilder
	WhereIn(column string, sub Builder) UpdateBuilder
	WhereNotIn(column string, sub Builder) UpdateBuilder
	Returning(columns ...string) UpdateBuilder
	ReturningObject(obj any) UpdateBuilder
	Build() (string, Params, error)
}

type updateBuilder struct {
	table       string
	sets        []paramClause
	wheres      []paramClause
	returnings  []string
	paramPrefix byte
	err         error
}

func NewUpdate() UpdateBuilder {
	return &updateBuilder{}
}

func (b *updateBuilder) clone() *updateBuilder {
	cp := *b
	cp.sets = copySlice(b.sets)
	cp.wheres = copySlice(b.wheres)
	cp.returnings = copySlice(b.returnings)
	return &cp
}

func (b *updateBuilder) setPrefix(prefix byte) error {
	return checkSetPrefix(&b.paramPrefix, prefix)
}
