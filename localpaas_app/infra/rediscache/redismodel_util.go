package rediscache

import (
	"fmt"
)

func ParseString(val any) string {
	str, ok := val.(string)
	if ok {
		return str
	}
	bytes, ok := val.([]byte)
	if ok {
		return string(bytes)
	}
	return fmt.Sprintf("%v", val)
}

func ParseBytes(val any) []byte {
	bytes, ok := val.([]byte)
	if ok {
		return bytes
	}
	str, ok := val.(string)
	if ok {
		return []byte(str)
	}
	return []byte(fmt.Sprintf("%v", val))
}
