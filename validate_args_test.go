package main

import (
	"log"
	"os"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		stdinData   string
		expected    []int
		expectError bool
		errorMsg    string
	}{
		// Args only scenarios
		{
			name:     "args only - single number",
			args:     []string{"5"},
			expected: []int{5},
		},
		{
			name:     "args only - multiple numbers",
			args:     []string{"1", "2", "3", "4", "5"},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "args only - space separated in single arg",
			args:     []string{"1 2 3", "4 5"},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "args only - floats converted to ints",
			args:     []string{"1.7", "2.3", "3.9"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "args only - mixed integers and floats",
			args:     []string{"1", "2.5", "3"},
			expected: []int{1, 2, 3},
		},

		// Stdin only scenarios
		{
			name:      "stdin only - single line single number",
			args:      []string{},
			stdinData: "42",
			expected:  []int{42},
		},
		{
			name:      "stdin only - single line multiple numbers",
			args:      []string{},
			stdinData: "1 2 3 4 5",
			expected:  []int{1, 2, 3, 4, 5},
		},
		{
			name:      "stdin only - multiple lines",
			args:      []string{},
			stdinData: "1 2\n3 4\n5",
			expected:  []int{1, 2, 3, 4, 5},
		},
		{
			name:      "stdin only - floats",
			args:      []string{},
			stdinData: "1.1 2.9 3.5",
			expected:  []int{1, 2, 3},
		},
		{
			name:      "stdin only - extra whitespace",
			args:      []string{},
			stdinData: "  1   2    3  ",
			expected:  []int{1, 2, 3},
		},
		{
			name:      "stdin only - tabs and mixed whitespace",
			args:      []string{},
			stdinData: "1\t2   3\n4",
			expected:  []int{1, 2, 3, 4},
		},

		// Args + stdin precedence (args should win)
		{
			name:      "args precedence - args take priority over stdin",
			args:      []string{"10", "20"},
			stdinData: "1 2 3",
			expected:  []int{10, 20},
		},
		{
			name:      "args precedence - single arg vs stdin",
			args:      []string{"99"},
			stdinData: "1 2 3 4 5",
			expected:  []int{99},
		},

		// Error cases
		{
			name:        "no input - empty args and no stdin",
			args:        []string{},
			stdinData:   "",
			expectError: true,
			errorMsg:    "requires at least 1 argument",
		},
		{
			name:        "invalid number in args",
			args:        []string{"1", "abc", "3"},
			expectError: true,
			errorMsg:    "abc is not a number",
		},
		{
			name:        "invalid number in stdin",
			args:        []string{},
			stdinData:   "1 xyz 3",
			expectError: true,
			errorMsg:    "xyz is not a number",
		},
		{
			name:        "only whitespace in stdin",
			args:        []string{},
			stdinData:   "   \n   \n",
			expectError: true,
			errorMsg:    "requires at least 1 argument",
		},

		// Edge cases
		{
			name:     "negative numbers",
			args:     []string{"-1", "-5", "10"},
			expected: []int{-1, -5, 10},
		},
		{
			name:     "zero values",
			args:     []string{"0", "1", "0", "2"},
			expected: []int{0, 1, 0, 2},
		},
		{
			name:      "stdin with negative numbers",
			args:      []string{},
			stdinData: "-5 -1 0 1 5",
			expected:  []int{-5, -1, 0, 1, 5},
		},
		{
			name:     "large numbers",
			args:     []string{"999999", "1000000"},
			expected: []int{999999, 1000000},
		},
		{
			name:      "stdin with decimal precision loss",
			args:      []string{},
			stdinData: "1.9999 2.0001",
			expected:  []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file to simulate stdin
			var file *os.File
			if tt.stdinData != "" {
				tmpFile, err := os.CreateTemp("", "test_stdin")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				defer func(name string) {
					if err := os.Remove(name); err != nil {
						log.Fatalf("failed to remove temp file: %v", err)
					}
				}(tmpFile.Name())
				defer func(tmpFile *os.File) {
					if err := tmpFile.Close(); err != nil {
						log.Fatalf("failed to close temp file: %v", err)
					}
				}(tmpFile)

				// Write test data and reset file position
				if _, err := tmpFile.WriteString(tt.stdinData); err != nil {
					t.Fatalf("failed to write to temp file: %v", err)
				}
				if _, err := tmpFile.Seek(0, 0); err != nil {
					t.Fatalf("failed to seek temp file: %v", err)
				}
				file = tmpFile
			} else {
				// Create a regular character device file (simulates no piped data)
				file = os.Stdin
			}

			result, err := validateArgs(tt.args, file)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !sliceEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to check if two int slices are equal
func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
