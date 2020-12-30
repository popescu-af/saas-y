package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

// KeyValueMock is a mock for the KeyValue interface.
type KeyValueMock struct {
	mock.Mock

	mutex   sync.RWMutex
	storage map[string][]byte
}

// Get implements the method with the same name from KeyValue.
func (k *KeyValueMock) Get(key string) ([]byte, error) {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	k.Called()
	if v, ok := k.storage[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("key not present")
}

// Set implements the method with the same name from KeyValue.
func (k *KeyValueMock) Set(key string, value []byte, expiration time.Duration) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.Called()
	k.storage[key] = value
	return nil
}

// Delete implements the method with the same name from KeyValue.
func (k *KeyValueMock) Delete(key string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.Called()
	delete(k.storage, key)
	return nil
}

// Ready implements the method with the same name from KeyValue.
func (k *KeyValueMock) Ready() error {
	return nil
}

// NewKeyValueMock creates a KeyValueMock instance.
func NewKeyValueMock() *KeyValueMock {
	return &KeyValueMock{
		storage: make(map[string][]byte),
	}
}
