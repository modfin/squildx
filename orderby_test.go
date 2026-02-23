package squildx

import (
	"testing"
)

func TestOrderByMultiple(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users").
		OrderBy("name ASC", "age DESC").
		OrderBy("id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users ORDER BY name ASC, age DESC, id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
