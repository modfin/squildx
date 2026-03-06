package squildx

import (
	"errors"
	"testing"
)

func TestDuplicateParamSameValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", Params{"val": 18}).
		Where("score > :val", Params{"val": 18}).
		Build()

	if err != nil {
		t.Errorf("expected no error for duplicate param with same value, got: %v", err)
	}
}

func TestDuplicateParamDifferentValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", Params{"val": 18}).
		Where("score > :val", Params{"val": 99}).
		Build()

	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestMissingParamValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND age < :max", Params{"min": 18}).
		Build()

	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestExtraParamKey(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min", Params{"min": 18, "extra": 65}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestAtPrefixParams(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > @min_age", Params{"min_age": 25}).
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
		Where("age > :min_age", Params{"min_age": 18}).
		Where("active = @active", Params{"active": true}).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestMixedPrefixInSameClause(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND name = @name", Params{"min": 18, "name": "test"}).
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
		Where("age::integer > :min_age", Params{"min_age": 18}).
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
		Where("age > :min_age", Params{"min_age": 18, "unused": "value"}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestMissingPlaceholderValueError(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age AND name = :name", Params{"min_age": 18}).
		Build()

	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestMultipleParamMapsMerged(t *testing.T) {
	q, params, err := New().Select("*").
		From("users").
		Where("age > :min_age AND name = :name", Params{"min_age": 18}, Params{"name": "alice"}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > :min_age AND name = :name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 18)
	assertParam(t, params, "name", "alice")
}

func TestMultipleParamMapsDuplicateKey(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", Params{"min_age": 18}, Params{"min_age": 99}).
		Build()

	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestSubqueryMixedPrefixError(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = @status", Params{"status": "active"})

	_, _, err := New().Select("*").
		From("users").
		Where("age > :min_age", Params{"min_age": 18}).
		WhereIn("id", sub).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestAtPrefixWithSubquery(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = @status", Params{"status": "active"})

	q, params, err := New().Select("*").
		From("users").
		Where("age > @min_age", Params{"min_age": 18}).
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
		Where("age > @min_age", Params{"min_age": 18}).
		Where("active = @active", Params{"active": true}).
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
		Where("@@session_var = true AND name = @name", Params{"name": "test"}).
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

	_ = base.Where("age > :min_age", Params{"min_age": 18})

	// base should not have its prefix set — @ prefix should still work on base
	q, params, err := base.Where("name = @name", Params{"name": "test"}).Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE name = @name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "test")
}

func TestNestedSubqueryMixedPrefix(t *testing.T) {
	inner := New().Select("1").From("c").Where("z = @z", Params{"z": 3})
	middle := New().Select("1").From("b").WhereExists(inner)

	_, _, err := New().Select("*").From("a").
		Where("y = :y", Params{"y": 2}).
		WhereExists(middle).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestSiblingSubqueriesMixedPrefix(t *testing.T) {
	sub1 := New().Select("1").From("a").Where("x = :x", Params{"x": 1})
	sub2 := New().Select("1").From("b").Where("y = @y", Params{"y": 2})

	_, _, err := New().Select("*").From("c").
		WhereExists(sub1).
		WhereExists(sub2).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestSiblingSubqueriesMixedPrefixNoParentParams(t *testing.T) {
	sub1 := New().Select("id").From("a").Where("x = :x", Params{"x": 1})
	sub2 := New().Select("id").From("b").Where("y = @y", Params{"y": 2})

	_, _, err := New().Select("*").From("c").
		WhereIn("id", sub1).
		WhereNotIn("id", sub2).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestDeeplyNestedSubqueryMixedPrefix(t *testing.T) {
	level3 := New().Select("1").From("d").Where("w = @w", Params{"w": 4})
	level2 := New().Select("1").From("c").WhereExists(level3)
	level1 := New().Select("id").From("b").WhereIn("id", level2)

	_, _, err := New().Select("*").From("a").
		Where("y = :y", Params{"y": 2}).
		WhereExists(level1).
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestLateralJoinNestedSubqueryMixedPrefix(t *testing.T) {
	inner := New().Select("1").From("c").Where("z = @z", Params{"z": 3})
	sub := New().Select("*").From("b").WhereExists(inner)

	_, _, err := New().Select("*").From("a").
		Where("y = :y", Params{"y": 2}).
		LeftJoinLateral(sub, "lat", "true").
		Build()

	if !errors.Is(err, ErrMixedPrefix) {
		t.Errorf("expected ErrMixedPrefix, got: %v", err)
	}
}

func TestSiblingSubqueriesSamePrefixOk(t *testing.T) {
	sub1 := New().Select("1").From("a").Where("x = :x", Params{"x": 1})
	sub2 := New().Select("1").From("b").Where("y = :y", Params{"y": 2})

	q, params, err := New().Select("*").From("c").
		WhereExists(sub1).
		WhereExists(sub2).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM c WHERE EXISTS (SELECT 1 FROM a WHERE x = :x) AND EXISTS (SELECT 1 FROM b WHERE y = :y)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "x", 1)
	assertParam(t, params, "y", 2)
}
