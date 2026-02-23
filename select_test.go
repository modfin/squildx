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
