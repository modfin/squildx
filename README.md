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
