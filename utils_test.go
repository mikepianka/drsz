package drsz

import "testing"

func TestIsCsvPath(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "csv path",
			path:     "some/dir/a.csv",
			expected: true,
		},

		{
			name:     "csv path with uppercase extension",
			path:     "some/dir/a.CSV",
			expected: true,
		},
		{

			name:     "csv path with leading period",
			path:     "some/dir/a.temp.csv",
			expected: true,
		},
		{
			name:     "non-csv path",
			path:     "some/dir/a.txt",
			expected: false,
		},
		{
			name:     "directory path",
			path:     "some/dir",
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if IsCsvPath(tc.path) != tc.expected {
				t.Errorf("expected %s to be identified as csv", tc.path)
			}
		})
	}
}
