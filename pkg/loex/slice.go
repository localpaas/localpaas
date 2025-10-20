package loex

// MakeSliceWithSkippingNil makes a slice from the given values with skipping `nil` ones
func MakeSliceWithSkippingNil[T any](values ...*T) []*T {
	result := make([]*T, 0, len(values))
	for _, v := range values {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

// MakeSliceWithSkippingZero makes a slice from the given values with skipping `zero` ones
func MakeSliceWithSkippingZero[T comparable](values ...T) []T {
	result := make([]T, 0, len(values))
	var zeroT T
	for _, v := range values {
		if v != zeroT {
			result = append(result, v)
		}
	}
	return result
}
