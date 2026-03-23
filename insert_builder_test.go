package squildx

import "testing"

func TestInsertBuilder_SimpleInsert(t *testing.T) {
	sql, params, err := NewInsert().
		Into("users").
		Columns("name", "email").
		Values(":name, :email", Params{"name": "Alice", "email": "a@b.com"}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) VALUES (:name, :email)"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "email", "a@b.com")
}

func TestInsertBuilder_MultiRow(t *testing.T) {
	sql, params, err := NewInsert().
		Into("users").
		Columns("name", "email").
		Values(":n1, :e1", Params{"n1": "Alice", "e1": "a@b.com"}).
		Values(":n2, :e2", Params{"n2": "Bob", "e2": "b@b.com"}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) VALUES (:n1, :e1), (:n2, :e2)"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "n1", "Alice")
	assertParam(t, params, "e1", "a@b.com")
	assertParam(t, params, "n2", "Bob")
	assertParam(t, params, "e2", "b@b.com")
}

func TestInsertBuilder_StructInsert(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	sql, params, err := NewInsert().
		Into("users").
		ValuesObject(User{Name: "Alice", Email: "a@b.com"}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) VALUES (:name, :email)"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "email", "a@b.com")
}

func TestInsertBuilder_InsertSelect(t *testing.T) {
	sub := New().Select("name", "email").From("temp_users").Where("active = :active", Params{"active": true})
	sql, params, err := NewInsert().
		Into("users").
		Columns("name", "email").
		Select(sub).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) SELECT name, email FROM temp_users WHERE active = :active"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "active", true)
}

func TestInsertBuilder_OnConflictDoNothing(t *testing.T) {
	sql, _, err := NewInsert().
		Into("users").
		Columns("id", "name", "email").
		Values(":id, :name, :email", Params{"id": 1, "name": "Alice", "email": "a@b.com"}).
		OnConflictDoNothing("id").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name, email) VALUES (:id, :name, :email) ON CONFLICT (id) DO NOTHING"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertBuilder_OnConflictDoUpdate(t *testing.T) {
	sql, _, err := NewInsert().
		Into("users").
		Columns("id", "name", "email").
		Values(":id, :name, :email", Params{"id": 1, "name": "Alice", "email": "a@b.com"}).
		OnConflictDoUpdate([]string{"id"}, "name = EXCLUDED.name, email = EXCLUDED.email").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name, email) VALUES (:id, :name, :email) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertBuilder_Returning(t *testing.T) {
	sql, _, err := NewInsert().
		Into("users").
		Columns("name", "email").
		Values(":name, :email", Params{"name": "Alice", "email": "a@b.com"}).
		Returning("id", "created_at").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) VALUES (:name, :email) RETURNING id, created_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertBuilder_FullCombination(t *testing.T) {
	sql, params, err := NewInsert().
		Into("users").
		Columns("id", "name", "email").
		Values(":id, :name, :email", Params{"id": 1, "name": "Alice", "email": "a@b.com"}).
		OnConflictDoUpdate([]string{"id"}, "name = EXCLUDED.name, updated_at = :now", Params{"now": "2024-01-01"}).
		Returning("id", "updated_at").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (id, name, email) VALUES (:id, :name, :email) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, updated_at = :now RETURNING id, updated_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "id", 1)
	assertParam(t, params, "name", "Alice")
	assertParam(t, params, "email", "a@b.com")
	assertParam(t, params, "now", "2024-01-01")
}

func TestInsertBuilder_Immutability(t *testing.T) {
	base := NewInsert().Into("users").Columns("name").
		Values(":name", Params{"name": "Alice"})

	withReturning := base.Returning("id")
	withConflict := base.OnConflictDoNothing("name")

	sql1, _, err1 := base.Build()
	sql2, _, err2 := withReturning.Build()
	sql3, _, err3 := withConflict.Build()

	if err1 != nil {
		t.Fatalf("base build error: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("withReturning build error: %v", err2)
	}
	if err3 != nil {
		t.Fatalf("withConflict build error: %v", err3)
	}

	if sql1 != "INSERT INTO users (name) VALUES (:name)" {
		t.Errorf("base SQL unexpected: %s", sql1)
	}
	if sql2 != "INSERT INTO users (name) VALUES (:name) RETURNING id" {
		t.Errorf("withReturning SQL unexpected: %s", sql2)
	}
	if sql3 != "INSERT INTO users (name) VALUES (:name) ON CONFLICT (name) DO NOTHING" {
		t.Errorf("withConflict SQL unexpected: %s", sql3)
	}
}

func TestInsertBuilder_ReturningObject(t *testing.T) {
	type Result struct {
		ID        int    `db:"id"`
		CreatedAt string `db:"created_at"`
	}
	sql, _, err := NewInsert().
		Into("users").
		Columns("name").
		Values(":name", Params{"name": "Alice"}).
		ReturningObject(Result{}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name) VALUES (:name) RETURNING id, created_at"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertBuilder_ColumnsObject(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	sql, _, err := NewInsert().
		Into("users").
		ColumnsObject(User{}).
		Values(":name, :email", Params{"name": "Alice", "email": "a@b.com"}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) VALUES (:name, :email)"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
}

func TestInsertBuilder_AtPrefix(t *testing.T) {
	sql, params, err := NewInsert().
		Into("users").
		Columns("name").
		Values("@name", Params{"name": "Alice"}).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name) VALUES (@name)"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "name", "Alice")
}
