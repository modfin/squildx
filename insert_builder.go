package squildx

// InsertBuilder provides a fluent, immutable API for constructing INSERT queries.
type InsertBuilder interface {
	Into(table string) InsertBuilder
	Columns(columns ...string) InsertBuilder
	ColumnsObject(obj any) InsertBuilder
	Values(sql string, params ...Params) InsertBuilder
	ValuesObject(obj any) InsertBuilder
	Select(sub Builder) InsertBuilder
	OnConflictDoNothing(columns ...string) InsertBuilder
	OnConflictDoUpdate(columns []string, set string, params ...Params) InsertBuilder
	Returning(columns ...string) InsertBuilder
	ReturningObject(obj any) InsertBuilder
	Build() (string, Params, error)
}

type insertBuilder struct {
	table       string
	columns     []string
	valueRows   []paramClause
	selectQuery Builder
	conflict    *conflictClause
	returnings  []string
	paramPrefix byte
	err         error
}

type conflictClause struct {
	columns  []string
	doUpdate bool
	set      string
	params   Params
}

func NewInsert() InsertBuilder {
	return &insertBuilder{}
}

func (b *insertBuilder) clone() *insertBuilder {
	cp := *b
	cp.columns = copySlice(b.columns)
	cp.valueRows = copySlice(b.valueRows)
	cp.returnings = copySlice(b.returnings)
	if b.conflict != nil {
		cc := *b.conflict
		cc.columns = copySlice(b.conflict.columns)
		cp.conflict = &cc
	}
	return &cp
}

func (b *insertBuilder) setPrefix(prefix byte) error {
	if prefix == 0 {
		return nil
	}
	if b.paramPrefix != 0 && b.paramPrefix != prefix {
		return ErrMixedPrefix
	}
	b.paramPrefix = prefix
	return nil
}
