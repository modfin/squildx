package squildx

import (
	"testing"
)

func TestMultipleWheres(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > :min_age", 25).
		Where("active = :active", true).
		Where("role = :role", "admin").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > :min_age AND active = :active AND role = :role"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_age", 25)
	assertParam(t, params, "active", true)
	assertParam(t, params, "role", "admin")
}

func TestMultipleParamsInOneClause(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > :min AND age < :max", 18, 65).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > :min AND age < :max"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min", 18)
	assertParam(t, params, "max", 65)
}

func TestWhereNoParams(t *testing.T) {
	q, params, err := New().
		Select("*").
		From("users").
		Where("active = true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE active = true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestDBFunctionInWhere(t *testing.T) {
	q, params, err := New().
		Select("*").
		From("users").
		Where("created_at > func_get_date(:name)", "test").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE created_at > func_get_date(:name)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "name", "test")
}
