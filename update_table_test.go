package squildx

import "testing"

func TestUpdateTable(t *testing.T) {
	b := NewUpdate().Table("users")
	ub := b.(*updateBuilder)
	if ub.table != "users" {
		t.Errorf("table = %q, want %q", ub.table, "users")
	}
}

func TestUpdateTable_Immutability(t *testing.T) {
	base := NewUpdate()
	_ = base.Table("users")
	ub := base.(*updateBuilder)
	if ub.table != "" {
		t.Errorf("base table mutated to %q, want empty", ub.table)
	}
}
