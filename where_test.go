package squildx

import (
	"errors"
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

func TestSingleParamMultipleTimes(t *testing.T) {
	q, params, err := New().Select("*").
		From("article").
		Where("(title ILIKE '%' || :search || '%' OR text ILIKE '%' || :search || '%')", "test").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM article WHERE (title ILIKE '%' || :search || '%' OR text ILIKE '%' || :search || '%')"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "search", "test")
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

func TestWhereExists(t *testing.T) {
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	q, params, err := New().Select("*").
		From("users").
		WhereExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestWhereNotExists(t *testing.T) {
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	q, params, err := New().Select("*").
		From("users").
		WhereNotExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestWhereIn(t *testing.T) {
	sub := New().Select("user_id").From("orders").Where("total > :min_total", 100)

	q, params, err := New().Select("*").
		From("users").
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > :min_total)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_total", 100)
}

func TestWhereNotIn(t *testing.T) {
	sub := New().Select("user_id").From("banned_users")

	q, params, err := New().Select("*").
		From("users").
		WhereNotIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestWhereSubqueryParamsMerged(t *testing.T) {
	sub := New().Select("id").From("orders").Where("status = :status", "active")

	q, params, err := New().Select("*").
		From("users").
		Where("age > :min_age", 18).
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE age > :min_age AND id IN (SELECT id FROM orders WHERE status = :status)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_age", 18)
	assertParam(t, params, "status", "active")
}

func TestWhereSubqueryParamCollision(t *testing.T) {
	sub := New().Select("id").From("orders").Where("name = :name", "order_name")

	_, _, err := New().Select("*").
		From("users").
		Where("name = :name", "user_name").
		WhereIn("id", sub).
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestWhereSubqueryParamCollisionOk(t *testing.T) {
	sub := New().Select("id").From("orders").Where("name = :name", "order_name")

	q, params, err := New().Select("*").
		From("users").
		Where("name = :name", "order_name").
		WhereIn("id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE name = :name AND id IN (SELECT id FROM orders WHERE name = :name)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "name", "order_name")
}

func TestWhereSubqueryBuildFailure(t *testing.T) {
	sub := New().Select("1") // missing FROM

	_, _, err := New().Select("*").
		From("users").
		WhereExists(sub).
		Build()

	if err == nil {
		t.Fatal("expected error from subquery build failure, got nil")
	}
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom, got: %v", err)
	}
}

func TestWhereExistsWithRegularWhere(t *testing.T) {
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	q, params, err := New().Select("*").
		From("users").
		Where("active = :active", true).
		WhereExists(sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE active = :active AND EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "active", true)
}

func TestWhereSubqueryImmutability(t *testing.T) {
	base := New().Select("*").From("users").Where("active = :active", true)
	sub := New().Select("1").From("orders").Where("orders.user_id = users.id")

	_ = base.WhereExists(sub)

	q, params, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE active = :active"
	if q != expected {
		t.Errorf("original builder mutated\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "active", true)
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

func TestReusedParamKeepsOrder(t *testing.T) {
	q, params, err := New().Select("*").
		From("article").
		Where("(title ILIKE '%' || :search_1 || '%' OR text ILIKE '%' || :search_2 || '%' OR body ILIKE '%' || :search_1 || '%')", "first", "second").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM article WHERE (title ILIKE '%' || :search_1 || '%' OR text ILIKE '%' || :search_2 || '%' OR body ILIKE '%' || :search_1 || '%')"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "search_1", "first")
	assertParam(t, params, "search_2", "second")
}

func TestReusedParamCrossClauseConflict(t *testing.T) {
	_, _, err := New().Select("*").
		From("article").
		Where("(title ILIKE '%' || :search || '%' OR text ILIKE '%' || :search || '%')", "first").
		Where("body = :search", "second").
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestReusedParamCrossClauseSameValue(t *testing.T) {
	q, params, err := New().Select("*").
		From("article").
		Where("(title ILIKE '%' || :search || '%' OR text ILIKE '%' || :search || '%')", "same").
		Where("body = :search", "same").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM article WHERE (title ILIKE '%' || :search || '%' OR text ILIKE '%' || :search || '%') AND body = :search"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "search", "same")
}

func TestWhereInQualifiedColumn(t *testing.T) {
	sub := New().Select("user_id").From("orders")

	q, _, err := New().Select("*").
		From("users").
		WhereIn("users.id", sub).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE users.id IN (SELECT user_id FROM orders)"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
