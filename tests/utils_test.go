package tests

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already_normalized_path",
			input:    "./some/path",
			expected: "./some/path",
		},
		{
			name:     "path_without_prefix",
			input:    "some/path",
			expected: "./some/path",
		},
		{
			name:     "parent_directory_path",
			input:    "../some/path",
			expected: "../some/path",
		},
		{
			name:     "absolute_path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "current_directory",
			input:    ".",
			expected: ".",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "double_dot_path",
			input:    "../../some/path",
			expected: "../../some/path",
		},
		{
			name:     "clean_double_slash",
			input:    "some//path",
			expected: "./some/path",
		},
		{
			name:     "clean_current_dir_references",
			input:    "some/./path",
			expected: "./some/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
