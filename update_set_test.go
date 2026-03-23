package squildx

import (
	"errors"
	"testing"
)

func TestUpdateSet(t *testing.T) {
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
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "id", 1)
}

func TestUpdateMultipleSets(t *testing.T) {
	q, params, err := NewUpdate().
		Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Set("age = :age", Params{"age": 30}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name, age = :age WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "age", 30)
	assertParam(t, params, "id", 1)
}

func TestUpdateSetObject(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(User{Name: "Alice", Email: "alice@example.com"}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name, email = :email WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "email", "alice@example.com")
	assertParam(t, params, "id", 1)
}

func TestUpdateSetObject_NilPointerSkipped(t *testing.T) {
	type UserPatch struct {
		Name  *string `db:"name"`
		Email *string `db:"email"`
		Age   *int    `db:"age"`
	}

	name := "Alice"
	patch := UserPatch{Name: &name}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(patch).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "id", 1)
	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}

func TestUpdateSetObject_AllPointersSet(t *testing.T) {
	type UserPatch struct {
		Name  *string `db:"name"`
		Email *string `db:"email"`
	}

	name := "Alice"
	email := "alice@example.com"
	patch := UserPatch{Name: &name, Email: &email}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(patch).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name, email = :email WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "email", "alice@example.com")
}

func TestUpdateSetObject_MixedPointerAndValue(t *testing.T) {
	type UserPatch struct {
		Name  string  `db:"name"`
		Email *string `db:"email"`
	}

	patch := UserPatch{Name: "Alice"}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(patch).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
}

func TestUpdateSetObject_EmbeddedStruct(t *testing.T) {
	type Timestamps struct {
		UpdatedAt string `db:"updated_at"`
	}
	type UserPatch struct {
		Timestamps
		Name *string `db:"name"`
	}

	name := "Alice"
	patch := UserPatch{
		Timestamps: Timestamps{UpdatedAt: "2024-01-01"},
		Name:       &name,
	}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(patch).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET updated_at = :updated_at, name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "updated_at", "2024-01-01")
	assertParam(t, params, "name", "Alice")
}

func TestUpdateSetObject_NotAStruct(t *testing.T) {
	_, _, err := NewUpdate().
		Table("users").
		SetObject("not a struct").
		Where("id = :id", Params{"id": 1}).
		Build()

	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}

func TestUpdateSetObject_SnakeCase(t *testing.T) {
	type User struct {
		FirstName string
		LastName  string
	}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(User{FirstName: "Alice", LastName: "Smith"}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET first_name = :first_name, last_name = :last_name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "first_name", "Alice")
	assertParam(t, params, "last_name", "Smith")
}

func TestUpdateSetAndSetObjectCombined(t *testing.T) {
	type Patch struct {
		Name *string `db:"name"`
	}

	name := "Alice"
	q, params, err := NewUpdate().
		Table("users").
		SetObject(Patch{Name: &name}).
		Set("updated_at = NOW()").
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name, updated_at = NOW() WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "id", 1)
}

func TestUpdateSet_Immutability(t *testing.T) {
	base := NewUpdate().Table("users").
		Set("name = :name", Params{"name": "Alice"}).
		Where("id = :id", Params{"id": 1})

	_ = base.Set("age = :age", Params{"age": 30})

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

func TestUpdateSetObject_SkipTaggedDash(t *testing.T) {
	type User struct {
		Name   string `db:"name"`
		Secret string `db:"-"`
	}

	q, params, err := NewUpdate().
		Table("users").
		SetObject(User{Name: "Alice", Secret: "hidden"}).
		Where("id = :id", Params{"id": 1}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "UPDATE users SET name = :name WHERE id = :id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
	assertParam(t, params, "name", "Alice")
	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}
