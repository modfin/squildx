package squildx

import (
	"errors"
	"testing"
)

func TestDuplicateParamSameValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", map[string]any{"val": 18}).
		Where("score > :val", map[string]any{"val": 18}).
		Build()

	if err != nil {
		t.Errorf("expected no error for duplicate param with same value, got: %v", err)
	}
}

func TestDuplicateParamDifferentValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", map[string]any{"val": 18}).
		Where("score > :val", map[string]any{"val": 99}).
		Build()

	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestMissingParamValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND age < :max", map[string]any{"min": 18}).
		Build()

	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestExtraParamKey(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min", map[string]any{"min": 18, "extra": 65}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestAtPrefixParams(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > @min_age", map[string]any{"min_age": 25}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > @min_age"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 25)
}

func TestMixedPrefixError(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", map[string]any{"min_age": 18}).
		Where("active = @active", map[string]any{"active": true}).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestMixedPrefixInSameClause(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND name = @name", map[string]any{"min": 18, "name": "test"}).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestDoubleColonNotParam(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age::integer > 18").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age::integer > 18"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestDoubleAtNotParam(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("@@session_var = true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE @@session_var = true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestDoubleColonWithRealParam(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age::integer > :min_age", map[string]any{"min_age": 18}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age::integer > :min_age"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 18)
}

func TestExtraMapKeyError(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", map[string]any{"min_age": 18, "unused": "value"}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestMissingPlaceholderValueError(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age AND name = :name", map[string]any{"min_age": 18}).
		Build()

	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestMultipleParamMapsError(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", map[string]any{"min_age": 18}, map[string]any{"extra": 1}).
		Build()

	if !errors.Is(err, ErrMultipleParamMaps) {
		t.Errorf("expected ErrMultipleParamMaps, got: %v", err)
	}
}

func TestSubqueryMixedPrefixError(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = @status", map[string]any{"status": "active"})

	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", map[string]any{"min_age": 18}).
		WhereIn("id", sub).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestAtPrefixWithSubquery(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = @status", map[string]any{"status": "active"})

	q, params, err := New().Select("*").
		From("users").
		Where("age > @min_age", map[string]any{"min_age": 18}).
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > @min_age AND id IN (SELECT id FROM orders WHERE status = @status)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_age", 18)
	assertParam(t, params, "status", "active")
}

func TestAtPrefixMultipleClauses(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > @min_age", map[string]any{"min_age": 18}).
		Where("active = @active", map[string]any{"active": true}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > @min_age AND active = @active"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_age", 18)
	assertParam(t, params, "active", true)
}

func TestDoubleAtWithRealAtParam(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("@@session_var = true AND name = @name", map[string]any{"name": "test"}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE @@session_var = true AND name = @name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "test")
}

func TestPrefixImmutability(t *testing.T) {
	base := New().Select("*").From("users")

	_ = base.Where("age > :min_age", map[string]any{"min_age": 18})

	// base should not have its prefix set — @ prefix should still work on base
	q, params, err := base.Where("name = @name", map[string]any{"name": "test"}).Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE name = @name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "test")
}
