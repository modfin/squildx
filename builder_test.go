package squildx

import (
	"errors"
	"testing"
)

func TestFullAPIExample(t *testing.T) {
	q, params, err := New().
		Select("u.name", "o.total").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		Where("age > :min_age", 18).
		Where("active = :active", true).
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

func TestErrNoColumns(t *testing.T) {
	_, _, err := New().Select().From("users").Build()
	if !errors.Is(err, ErrNoColumns) {
		t.Errorf("expected ErrNoColumns, got: %v", err)
	}
}

func TestErrNoFrom(t *testing.T) {
	_, _, err := New().Select("id").Build()
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom, got: %v", err)
	}
}

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

func TestJoinTypes(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		LeftJoin("profiles p ON p.user_id = u.id").
		RightJoin("accounts a ON a.user_id = u.id").
		FullJoin("logs l ON l.user_id = u.id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id LEFT JOIN profiles p ON p.user_id = u.id RIGHT JOIN accounts a ON a.user_id = u.id FULL JOIN logs l ON l.user_id = u.id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestJoinWithParams(t *testing.T) {
	q, params, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id AND o.status = :order_status", "complete").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id AND o.status = :order_status"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "order_status", "complete")
}

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

func TestLimitZero(t *testing.T) {
	q, _, err := New().Select("*").From("users").Limit(0).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users LIMIT 0"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestOffsetZero(t *testing.T) {
	q, _, err := New().Select("*").From("users").Offset(0).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users OFFSET 0"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDuplicateParamSameValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", 18).
		Where("score > :val", 18).
		Build()

	if err != nil {
		t.Errorf("expected no error for duplicate param with same value, got: %v", err)
	}
}

func TestDuplicateParamDifferentValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", 18).
		Where("score > :val", 99).
		Build()

	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestParamCountMismatch(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND age < :max", 18).
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Errorf("expected ErrParamMismatch, got: %v", err)
	}
}

func TestParamCountMismatchTooManyValues(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min", 18, 65).
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Errorf("expected ErrParamMismatch, got: %v", err)
	}
}

func TestImmutability(t *testing.T) {
	base := New().Select("*").From("users")

	q1, _, err := base.Where("active = :active", true).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := base.Where("role = :role", "admin").Build()
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
		q = q.Where("age = :age", f.Age)
	}
	if f.Name != "" {
		q = q.Where("name = :name", f.Name)
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

func TestParamMismatchInJoin(t *testing.T) {
	_, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.status = :s1 AND o.type = :s2", "active").
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Errorf("expected ErrParamMismatch, got: %v", err)
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
