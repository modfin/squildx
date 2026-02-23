package squildx

import (
	"errors"
	"testing"
)

func TestErrNoFrom(t *testing.T) {
	_, _, err := New().Select("id").Build()
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom, got: %v", err)
	}
}

func TestNamedTableFrom(t *testing.T) {
	q, _, err := New().Select("u.id").From("users u").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.id FROM users u"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
