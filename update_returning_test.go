package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestUpdateReturning(t *testing.T) {
	q, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		Returning("id", "name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateReturning_Chained(t *testing.T) {
	q, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		Returning("id").
		Returning("name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateReturning_Immutability(t *testing.T) {
	base := NewUpdate().Returning("id")
	_ = base.Returning("name")
	ub := base.(*updateBuilder)
	want := []string{"id"}
	if !reflect.DeepEqual(ub.returnings, want) {
		t.Errorf("base returnings = %v, want %v", ub.returnings, want)
	}
}

func TestUpdateReturningObject(t *testing.T) {
	type Result struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	q, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		ReturningObject(Result{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateReturningObject_NotAStruct(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		ReturningObject("not a struct").
		Build()

	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}
