package squildx

import (
	"errors"
	"testing"
)

func TestInsertOnConflictDoNothing(t *testing.T) {
	q := NewInsert().Into("users").Columns("id", "name").
		Values(":id, :name", Params{"id": 1, "name": "Alice"}).
		OnConflictDoNothing("id")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name) VALUES (:id, :name) ON CONFLICT (id) DO NOTHING"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertOnConflictDoNothing_MultipleColumns(t *testing.T) {
	q := NewInsert().Into("users").Columns("id", "email", "name").
		Values(":id, :email, :name", Params{"id": 1, "email": "a@b.com", "name": "Alice"}).
		OnConflictDoNothing("id", "email")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, email, name) VALUES (:id, :email, :name) ON CONFLICT (id, email) DO NOTHING"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertOnConflictDoUpdate(t *testing.T) {
	q := NewInsert().Into("users").Columns("id", "name").
		Values(":id, :name", Params{"id": 1, "name": "Alice"}).
		OnConflictDoUpdate([]string{"id"}, "name = EXCLUDED.name")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name) VALUES (:id, :name) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertOnConflictDoUpdate_WithParams(t *testing.T) {
	q := NewInsert().Into("users").Columns("id", "name").
		Values(":id, :name", Params{"id": 1, "name": "Alice"}).
		OnConflictDoUpdate([]string{"id"}, "name = :conflict_name", Params{"conflict_name": "Bob"})
	sql, params, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name) VALUES (:id, :name) ON CONFLICT (id) DO UPDATE SET name = :conflict_name"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "conflict_name", "Bob")
}

func TestInsertOnConflictDoUpdate_MissingParam(t *testing.T) {
	q := NewInsert().Into("users").Columns("id").
		Values(":id", Params{"id": 1}).
		OnConflictDoUpdate([]string{"id"}, "name = :conflict_name")
	_, _, err := q.Build()
	if !errors.Is(err, ErrMissingParam) {
		t.Errorf("expected ErrMissingParam, got: %v", err)
	}
}

func TestInsertOnConflict_LastWriteWins(t *testing.T) {
	q := NewInsert().Into("users").Columns("id", "name").
		Values(":id, :name", Params{"id": 1, "name": "Alice"}).
		OnConflictDoNothing("id").
		OnConflictDoUpdate([]string{"id"}, "name = EXCLUDED.name")
	sql, _, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name) VALUES (:id, :name) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertOnConflict_Immutability(t *testing.T) {
	base := NewInsert().Into("users").Columns("id", "name").
		Values(":id, :name", Params{"id": 1, "name": "Alice"})
	_ = base.OnConflictDoNothing("id")
	ib := base.(*insertBuilder)
	if ib.conflict != nil {
		t.Error("base should not have conflict set")
	}
}
