package squildx

import (
	"errors"
	"testing"
)

func TestUpdateBuild_NoTable(t *testing.T) {
	_, _, err := NewUpdate().
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if !errors.Is(err, ErrUpdateNoTable) {
		t.Errorf("expected ErrUpdateNoTable, got: %v", err)
	}
}

func TestUpdateBuild_NoSet(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		Where("id = :id", Params{"id": 1}).
		Build()

	if !errors.Is(err, ErrUpdateNoSet) {
		t.Errorf("expected ErrUpdateNoSet, got: %v", err)
	}
}

func TestUpdateBuild_NoWhere(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Build()

	if !errors.Is(err, ErrUpdateNoWhere) {
		t.Errorf("expected ErrUpdateNoWhere, got: %v", err)
	}
}

func TestUpdateBuild_StoredError(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice", "extra": 2}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if !errors.Is(err, ErrExtraParam) {
		t.Errorf("expected ErrExtraParam, got: %v", err)
	}
}

func TestUpdateBuild_Full(t *testing.T) {
	q, params, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Set("age = :age", Params{"age": 30}).
		Where("id = :id", Params{"id": 1}).
		Where("active = :active", Params{"active": true}).
		Returning("id", "name", "age").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name, age = :age WHERE id = :id AND active = :active RETURNING id, name, age"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "age", 30)
	assertParam(t, params, "id", 1)
	assertParam(t, params, "active", true)
}

func TestUpdateBuild_SubqueryBuildFailure(t *testing.T) {
	sub := New().Select("1") // missing FROM

	_, _, err := NewUpdate().
		Table("users").
		Set("active = false").
		WhereExists(sub).
		Build()

	if err == nil {
		t.Fatal("expected error from subquery build failure, got nil")
	}
	if !errors.Is(err, ErrNoFrom) {
		t.Errorf("expected ErrNoFrom, got: %v", err)
	}
}

func TestUpdateBuild_SetWhereParamCollision(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		Set("name = :val", Params{"val": "Alice"}).
		Where("id = :val", Params{"val": 1}).
		Build()

	if err == nil {
		t.Fatal("expected ErrDuplicateParam, got nil")
	}
	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}
