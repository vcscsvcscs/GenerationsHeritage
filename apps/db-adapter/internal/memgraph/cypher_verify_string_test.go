package memgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyString(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"MATCH (n) RETURN n", true}, // Contains Cypher keyword
		{"Hello 'World'", true},      // Contains Cypher delimiter
		{"Hello World", false},       // Valid string
	}

	for _, test := range tests {
		err := VerifyString(test.input)
		if test.expectError {
			require.Error(t, err, "expected error for input: %s", test.input)
		} else {
			require.NoError(t, err, "did not expect error for input: %s", test.input)
		}
	}
}

func TestEscapeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello 'World'", "Hello \\'World\\'"},
		{"Hello \"World\"", "Hello \\\"World\\\""},
		{"Hello `World`", "Hello ``World``"},
		{"Hello World", "Hello World"}, // No delimiters to escape
	}

	for _, test := range tests {
		result := EscapeString(test.input)
		assert.Equal(t, test.expected, result, "unexpected result for input: %s", test.input)
	}
}
