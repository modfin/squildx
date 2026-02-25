package squildx

// Params is a named parameter map for SQL query building.
// It is interchangeable with map[string]any.
type Params map[string]any

type Builder interface {
	Select(columns ...string) Builder
	SelectObject(obj any, table ...string) Builder
	RemoveSelect(columns ...string) Builder
	Distinct() Builder

	From(table string) Builder

	InnerJoin(sql string, params ...Params) Builder
	LeftJoin(sql string, params ...Params) Builder
	RightJoin(sql string, params ...Params) Builder
	FullJoin(sql string, params ...Params) Builder
	CrossJoin(sql string, params ...Params) Builder

	InnerJoinLateral(sub Builder, alias string, on string, params ...Params) Builder
	LeftJoinLateral(sub Builder, alias string, on string, params ...Params) Builder
	CrossJoinLateral(sub Builder, alias string) Builder

	Where(sql string, params ...Params) Builder
	WhereExists(sub Builder) Builder
	WhereNotExists(sub Builder) Builder
	WhereIn(column string, sub Builder) Builder
	WhereNotIn(column string, sub Builder) Builder

	GroupBy(exprs ...string) Builder
	Having(sql string, params ...Params) Builder

	OrderBy(expr string, params ...Params) Builder

	Limit(n uint64) Builder
	Offset(n uint64) Builder

	Build() (string, Params, error)
}

type builder struct {
	columns     []string
	distinct    bool
	from        string
	joins       []joinClause
	wheres      []paramClause
	groupBys    []string
	havings     []paramClause
	orderBys    []paramClause
	limit       *uint64
	offset      *uint64
	paramPrefix byte // ':' or '@', 0 = not yet detected
	err         error
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
	params    Params
	subQuery  Builder
	subPrefix string
}
