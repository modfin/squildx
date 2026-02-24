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
- `Builder` interface (`builder.go`) — public API with `Select`, `Distinct`, `From`, `Where`, `WhereExists`/`WhereNotExists`/`WhereIn`/`WhereNotIn`, `InnerJoin`/`LeftJoin`/`RightJoin`/`FullJoin`/`CrossJoin`, `InnerJoinLateral`/`LeftJoinLateral`/`CrossJoinLateral`, `GroupBy`, `Having`, `OrderBy`, `Limit`, `Offset`, `Build`
- `builder` struct (`builder.go`) — internal state; created via `New()`
- `joinClause` / `paramClause` — internal clause representations. Where subquery methods (`WhereExists`, `WhereIn`, etc.) use `paramClause.subQuery` to embed a nested `Builder`.
- `Build()` (`build.go`) — assembles final SQL string and merged `map[string]any` params

**Parameter system** (`params.go`): Named placeholders are extracted via regex, matched positionally against variadic `values ...any` args, and merged across all clauses at build time. Duplicate param names with different values produce `ErrDuplicateParam`.

**Error handling**: Errors from parameter parsing are stored in the builder and only surfaced when `Build()` is called, allowing uninterrupted method chaining.

**OrderBy with parameters**: `OrderBy` accepts named parameters (e.g., `OrderBy("similarity(embedding, :vec) DESC", vec)`), following the same parameter system as `Where`.

**Code layout**: Each SQL clause (select, from, where, join, orderby, groupby, having, limit, offset) has its own file and corresponding `_test.go` file.

## Style Rules

- **Never use `if/else` or `else if`**. Use `switch` statements instead for multi-branch logic, but only when absolutely necessary — prefer early returns or guard clauses to avoid branching altogether.
