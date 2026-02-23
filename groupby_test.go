package squildx

import (
	"testing"
)

func TestGroupByBasic(t *testing.T) {
	q, _, err := New().
		Select("department", "COUNT(*) AS cnt").
		From("employees").
		GroupBy("department").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, COUNT(*) AS cnt FROM employees GROUP BY department"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestGroupByMultipleCalls(t *testing.T) {
	q, _, err := New().
		Select("department", "role", "COUNT(*)").
		From("employees").
		GroupBy("department").
		GroupBy("role").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, role, COUNT(*) FROM employees GROUP BY department, role"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestGroupByMultipleExprs(t *testing.T) {
	q, _, err := New().
		Select("department", "role", "COUNT(*)").
		From("employees").
		GroupBy("department", "role").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, role, COUNT(*) FROM employees GROUP BY department, role"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestGroupByImmutability(t *testing.T) {
	base := New().Select("department", "COUNT(*)").From("employees")

	withGroup := base.GroupBy("department")

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := withGroup.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 == q2 {
		t.Error("expected different SQL for base and grouped builder")
	}

	expected := "SELECT department, COUNT(*) FROM employees"
	if q1 != expected {
		t.Errorf("base builder was mutated\n got: %s\nwant: %s", q1, expected)
	}
}
