package squildx

import (
	"errors"
	"testing"
)

func TestSelectOnly(t *testing.T) {
	q, params, err := New().Select("*").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q != "SELECT * FROM users" {
		t.Errorf("got: %s", q)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestNoWhereClause(t *testing.T) {
	q, params, err := New().Select("id", "name").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q != "SELECT id, name FROM users" {
		t.Errorf("got: %s", q)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestErrNoColumns(t *testing.T) {
	_, _, err := New().Select().From("users").Build()
	if !errors.Is(err, ErrNoColumns) {
		t.Errorf("expected ErrNoColumns, got: %v", err)
	}
}

func TestSelectAppends(t *testing.T) {
	q, _, err := New().Select("id").Select("name", "email").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestRemoveSelect(t *testing.T) {
	q, _, err := New().Select("id", "name", "email").RemoveSelect("name").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestRemoveSelectImmutability(t *testing.T) {
	base := New().Select("id", "name", "email").From("users")
	reduced := base.RemoveSelect("name")

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q2, _, err := reduced.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 != "SELECT id, name, email FROM users" {
		t.Errorf("base was mutated: %s", q1)
	}
	if q2 != "SELECT id, email FROM users" {
		t.Errorf("reduced mismatch: %s", q2)
	}
}
