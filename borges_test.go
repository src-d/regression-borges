package borges

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitURLToProjectPath(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{"git://github.com/foo/bar.git", "github.com/foo/bar"},
		{"git@github.com:foo/bar.git", "github.com/foo/bar"},
		{"git@github.com:foo/bar", "github.com/foo/bar"},
		{"https://github.com/foo/bar.git", "github.com/foo/bar"},
		{"https://github.com/foo/bar", "github.com/foo/bar"},
		{"http://github.com/foo/bar.git", "github.com/foo/bar"},
		{"http://github.com/foo/bar", "github.com/foo/bar"},
		{"yada://github.com/foo/bar.git", "github.com/foo/bar"},
		{"yada://github.com/foo/bar", "github.com/foo/bar"},
	}

	for _, tt := range testCases {
		t.Run(tt.url, func(t *testing.T) {
			result, err := gitURLToProjectPath(tt.url)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
