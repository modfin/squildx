package squildx

type Builder interface {
	Select(columns ...string) Builder
	SelectObject(obj any, table ...string) Builder
	RemoveSelect(columns ...string) Builder
	Distinct() Builder

	From(table string) Builder

	InnerJoin(sql string, values ...any) Builder
	LeftJoin(sql string, values ...any) Builder
	RightJoin(sql string, values ...any) Builder
	FullJoin(sql string, values ...any) Builder
	CrossJoin(sql string, values ...any) Builder

	InnerJoinLateral(sub Builder, alias string, on string, values ...any) Builder
	LeftJoinLateral(sub Builder, alias string, on string, values ...any) Builder
	CrossJoinLateral(sub Builder, alias string) Builder

	Where(sql string, values ...any) Builder
	WhereExists(sub Builder) Builder
	WhereNotExists(sub Builder) Builder
	WhereIn(column string, sub Builder) Builder
	WhereNotIn(column string, sub Builder) Builder

	GroupBy(exprs ...string) Builder
	Having(sql string, values ...any) Builder

	OrderBy(expr string, values ...any) Builder

	Limit(n uint64) Builder
	Offset(n uint64) Builder

	Build() (string, map[string]any, error)
}

type builder struct {
	columns  []string
	distinct bool
	from     string
	joins    []joinClause
	wheres   []paramClause
	groupBys []string
	havings  []paramClause
	orderBys []paramClause
	limit    *uint64
	offset   *uint64
	err      error
}

func New() Builder {
	return &builder{}
}

// clone performs a shallow copy of the builder with fresh slices.
// Fields containing Builder interfaces (e.g. subQuery in joinClause/paramClause)
// are shared, which is safe because the Builder is immutable — every method clones before mutating.
func (b *builder) clone() *builder {
	cp := *b
	cp.columns = copySlice(b.columns)
	cp.joins = copySlice(b.joins)
	cp.wheres = copySlice(b.wheres)
	cp.groupBys = copySlice(b.groupBys)
	cp.havings = copySlice(b.havings)
	cp.orderBys = copySlice(b.orderBys)
	return &cp
}

type paramClause struct {
	sql       string
	params    map[string]any
	subQuery  Builder
	subPrefix string
}
