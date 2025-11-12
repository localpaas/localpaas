package unitutil

import (
	"testing"
)

func Test_GetSizeString(t *testing.T) {
	tests := []struct {
		size       int64
		roundToInt bool
		expected   string
	}{
		{1023, true, "1023B"},
		{1024, true, "1KB"},
		{1048575, true, "1023KB"},
		{1048576, true, "1MB"},
		{1048576, false, "1MB"},
		{1073741824, true, "1024MB"},
		{1073741824, false, "1024MB"},
	}

	for _, tt := range tests {
		if got := GetSizeString(tt.size, tt.roundToInt); got != tt.expected {
			t.Errorf("GetSizeString(%d, %v) = %v; want %v", tt.size, tt.roundToInt, got, tt.expected)
		}
	}
}
