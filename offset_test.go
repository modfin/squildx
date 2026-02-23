package squildx

import (
	"testing"
)

func TestOffsetZero(t *testing.T) {
	q, _, err := New().Select("*").From("users").Offset(0).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users OFFSET 0"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}
