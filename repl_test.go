package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{""},
		},
		{
			input:    "1234",
			expected: []string{"1234"},
		},
	}

	for _, c := range cases {
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("error with string length: actual %v != expected %v", len(actual), len(c.expected))
		}

		for i := range actual {
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word %v != expected word %v", word, expectedWord)
			}

		}
	}

}
