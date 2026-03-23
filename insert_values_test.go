package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestInsertValues(t *testing.T) {
	q := NewInsert().Into("users").Columns("name", "email").
		Values(":name, :email", Params{"name": "Alice", "email": "a@b.com"})
	ib := q.(*insertBuilder)
	if len(ib.valueRows) != 1 {
		t.Fatalf("expected 1 value row, got %d", len(ib.valueRows))
	}
	if ib.valueRows[0].sql != ":name, :email" {
		t.Errorf("unexpected sql: %s", ib.valueRows[0].sql)
	}
}

func TestInsertValues_MultiRow(t *testing.T) {
	q := NewInsert().Into("users").Columns("name", "email").
		Values(":n1, :e1", Params{"n1": "Alice", "e1": "a@b.com"}).
		Values(":n2, :e2", Params{"n2": "Bob", "e2": "b@b.com"})
	ib := q.(*insertBuilder)
	if len(ib.valueRows) != 2 {
		t.Fatalf("expected 2 value rows, got %d", len(ib.valueRows))
	}
}

func TestInsertValues_Immutability(t *testing.T) {
	base := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"})
	_ = base.Values(":name2", Params{"name2": "Bob"})
	ib := base.(*insertBuilder)
	if len(ib.valueRows) != 1 {
		t.Errorf("base should have 1 value row, got %d", len(ib.valueRows))
	}
}

func TestInsertValues_MissingParam(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{})
	_, _, err := q.Build()
	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestInsertValues_ExtraParam(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice", "extra": "oops"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestInsertValuesObject(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	q := NewInsert().Into("users").
		ValuesObject(User{Name: "Alice", Email: "a@b.com"})
	ib := q.(*insertBuilder)

	wantCols := []string{"name", "email"}
	if !reflect.DeepEqual(ib.columns, wantCols) {
		t.Errorf("columns = %v, want %v", ib.columns, wantCols)
	}
	if len(ib.valueRows) != 1 {
		t.Fatalf("expected 1 value row, got %d", len(ib.valueRows))
	}
	if ib.valueRows[0].sql != ":name, :email" {
		t.Errorf("unexpected sql: %s", ib.valueRows[0].sql)
	}
	assertParam(t, ib.valueRows[0].params, "name", "Alice")
	assertParam(t, ib.valueRows[0].params, "email", "a@b.com")
}

func TestInsertValuesObject_NotAStruct(t *testing.T) {
	q := NewInsert().Into("users").ValuesObject(42)
	_, _, err := q.Build()
	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}

func TestInsertValuesObject_ColumnMismatch(t *testing.T) {
	type User struct {
		Name string `db:"name"`
	}
	q := NewInsert().Into("users").Columns("email").
		ValuesObject(User{Name: "Alice"})
	_, _, err := q.Build()
	if !errors.Is(err, ErrNoInsertColumns) {
		t.Errorf("expected ErrNoInsertColumns on column mismatch, got: %v", err)
	}
}

func TestInsertValuesObject_SetsColumnsOnce(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	q := NewInsert().Into("users").
		ValuesObject(User{Name: "Alice", Email: "a@b.com"}).
		ValuesObject(User{Name: "Bob", Email: "b@b.com"})
	ib := q.(*insertBuilder)
	if len(ib.valueRows) != 2 {
		t.Fatalf("expected 2 value rows, got %d", len(ib.valueRows))
	}
	wantCols := []string{"name", "email"}
	if !reflect.DeepEqual(ib.columns, wantCols) {
		t.Errorf("columns = %v, want %v", ib.columns, wantCols)
	}
}

func TestInsertValuesObject_EmbeddedStruct(t *testing.T) {
	type Base struct {
		ID int `db:"id"`
	}
	type User struct {
		Base
		Name string `db:"name"`
	}
	q := NewInsert().Into("users").ValuesObject(User{Base: Base{ID: 1}, Name: "Alice"})
	ib := q.(*insertBuilder)
	wantCols := []string{"id", "name"}
	if !reflect.DeepEqual(ib.columns, wantCols) {
		t.Errorf("columns = %v, want %v", ib.columns, wantCols)
	}
	assertParam(t, ib.valueRows[0].params, "id", 1)
	assertParam(t, ib.valueRows[0].params, "name", "Alice")
}
