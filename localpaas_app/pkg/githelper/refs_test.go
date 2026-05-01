package githelper

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeRepoRef(t *testing.T) {
	tests := []struct {
		input    string
		expected plumbing.ReferenceName
	}{
		{"", "HEAD"},
		{"refs/heads/main", "refs/heads/main"},
		{"tags/v1.0", plumbing.NewTagReferenceName("v1.0")},
		{"heads/feature", plumbing.NewBranchReferenceName("feature")},
		{"somebranch", plumbing.NewBranchReferenceName("somebranch")},
	}

	for _, tt := range tests {
		result := NormalizeRepoRef(tt.input)
		assert.Equal(t, tt.expected, result, "NormalizeRepoRef(%s)", tt.input)
	}
}
