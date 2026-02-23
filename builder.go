package squildx

type Builder interface {
	Select(columns ...string) Builder
	From(table string) Builder
	InnerJoin(sql string, values ...any) Builder
	LeftJoin(sql string, values ...any) Builder
	RightJoin(sql string, values ...any) Builder
	FullJoin(sql string, values ...any) Builder
	Where(sql string, values ...any) Builder
	OrderBy(exprs ...string) Builder
	Limit(n uint64) Builder
	Offset(n uint64) Builder
	Build() (string, map[string]any, error)
}

type builder struct {
	columns  []string
	from     string
	joins    []joinClause
	wheres   []whereClause
	ors      []whereClause
	orderBys []string
	limit    *uint64
	offset   *uint64
	err      error
}

func New() Builder {
	return &builder{}
}

func (b *builder) clone() *builder {
	cp := *b
	cp.columns = copySlice(b.columns)
	cp.joins = copySlice(b.joins)
	cp.wheres = copySlice(b.wheres)
	cp.ors = copySlice(b.ors)
	cp.orderBys = copySlice(b.orderBys)
	return &cp
}
