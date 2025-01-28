package utils

import "testing"

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "empty string",
			email:    "",
			expected: false,
		},
		{
			name:     "@ at start",
			email:    "@invalid.com",
			expected: false,
		},
		{
			name:     "no @",
			email:    "invalidemail.com",
			expected: false,
		},
		{
			name:     "multiple @",
			email:    "test@multiple@example.com",
			expected: true,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			got := IsEmail(ts.email)
			if got != ts.expected {
				t.Errorf("IsEmail(%q) = %v, want %v", ts.email, got, ts.expected)
			}
		})
	}
}
