package rediscache

import (
	"encoding/json"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

// JSONValue implements `Value[T]` interface to provide a method to store
// value in redis by using json format for the data.
type JSONValue[T any] struct {
	Value[T]
	val T
}

func (m JSONValue[T]) RedisMarshal() ([]byte, error) {
	bytes, err := json.Marshal(m.val)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return bytes, nil
}

func (m *JSONValue[T]) RedisUnmarshal(data []byte) error {
	err := json.Unmarshal(data, &m.val)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (m JSONValue[T]) GetData() T {
	return m.val
}

func NewJSONValue[T any](val T) Value[T] {
	return &JSONValue[T]{val: val}
}
