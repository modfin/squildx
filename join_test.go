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
