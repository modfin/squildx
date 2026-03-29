package squildx

import "testing"

func TestDeleteFrom(t *testing.T) {
	b := NewDelete().From("users")
	db := b.(*deleteBuilder)
	if db.table != "users" {
		t.Errorf("table = %q, want %q", db.table, "users")
	}
}

func TestDeleteFrom_Immutability(t *testing.T) {
	base := NewDelete()
	_ = base.From("users")
	db := base.(*deleteBuilder)
	if db.table != "" {
		t.Errorf("base table mutated to %q, want empty", db.table)
	}
}
