package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestInsertReturning(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		Returning("id", "created_at")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name) VALUES (:name) RETURNING id, created_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertReturning_Chained(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		Returning("id").Returning("created_at")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name) VALUES (:name) RETURNING id, created_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertReturning_Immutability(t *testing.T) {
	base := NewInsert().Returning("id")
	_ = base.Returning("created_at")
	ib := base.(*insertBuilder)
	want := []string{"id"}
	if !reflect.DeepEqual(ib.returnings, want) {
		t.Errorf("base returnings = %v, want %v", ib.returnings, want)
	}
}

func TestInsertReturningObject(t *testing.T) {
	type Result struct {
		ID        int    `db:"id"`
		CreatedAt string `db:"created_at"`
	}
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		ReturningObject(Result{})
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name) VALUES (:name) RETURNING id, created_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertReturningObject_NotAStruct(t *testing.T) {
	q := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"}).
		ReturningObject("not a struct")
	_, _, err := q.Build()
	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}
