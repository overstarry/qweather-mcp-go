package utils

import (
	"testing"
)

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		sep      string
		expected string
	}{
		{
			name:     "join with comma separator",
			strs:     []string{"a", "b", "c"},
			sep:      ",",
			expected: "a,b,c",
		},
		{
			name:     "filter empty strings",
			strs:     []string{"a", "", "b", "", "c"},
			sep:      ",",
			expected: "a,b,c",
		},
		{
			name:     "all empty strings",
			strs:     []string{"", "", ""},
			sep:      ",",
			expected: "",
		},
		{
			name:     "empty slice",
			strs:     []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "single non-empty string",
			strs:     []string{"hello"},
			sep:      ",",
			expected: "hello",
		},
		{
			name:     "single empty string",
			strs:     []string{""},
			sep:      ",",
			expected: "",
		},
		{
			name:     "different separator",
			strs:     []string{"a", "b", "c"},
			sep:      " | ",
			expected: "a | b | c",
		},
		{
			name:     "mixed empty and non-empty with different separator",
			strs:     []string{"hello", "", "world", ""},
			sep:      " ",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinStrings(tt.strs, tt.sep)
			if result != tt.expected {
				t.Errorf("JoinStrings(%v, %q) = %q, want %q", tt.strs, tt.sep, result, tt.expected)
			}
		})
	}
}
