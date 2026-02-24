# squildx

A Go SQL query builder for use with [sqlx](https://github.com/jmoiron/sqlx) using named parameters (`:param` style). Immutable builder pattern — every method returns a new builder, so partial queries can be safely reused.

## Install

```bash
go get github.com/modfin/squildx
```

## Usage

Basic query:

```go
query, params, err := squildx.New().
    Select("id", "name", "email").
    From("users").
    Where("active = :active", true).
    Build()

// query:  SELECT id, name, email FROM users WHERE active = :active
// params: map[active:true]
```

Conditionally building a query:

```go
q := squildx.New().
    Select("id", "name", "email").
    From("users").
    Where("active = :active", true)

if nameFilter != "" {
    q = q.Where("name ILIKE :name", nameFilter)
}

if sortByName {
    q = q.OrderBy("name ASC")
}

query, params, err := q.Build()
```

OrderBy with parameters:

```go
query, params, err := squildx.New().
    Select("id", "title").
    From("documents").
    OrderBy("similarity(embedding, :query_vec) DESC", vec).
    Build()

// query:  SELECT id, title FROM documents ORDER BY similarity(embedding, :query_vec) DESC
// params: map[query_vec:<vec>]
```

Where subqueries:

```go
// WHERE EXISTS
sub := squildx.New().Select("1").From("orders").Where("orders.user_id = users.id")

query, params, err := squildx.New().
    Select("*").
    From("users").
    WhereExists(sub).
    Build()

// query: SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)
```

```go
// WHERE IN with parameters
sub := squildx.New().Select("user_id").From("orders").Where("total > :min_total", 100)

query, params, err := squildx.New().
    Select("*").
    From("users").
    WhereIn("id", sub).
    Build()

// query:  SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > :min_total)
// params: map[min_total:100]
```

Also available: `WhereNotExists` and `WhereNotIn`.

Other features: `Distinct()`, `InnerJoinLateral`/`LeftJoinLateral`/`CrossJoinLateral`.
