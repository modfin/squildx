package squildx

import (
	"errors"
	"testing"
	"time"
)

func TestSelectOnly(t *testing.T) {
	q, params, err := New().Select("*").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q != "SELECT * FROM users" {
		t.Errorf("got: %s", q)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestNoWhereClause(t *testing.T) {
	q, params, err := New().Select("id", "name").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q != "SELECT id, name FROM users" {
		t.Errorf("got: %s", q)
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestErrNoColumns(t *testing.T) {
	_, _, err := New().Select().From("users").Build()
	if !errors.Is(err, ErrNoColumns) {
		t.Errorf("expected ErrNoColumns, got: %v", err)
	}
}

func TestSelectAppends(t *testing.T) {
	q, _, err := New().Select("id").Select("name", "email").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestRemoveSelect(t *testing.T) {
	q, _, err := New().Select("id", "name", "email").RemoveSelect("name").From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestRemoveSelectImmutability(t *testing.T) {
	base := New().Select("id", "name", "email").From("users")
	reduced := base.RemoveSelect("name")

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q2, _, err := reduced.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 != "SELECT id, name, email FROM users" {
		t.Errorf("base was mutated: %s", q1)
	}
	if q2 != "SELECT id, email FROM users" {
		t.Errorf("reduced mismatch: %s", q2)
	}
}

type testUser struct {
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

func TestSelectObjectBasic(t *testing.T) {
	// CreatedAt is a struct (time.Time) so it gets skipped
	q, _, err := New().SelectObject(testUser{}).From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT id, first_name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectPointer(t *testing.T) {
	var u *testUser
	q, _, err := New().SelectObject(u).From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT id, first_name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectTablePrefix(t *testing.T) {
	q, _, err := New().SelectObject(testUser{}, "u").From("users u").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT u.id, u.first_name, u.email FROM users u"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectTagPriority(t *testing.T) {
	type row struct {
		A string `squildx:"col_a" db:"wrong_a" json:"also_wrong_a"`
		B string `db:"col_b" json:"wrong_b"`
		C string `json:"col_c,omitempty"`
		D string
	}
	q, _, err := New().SelectObject(row{}).From("t").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT col_a, col_b, col_c, d FROM t"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectSkipDash(t *testing.T) {
	type row struct {
		ID      int    `db:"id"`
		Secret  string `db:"-"`
		Visible string `db:"visible"`
	}
	q, _, err := New().SelectObject(row{}).From("t").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT id, visible FROM t"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectUnexportedSkipped(t *testing.T) {
	type row struct {
		ID       int `db:"id"`
		internal string
	}
	_ = row{internal: "x"}
	q, _, err := New().SelectObject(row{}).From("t").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT id FROM t"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectNestedStructSkipped(t *testing.T) {
	type Nested struct {
		Val string
	}
	type row struct {
		ID     int    `db:"id"`
		Nested Nested // struct field, should be skipped
	}
	q, _, err := New().SelectObject(row{}).From("t").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT id FROM t"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectSnakeCaseFallback(t *testing.T) {
	type row struct {
		UserID    int
		FirstName string
		HTTPCode  int
	}
	q, _, err := New().SelectObject(row{}).From("t").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT user_id, first_name, http_code FROM t"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectAppendWithSelect(t *testing.T) {
	q, _, err := New().Select("COUNT(*) AS total").SelectObject(testUser{}).From("users").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "SELECT COUNT(*) AS total, id, first_name, email FROM users"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestSelectObjectNonStructError(t *testing.T) {
	_, _, err := New().SelectObject("not a struct").From("t").Build()
	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}

func TestSelectObjectNilError(t *testing.T) {
	_, _, err := New().SelectObject(nil).From("t").Build()
	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}

func TestSelectObjectImmutability(t *testing.T) {
	base := New().Select("extra").From("users")
	withObj := base.SelectObject(testUser{})

	q1, _, err := base.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q2, _, err := withObj.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q1 != "SELECT extra FROM users" {
		t.Errorf("base was mutated: %s", q1)
	}
	expected := "SELECT extra, id, first_name, email FROM users"
	if q2 != expected {
		t.Errorf("withObj mismatch\n got: %s\nwant: %s", q2, expected)
	}
}
