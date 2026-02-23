package squildx

import (
	"testing"
)

func TestLimitZero(t *testing.T) {
	q, _, err := New().Select("*").From("users").Limit(0).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users LIMIT 0"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
