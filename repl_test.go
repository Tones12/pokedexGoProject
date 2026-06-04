package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name	 string
		input    string
		expected []string
	}{
		{
			name:	  "whitespace and caps test",
			input:    "  Hello  World  ",
			expected: []string{"hello", "world"},
		},
		{
			name:	  "just white space test",
			input:    "   ",
			expected: []string{},
		},
		{
			name:	  "one word test",
			input:    "1234",
			expected: []string{"1234"},
		},
		{
			name:	  "multiple single character test",
			input:    "a b c d e f g",
			expected: []string{"a", "b", "c", "d", "e", "f", "g"},
		},
	}

	for _, c := range cases {
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Test: %v\nerror with string length: actual %v != expected %v", c.name, len(actual), len(c.expected))
		}

		for i := range actual {
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test: %v\nerror: word %v != expected word %v", c.name, word, expectedWord)
			}

		}
	}

}
