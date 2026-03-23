package squildx

import "testing"

func TestInsertSelect(t *testing.T) {
	sub := New().Select("name", "email").From("temp_users").Where("active = :active", Params{"active": true})
	q := NewInsert().Into("users").Columns("name", "email").Select(sub)
	sql, params, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "INSERT INTO users (name, email) SELECT name, email FROM temp_users WHERE active = :active"
	if sql != want {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", sql, want)
	}
	assertParam(t, params, "active", true)
}

func TestInsertSelect_Immutability(t *testing.T) {
	sub := New().Select("name").From("temp")
	base := NewInsert().Into("users").Columns("name")
	_ = base.Select(sub)
	ib := base.(*insertBuilder)
	if ib.selectQuery != nil {
		t.Error("base should not have selectQuery set")
	}
}
