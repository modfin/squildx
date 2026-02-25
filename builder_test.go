package squildx

import (
	"testing"
)

func TestFullAPIExample(t *testing.T) {
	q, params, err := New().
		Select("u.name", "o.total").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		Where("age > :min_age", map[string]any{"min_age": 18}).
		Where("active = :active", map[string]any{"active": true}).
		OrderBy("u.name ASC").
		Limit(10).
		Offset(20).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, o.total FROM users u INNER JOIN orders o ON o.user_id = u.id WHERE age > :min_age AND active = :active ORDER BY u.name ASC LIMIT 10 OFFSET 20"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_age", 18)
	assertParam(t, params, "active", true)

	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}

func TestImmutability(t *testing.T) {
	base := New().Select("*").From("users")

	q1, _, err := base.Where("active = :active", map[string]any{"active": true}).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := base.Where("role = :role", map[string]any{"role": "admin"}).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 != "SELECT * FROM users WHERE active = :active" {
		t.Errorf("q1 mismatch: %s", q1)
	}
	if q2 != "SELECT * FROM users WHERE role = :role" {
		t.Errorf("q2 mismatch: %s", q2)
	}
}

func TestConditionalFiltering(t *testing.T) {
	type filter struct {
		Age  int
		Name string
	}

	f := filter{Age: 25}

	q := New().Select("*").From("users")
	if f.Age != 0 {
		q = q.Where("age = :age", map[string]any{"age": f.Age})
	}
	if f.Name != "" {
		q = q.Where("name = :name", map[string]any{"name": f.Name})
	}

	sql, params, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age = :age"
	if sql != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, expected)
	}

	assertParam(t, params, "age", 25)
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

func assertParam(t *testing.T, params map[string]any, key string, expected any) {
	t.Helper()
	val, ok := params[key]
	if !ok {
		t.Errorf("missing param %q", key)
		return
	}
	if val != expected {
		t.Errorf("param %q = %v, want %v", key, val, expected)
	}
}
