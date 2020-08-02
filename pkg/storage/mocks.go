package storage

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

// KeyValueMock is a mock for the KeyValue interface.
type KeyValueMock struct {
	mock.Mock

	storage map[string][]byte
}

// Get implements the method with the same name from ChannelListener.
func (k *KeyValueMock) Get(key string) ([]byte, error) {
	k.Called()
	if v, ok := k.storage[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("key not present")
}

// Set implements the method with the same name from ChannelListener.
func (k *KeyValueMock) Set(key string, value []byte, expiration time.Duration) error {
	k.Called()
	k.storage[key] = value
	return nil
}

// Delete implements the method with the same name from ChannelListener.
func (k *KeyValueMock) Delete(key string) error {
	k.Called()
	delete(k.storage, key)
	return nil
}

// NewKeyValueMock creates a KeyValueMock instance.
func NewKeyValueMock() *KeyValueMock {
	return &KeyValueMock{
		storage: make(map[string][]byte),
	}
}
