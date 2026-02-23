# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

squildx is a Go SQL query builder library for use with `sqlx` using named parameters (`:param_name` style). It provides a fluent, chainable API for constructing SELECT queries.

## Commands

```bash
go test ./...          # Run all tests
go test -v ./...       # Run all tests with verbose output
go test -run TestName  # Run a single test by name
go build ./...         # Build/check compilation
```

## Architecture

Single-package library using the **immutable Builder pattern** — every method clones the builder before modifying, so partial queries can be safely reused.

**Key types:**
- `Builder` interface (`builder.go`) — public API with `Select`, `From`, `Where`, `InnerJoin`/`LeftJoin`/`RightJoin`/`FullJoin`, `OrderBy`, `Limit`, `Offset`, `Build`
- `builder` struct (`builder.go`) — internal state; created via `New()`
- `joinClause` / `whereClause` — internal clause representations
- `Build()` (`build.go`) — assembles final SQL string and merged `map[string]any` params

**Parameter system** (`params.go`): Named placeholders are extracted via regex, matched positionally against variadic `values ...any` args, and merged across all clauses at build time. Duplicate param names with different values produce `ErrDuplicateParam`.

**Error handling**: Errors from parameter parsing are stored in the builder and only surfaced when `Build()` is called, allowing uninterrupted method chaining.

**Code layout**: Each SQL clause (select, from, where, join, orderby, limit, offset) has its own file and corresponding `_test.go` file.
