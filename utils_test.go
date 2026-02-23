package squildx

import "testing"

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
