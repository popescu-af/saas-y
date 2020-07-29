package storage

import "time"

// KeyValueStore is the interface for a storage holding key-value pairs.
type KeyValueStore interface {
	Get(key string) (string, error)
	Set(key string, value []byte, expiration time.Duration) error
	Delete(key string) error
}
