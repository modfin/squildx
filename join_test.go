package squildx

import (
	"errors"
	"testing"
)

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

func TestDoubleJoin(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		InnerJoin("orders o ON o.user_id = u.id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDoubleJoinWithMatchingParams(t *testing.T) {
	q, params, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id AND o.status = :status", "active").
		InnerJoin("orders o ON o.user_id = u.id AND o.status = :status", "active").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id AND o.status = :status"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "status", "active")
}

func TestDoubleJoinWithConflictingParams(t *testing.T) {
	_, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id AND o.status = :status", "active").
		InnerJoin("orders o ON o.user_id = u.id AND o.status = :status", "inactive").
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestDifferentJoinTypesSameSql(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		LeftJoin("orders o ON o.user_id = u.id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id LEFT JOIN orders o ON o.user_id = u.id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDoubleLeftJoin(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		LeftJoin("profiles p ON p.user_id = u.id").
		LeftJoin("profiles p ON p.user_id = u.id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u LEFT JOIN profiles p ON p.user_id = u.id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDoubleCrossJoin(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		CrossJoin("colors c").
		CrossJoin("colors c").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u CROSS JOIN colors c"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestMixedDuplicateAndUniqueJoins(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		InnerJoin("orders o ON o.user_id = u.id").
		LeftJoin("profiles p ON p.user_id = u.id").
		InnerJoin("orders o ON o.user_id = u.id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u INNER JOIN orders o ON o.user_id = u.id LEFT JOIN profiles p ON p.user_id = u.id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDoubleLeftJoinLateral(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(3)

	q, _, err := New().
		Select("u.name", "recent.*").
		From("users u").
		LeftJoinLateral(sub, "recent", "true").
		LeftJoinLateral(sub, "recent", "true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, recent.* FROM users u LEFT JOIN LATERAL (SELECT * FROM orders o WHERE o.user_id = u.id LIMIT 3) recent ON true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDoubleJoinLateralConflictingSubquery(t *testing.T) {
	sub1 := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(3)
	sub2 := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(5)

	_, _, err := New().
		Select("u.name", "recent.*").
		From("users u").
		LeftJoinLateral(sub1, "recent", "true").
		LeftJoinLateral(sub2, "recent", "true").
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
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

func TestCrossJoin(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users u").
		CrossJoin("colors c").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u CROSS JOIN colors c"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestCrossJoinWithParams(t *testing.T) {
	q, params, err := New().
		Select("*").
		From("users u").
		CrossJoin("(SELECT * FROM sizes WHERE active = :active) s", true).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users u CROSS JOIN (SELECT * FROM sizes WHERE active = :active) s"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "active", true)
}

func TestLeftJoinLateral(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.user_id = u.id").OrderBy("o.created_at DESC").Limit(3)

	q, _, err := New().
		Select("u.name", "recent.*").
		From("users u").
		LeftJoinLateral(sub, "recent", "true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, recent.* FROM users u LEFT JOIN LATERAL (SELECT * FROM orders o WHERE o.user_id = u.id ORDER BY o.created_at DESC LIMIT 3) recent ON true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestInnerJoinLateral(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(1)

	q, _, err := New().
		Select("u.name", "latest.*").
		From("users u").
		InnerJoinLateral(sub, "latest", "true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, latest.* FROM users u INNER JOIN LATERAL (SELECT * FROM orders o WHERE o.user_id = u.id LIMIT 1) latest ON true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestCrossJoinLateral(t *testing.T) {
	sub := New().Select("*").From("generate_series(1, 3)")

	q, _, err := New().
		Select("u.name", "s.*").
		From("users u").
		CrossJoinLateral(sub, "s").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, s.* FROM users u CROSS JOIN LATERAL (SELECT * FROM generate_series(1, 3)) s"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestLateralJoinWithOnParams(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(3)

	q, params, err := New().
		Select("u.name", "recent.*").
		From("users u").
		LeftJoinLateral(sub, "recent", "recent.amount > :min", 100).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, recent.* FROM users u LEFT JOIN LATERAL (SELECT * FROM orders o WHERE o.user_id = u.id LIMIT 3) recent ON recent.amount > :min"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min", 100)
}

func TestLateralJoinWithSubqueryParams(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.status = :status", "active").Limit(3)

	q, params, err := New().
		Select("u.name", "recent.*").
		From("users u").
		LeftJoinLateral(sub, "recent", "true").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, recent.* FROM users u LEFT JOIN LATERAL (SELECT * FROM orders o WHERE o.status = :status LIMIT 3) recent ON true"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "status", "active")
}

func TestLateralJoinParamCollision(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.status = :status", "active")

	_, _, err := New().
		Select("*").
		From("users u").
		Where("u.status = :status", "inactive").
		LeftJoinLateral(sub, "recent", "true").
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestLateralJoinCombinedWithRegularJoins(t *testing.T) {
	sub := New().Select("*").From("orders o").Where("o.user_id = u.id").Limit(3)

	q, params, err := New().
		Select("u.name", "p.bio", "recent.*").
		From("users u").
		LeftJoin("profiles p ON p.user_id = u.id").
		LeftJoinLateral(sub, "recent", "true").
		Where("u.active = :active", true).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT u.name, p.bio, recent.* FROM users u LEFT JOIN profiles p ON p.user_id = u.id LEFT JOIN LATERAL (SELECT * FROM orders o WHERE o.user_id = u.id LIMIT 3) recent ON true WHERE u.active = :active"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "active", true)
}

func TestLateralJoinSubqueryBuildError(t *testing.T) {
	sub := New().Select("*") // missing From — will error on Build

	_, _, err := New().
		Select("u.name").
		From("users u").
		LeftJoinLateral(sub, "recent", "true").
		Build()

	if err == nil {
		t.Fatal("expected ErrNoFrom, got nil")
	}
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom, got: %v", err)
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
