package squildx

import (
	"errors"
	"testing"
)

func TestDeleteBuild_NoTable(t *testing.T) {
	_, _, err := NewDelete().
		Where("id = :id", Params{"id": 1}).
		Build()

	if !errors.Is(err, ErrDeleteNoTable) {
		t.Errorf("expected ErrDeleteNoTable, got: %v", err)
	}
}

func TestDeleteBuild_NoWhere(t *testing.T) {
	_, _, err := NewDelete().
		From("users").
		Build()

	if !errors.Is(err, ErrDeleteNoWhere) {
		t.Errorf("expected ErrDeleteNoWhere, got: %v", err)
	}
}

func TestDeleteBuild_StoredError(t *testing.T) {
	_, _, err := NewDelete().
		From("users").
		Where("id = :id", Params{"id": 1, "extra": 2}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestDeleteBuild_Full(t *testing.T) {
	q, params, err := NewDelete().
		From("users").
		Where("age > :min_age", Params{"min_age": 25}).
		Where("active = :active", Params{"active": false}).
		Returning("id", "name").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "DELETE FROM users WHERE age > :min_age AND active = :active RETURNING id, name"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "min_age", 25)
	assertParam(t, params, "active", false)
}

func TestDeleteBuild_SubqueryBuildFailure(t *testing.T) {
	sub := New().Select("1") // missing FROM

	_, _, err := NewDelete().
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
