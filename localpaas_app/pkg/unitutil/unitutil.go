package unitutil

import (
	"fmt"
)

const (
	UnitKB = 1024
	UnitMB = 1024 * UnitKB
)

func GetSizeString(size int64, roundToInt bool) (str string) {
	if size >= UnitMB {
		if roundToInt || size%UnitMB == 0 {
			return fmt.Sprintf("%dMB", size/UnitMB)
		}
		return fmt.Sprintf("%.2fMB", float64(size/UnitMB))
	}
	if size >= UnitKB {
		if roundToInt || size%UnitKB == 0 {
			return fmt.Sprintf("%dKB", size/UnitKB)
		}
		return fmt.Sprintf("%.2fKB", float64(size/UnitKB))
	}
	return fmt.Sprintf("%dB", size)
}
