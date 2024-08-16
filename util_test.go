package m

import (
	"testing"
)

func TestFirstNonZeroValue(t *testing.T) {
	res := *firstNonZeroValue("", "1")
	if res != "1" {
		t.Errorf("Expected: %v, Got: %v", "1", res)
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		pattern, path string
		expected      bool
	}{
		{"/swagger/*", "/swagger/index.html", true},
		{"/some/*", "/swagger/index.html", false},
		{"/swagger/*/*", "/swagger/some/about.html", true},
		{"/swagger/*", "/swagger/some/about.html", false},
	}

	for _, test := range tests {
		result := matchPattern(test.pattern, test.path)
		if result != test.expected {
			t.Errorf("Pattern: %s, Path: %s - Expected: %t, Got: %t", test.pattern, test.path, test.expected, result)
		}
	}
}
