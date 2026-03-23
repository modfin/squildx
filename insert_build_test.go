package squildx

import (
	"errors"
	"testing"
)

func TestInsertBuild_NoTable(t *testing.T) {
	q := NewInsert().Columns("name").Values(":name", Params{"name": "Alice"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrNoTable) {
		t.Errorf("expected ErrNoTable, got: %v", err)
	}
}

func TestInsertBuild_NoColumns(t *testing.T) {
	q := NewInsert().Into("users").Values(":name", Params{"name": "Alice"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrNoInsertColumns) {
		t.Errorf("expected ErrNoInsertColumns, got: %v", err)
	}
}

func TestInsertBuild_NoValues(t *testing.T) {
	q := NewInsert().Into("users").Columns("name")
	_, _, err := q.Build()
	if !errors.Is(err, ErrNoInsertValues) {
		t.Errorf("expected ErrNoInsertValues, got: %v", err)
	}
}

func TestInsertBuild_ValuesAndSelect(t *testing.T) {
	sub := New().Select("name").From("temp")
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		Select(sub)
	_, _, err := q.Build()
	if !errors.Is(err, ErrValuesAndSelect) {
		t.Errorf("expected ErrValuesAndSelect, got: %v", err)
	}
}

func TestInsertBuild_DeferredError(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name") // missing param
	_, _, err := q.Build()
	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestInsertBuild_DuplicateParamConflict(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		OnConflictDoUpdate([]string{"id"}, "name = :name", Params{"name": "Bob"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestInsertBuild_DuplicateParamSameValue(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		OnConflictDoUpdate([]string{"id"}, "name = :name", Params{"name": "Alice"})
	_, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInsertBuild_SelectSubqueryError(t *testing.T) {
	sub := New().Select("name") // missing FROM
	q := NewInsert().Into("users").Columns("name").Select(sub)
	_, _, err := q.Build()
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom from subquery, got: %v", err)
	}
}

func TestInsertBuild_MixedPrefix(t *testing.T) {
	q := NewInsert().Into("users").Columns("name", "email").
		Values(":name, @email", Params{"name": "Alice", "email": "a@b.com"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}
