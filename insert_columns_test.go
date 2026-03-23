package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestInsertColumns(t *testing.T) {
	q := NewInsert().Columns("name", "email")
	ib := q.(*insertBuilder)
	want := []string{"name", "email"}
	if !reflect.DeepEqual(ib.columns, want) {
		t.Errorf("columns = %v, want %v", ib.columns, want)
	}
}

func TestInsertColumns_Chained(t *testing.T) {
	q := NewInsert().Columns("name").Columns("email")
	ib := q.(*insertBuilder)
	want := []string{"name", "email"}
	if !reflect.DeepEqual(ib.columns, want) {
		t.Errorf("columns = %v, want %v", ib.columns, want)
	}
}

func TestInsertColumns_Immutability(t *testing.T) {
	base := NewInsert().Columns("name")
	_ = base.Columns("email")
	ib := base.(*insertBuilder)
	want := []string{"name"}
	if !reflect.DeepEqual(ib.columns, want) {
		t.Errorf("base columns = %v, want %v", ib.columns, want)
	}
}

func TestInsertColumnsObject(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	q := NewInsert().ColumnsObject(User{})
	ib := q.(*insertBuilder)
	want := []string{"name", "email"}
	if !reflect.DeepEqual(ib.columns, want) {
		t.Errorf("columns = %v, want %v", ib.columns, want)
	}
}

func TestInsertColumnsObject_NotAStruct(t *testing.T) {
	q := NewInsert().Into("users").ColumnsObject("not a struct").
		Values(":x", Params{"x": 1})
	_, _, err := q.Build()
	if !errors.Is(err, ErrNotAStruct) {
		t.Errorf("expected ErrNotAStruct, got: %v", err)
	}
}
