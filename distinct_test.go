package squildx

import (
	"testing"
)

func TestDistinctBasic(t *testing.T) {
	q, _, err := New().
		Select("name", "email").
		From("users").
		Distinct().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctWithWhere(t *testing.T) {
	q, params, err := New().
		Select("name", "email").
		From("users").
		Distinct().
		Where("active = :active", true).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT name, email FROM users WHERE active = :active"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	if params["active"] != true {
		t.Errorf("expected active=true, got %v", params["active"])
	}
}

func TestDistinctIdempotent(t *testing.T) {
	q, _, err := New().
		Select("name", "email").
		From("users").
		Distinct().
		Distinct().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctImmutability(t *testing.T) {
	base := New().Select("name", "email").From("users")
	withDistinct := base.Distinct()

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := withDistinct.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 == q2 {
		t.Error("expected different SQL for base and distinct builder")
	}

	expected := "SELECT name, email FROM users"
	if q1 != expected {
		t.Errorf("base builder was mutated\n got: %s\nwant: %s", q1, expected)
	}
}

func TestDistinctOnBasic(t *testing.T) {
	q, _, err := New().
		Select("name", "email").
		From("users").
		DistinctOn("name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name) name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctOnMultipleColumns(t *testing.T) {
	q, _, err := New().
		Select("name", "email", "department").
		From("users").
		DistinctOn("name", "department").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name, department) name, email, department FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctOnMultipleCalls(t *testing.T) {
	q, _, err := New().
		Select("name", "email", "department").
		From("users").
		DistinctOn("name").
		DistinctOn("department").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name, department) name, email, department FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctOnImmutability(t *testing.T) {
	base := New().Select("name", "email").From("users")
	withDistinctOn := base.DistinctOn("name")

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q2, _, err := withDistinctOn.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 == q2 {
		t.Error("expected different SQL for base and distinctOn builder")
	}

	expected := "SELECT name, email FROM users"
	if q1 != expected {
		t.Errorf("base builder was mutated\n got: %s\nwant: %s", q1, expected)
	}
}

func TestDistinctOnOverridesDistinct(t *testing.T) {
	q, _, err := New().
		Select("name", "email").
		From("users").
		Distinct().
		DistinctOn("name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name) name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctOnOverridesDistinctReverseOrder(t *testing.T) {
	q, _, err := New().
		Select("name", "email").
		From("users").
		DistinctOn("name").
		Distinct().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name) name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctWithSelectObject(t *testing.T) {
	q, _, err := New().
		SelectObject(testUser{}).
		From("users").
		Distinct().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT id, first_name, email, created_at FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestDistinctOnWithOrderBy(t *testing.T) {
	q, _, err := New().
		Select("name", "email", "created_at").
		From("users").
		DistinctOn("name").
		OrderBy("name ASC").
		OrderBy("created_at DESC").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT DISTINCT ON (name) name, email, created_at FROM users ORDER BY name ASC, created_at DESC"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
