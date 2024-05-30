package maputil

import "testing"

func TestAnyKeys(T *testing.T) {
	testCases := []struct {
		name     string
		m        interface{}
		keys     []string
		expected bool
	}{
		{
			name:     "map contains all keys",
			m:        map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "map contains extra keys",
			m:        map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "b"},
			expected: false,
		},
		{
			name:     "map contains no keys",
			m:        map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"x", "y"},
			expected: false,
		},
		{
			name:     "empty map and empty keys",
			m:        map[string]int{},
			keys:     []string{},
			expected: true,
		},
		{
			name:     "empty map",
			m:        map[string]int{},
			keys:     []string{"a", "b"},
			expected: false,
		},
		{
			name:     "at least one keys",
			m:        map[string]int{"a": 1},
			keys:     []string{"a", "b", "c", "d"},
			expected: true,
		},
		{
			name:     "non-map input",
			m:        123,
			keys:     []string{"a", "b"},
			expected: false,
		},
	}
	for _, tc := range testCases {
		T.Run(tc.name, func(T *testing.T) {
			actual := AnyKeys(tc.m, tc.keys...)
			if actual != tc.expected {
				T.Errorf("AnyKeys(%v,%v)=%v, expected %v", tc.m, tc.keys, actual, tc.expected)
			}
		})
	}
}
