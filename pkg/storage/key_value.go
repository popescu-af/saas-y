package storage

import "time"

// KeyValue is the interface for a storage holding key-value pairs.
type KeyValue interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration time.Duration) error
	Delete(key string) error
}
