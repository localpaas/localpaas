package fileutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSubpath(t *testing.T) {
	tests := []struct {
		base     string
		target   string
		expected bool
		wantErr  bool
	}{
		{"/a/b", "/a/b/c", true, false},
		{"/a/b", "/a/b", false, false}, // IsSubpath should be false for equal paths
		{"/a/b", "/a/x", false, false},
		{"/a/b", "/a/b/../x", false, false},
		{"a/b", "a/b/c", true, false},
		{"a/b", "a/x", false, false},
		{"/a/b", "a/b/c", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.base+"_"+tt.target, func(t *testing.T) {
			got, err := IsSubpath(tt.base, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestIsEqualOrSubpath(t *testing.T) {
	tests := []struct {
		base     string
		target   string
		expected bool
		wantErr  bool
	}{
		{"/a/b", "/a/b/c", true, false},
		{"/a/b", "/a/b", true, false}, // IsEqualOrSubpath should be true for equal paths
		{"/a/b", "/a/x", false, false},
		{"/a/b", "/a/b/../x", false, false},
		{"a/b", "a/b/c", true, false},
		{"a/b", "a/x", false, false},
		{"/a/b", "a/b/c", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.base+"_"+tt.target, func(t *testing.T) {
			got, err := IsEqualOrSubpath(tt.base, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
