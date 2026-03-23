package squildx

import (
	"errors"
	"testing"
)

func TestUpdateWhere(t *testing.T) {
	q, params, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "id", 1)
}

func TestUpdateMultipleWheres(t *testing.T) {
	q, params, err := NewUpdate().
		Table("users").
		Set("active = :active", Params{"active": false}).
		Where("age > :min_age", Params{"min_age": 25}).
		Where("role = :role", Params{"role": "admin"}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET active = :active WHERE age > :min_age AND role = :role"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 25)
	assertParam(t, params, "role", "admin")
}

func TestUpdateWhereExists(t *testing.T) {
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	q, _, err := NewUpdate().
		Table("users").
		Set("has_orders = true").
		WhereExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET has_orders = true WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateWhereNotExists(t *testing.T) {
	sub := New().Select("1").From("sessions").Where("sessions.user_id = users.id")

	q, _, err := NewUpdate().
		Table("users").
		Set("active = false").
		WhereNotExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET active = false WHERE NOT EXISTS (SELECT 1 FROM sessions WHERE sessions.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateWhereIn(t *testing.T) {
	sub := New().Select("user_id").From("vip_users")

	q, _, err := NewUpdate().
		Table("users").
		Set("tier = 'vip'").
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET tier = 'vip' WHERE id IN (SELECT user_id FROM vip_users)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateWhereNotIn(t *testing.T) {
	sub := New().Select("user_id").From("active_users")

	q, _, err := NewUpdate().
		Table("users").
		Set("active = false").
		WhereNotIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET active = false WHERE id NOT IN (SELECT user_id FROM active_users)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestUpdateWhereSubqueryParamCollision(t *testing.T) {
	sub := New().Select("id").From("orders").Where("name = :name", Params{"name": "order_name"})

	_, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "user_name"}).
		WhereIn("id", sub).
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestUpdateWhere_Immutability(t *testing.T) {
	base := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1})

	_ = base.Where("active = :active", Params{"active": true})

	q, params, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("original builder mutated\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}
