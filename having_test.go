package squildx

import (
	"errors"
	"testing"
)

func TestHavingBasic(t *testing.T) {
	q, params, err := New().
		Select("department", "COUNT(*) AS cnt").
		From("employees").
		GroupBy("department").
		Having("COUNT(*) > :min_count", 5).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, COUNT(*) AS cnt FROM employees GROUP BY department HAVING COUNT(*) > :min_count"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_count", 5)
}

func TestHavingMultiple(t *testing.T) {
	q, params, err := New().
		Select("department", "COUNT(*) AS cnt", "AVG(salary) AS avg_sal").
		From("employees").
		GroupBy("department").
		Having("COUNT(*) > :min_count", 5).
		Having("AVG(salary) > :min_salary", 50000).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, COUNT(*) AS cnt, AVG(salary) AS avg_sal FROM employees GROUP BY department HAVING COUNT(*) > :min_count AND AVG(salary) > :min_salary"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	assertParam(t, params, "min_count", 5)
	assertParam(t, params, "min_salary", 50000)
}

func TestHavingWithoutGroupByError(t *testing.T) {
	_, _, err := New().
		Select("department", "COUNT(*)").
		From("employees").
		Having("COUNT(*) > :min_count", 5).
		Build()

	if !errors.Is(err, ErrHavingWithoutGroupBy) {
		t.Errorf("expected ErrHavingWithoutGroupBy, got: %v", err)
	}
}

func TestHavingNoParams(t *testing.T) {
	q, params, err := New().
		Select("department", "COUNT(*) AS cnt").
		From("employees").
		GroupBy("department").
		Having("COUNT(*) > 5").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT department, COUNT(*) AS cnt FROM employees GROUP BY department HAVING COUNT(*) > 5"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestHavingImmutability(t *testing.T) {
	base := New().Select("department", "COUNT(*)").From("employees").GroupBy("department")

	withHaving := base.Having("COUNT(*) > :min_count", 5)

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := withHaving.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 == q2 {
		t.Error("expected different SQL for base and having builder")
	}

	expected := "SELECT department, COUNT(*) FROM employees GROUP BY department"
	if q1 != expected {
		t.Errorf("base builder was mutated\n got: %s\nwant: %s", q1, expected)
	}
}
