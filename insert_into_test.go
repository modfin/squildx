package squildx

import "testing"

func TestInsertInto(t *testing.T) {
	q := NewInsert().Into("users")
	ib := q.(*insertBuilder)
	if ib.table != "users" {
		t.Errorf("expected table %q, got %q", "users", ib.table)
	}
}

func TestInsertInto_Immutability(t *testing.T) {
	base := NewInsert()
	_ = base.Into("users")
	ib := base.(*insertBuilder)
	if ib.table != "" {
		t.Errorf("expected base table to be empty, got %q", ib.table)
	}
}

func TestInsertInto_Override(t *testing.T) {
	q := NewInsert().Into("users").Into("orders")
	ib := q.(*insertBuilder)
	if ib.table != "orders" {
		t.Errorf("expected table %q, got %q", "orders", ib.table)
	}
}
