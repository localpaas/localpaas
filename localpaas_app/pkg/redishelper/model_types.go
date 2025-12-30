package redishelper

// Value base interface for storing value in redis
type Value[T any] interface {
	RedisMarshal() ([]byte, error)
	RedisUnmarshal([]byte) error
	GetData() T
}

// ValueCreator creator function to create value when unmarshal data from redis
type ValueCreator[T any] func(T) Value[T]
