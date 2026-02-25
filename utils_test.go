package squildx

import "testing"

func TestValueEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"equal ints", 1, 1, true},
		{"different ints", 1, 2, false},
		{"equal strings", "foo", "foo", true},
		{"different strings", "foo", "bar", false},
		{"different types", 1, "1", false},
		{"both nil", nil, nil, true},
		{"nil vs value", nil, 1, false},
		{"equal slices", []int{1, 2}, []int{1, 2}, true},
		{"different slices", []int{1, 2}, []int{1, 3}, false},
		{"equal booleans", true, true, true},
		{"different booleans", true, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := valueEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("valueEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestParamsEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b map[string]any
		want bool
	}{
		{"both nil", nil, nil, true},
		{"both empty", map[string]any{}, map[string]any{}, true},
		{"nil vs empty", nil, map[string]any{}, true},
		{"equal params", map[string]any{"x": 1}, map[string]any{"x": 1}, true},
		{"different values", map[string]any{"x": 1}, map[string]any{"x": 2}, false},
		{"different keys", map[string]any{"x": 1}, map[string]any{"y": 1}, false},
		{"different lengths", map[string]any{"x": 1}, map[string]any{"x": 1, "y": 2}, false},
		{"multiple equal", map[string]any{"a": "foo", "b": 42}, map[string]any{"a": "foo", "b": 42}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := paramsEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("paramsEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestBuildersEqual(t *testing.T) {
	t.Run("identical builders", func(t *testing.T) {
		a := New().Select("*").From("users").Where("id = :id", map[string]any{"id": 1})
		b := New().Select("*").From("users").Where("id = :id", map[string]any{"id": 1})
		if !buildersEqual(a, b) {
			t.Error("expected identical builders to be equal")
		}
	})

	t.Run("different SQL", func(t *testing.T) {
		a := New().Select("*").From("users")
		b := New().Select("*").From("orders")
		if buildersEqual(a, b) {
			t.Error("expected builders with different SQL to not be equal")
		}
	})

	t.Run("different params", func(t *testing.T) {
		a := New().Select("*").From("users").Where("id = :id", map[string]any{"id": 1})
		b := New().Select("*").From("users").Where("id = :id", map[string]any{"id": 2})
		if buildersEqual(a, b) {
			t.Error("expected builders with different params to not be equal")
		}
	})

	t.Run("first builder errors", func(t *testing.T) {
		a := New().Select("*") // missing From
		b := New().Select("*").From("users")
		if buildersEqual(a, b) {
			t.Error("expected false when first builder errors")
		}
	})

	t.Run("second builder errors", func(t *testing.T) {
		a := New().Select("*").From("users")
		b := New().Select("*") // missing From
		if buildersEqual(a, b) {
			t.Error("expected false when second builder errors")
		}
	})

	t.Run("both builders error", func(t *testing.T) {
		a := New().Select("*") // missing From
		b := New().Select("*") // missing From
		if buildersEqual(a, b) {
			t.Error("expected false when both builders error")
		}
	})
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"FirstName", "first_name"},
		{"ID", "id"},
		{"HTTPCode", "http_code"},
		{"UserID", "user_id"},
		{"CreatedAt", "created_at"},
		{"A", "a"},
		{"already", "already"},
		{"HTMLParser", "html_parser"},
		{"myURL", "my_url"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
