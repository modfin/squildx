package squildx

import (
	"errors"
	"testing"
)

func TestDeleteWhere(t *testing.T) {
	q, params, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "id", 1)
}

func TestDeleteMultipleWheres(t *testing.T) {
	q, params, err := NewDelete().
		From("users").
		Where("age > :min_age", Params{"min_age": 25}).
		Where("active = :active", Params{"active": false}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE age > :min_age AND active = :active"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 25)
	assertParam(t, params, "active", false)
}

func TestDeleteWhereExists(t *testing.T) {
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	q, params, err := NewDelete().
		From("users").
		WhereExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestDeleteWhereNotExists(t *testing.T) {
	sub := New().Select("1").From("sessions").Where("sessions.user_id = users.id")

	q, _, err := NewDelete().
		From("users").
		WhereNotExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM sessions WHERE sessions.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteWhereIn(t *testing.T) {
	sub := New().Select("user_id").From("banned_users")

	q, _, err := NewDelete().
		From("users").
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id IN (SELECT user_id FROM banned_users)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteWhereNotIn(t *testing.T) {
	sub := New().Select("user_id").From("active_users")

	q, _, err := NewDelete().
		From("users").
		WhereNotIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE id NOT IN (SELECT user_id FROM active_users)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDeleteWhereSubqueryParamsMerged(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = :status", Params{"status": "cancelled"})

	q, params, err := NewDelete().
		From("users").
		Where("role = :role", Params{"role": "guest"}).
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE role = :role AND id IN (SELECT id FROM orders WHERE status = :status)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "role", "guest")
	assertParam(t, params, "status", "cancelled")
}

func TestDeleteWhereSubqueryParamCollision(t *testing.T) {
	sub := New().Select("id").From("orders").Where("name = :name", Params{"name": "order_name"})

	_, _, err := NewDelete().
		From("users").
		Where("name = :name", Params{"name": "user_name"}).
		WhereIn("id", sub).
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestDeleteWhereImmutability(t *testing.T) {
	base := NewDelete().From("users").Where("active = :active", Params{"active": false})
	_ = base.Where("role = :role", Params{"role": "admin"})

	q, params, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE active = :active"
	if q != expected {
		t.Errorf("original builder mutated\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}
