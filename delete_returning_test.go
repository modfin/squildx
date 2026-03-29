package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestDeleteReturning(t *testing.T) {
	q, _, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1}).
		Returning("id", "name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteReturning_Chained(t *testing.T) {
	q, _, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1}).
		Returning("id").
		Returning("name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteReturning_Immutability(t *testing.T) {
	base := NewDelete().Returning("id")
	_ = base.Returning("name")
	db := base.(*deleteBuilder)
	want := []string{"id"}
	if !reflect.DeepEqual(db.returnings, want) {
		t.Errorf("base returnings = %v, want %v", db.returnings, want)
	}
}

func TestDeleteReturningObject(t *testing.T) {
	type Result struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	q, _, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1}).
		ReturningObject(Result{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id = :id RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteReturningObject_NotAStruct(t *testing.T) {
	_, _, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1}).
		ReturningObject("not a struct").
		Build()

	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}
