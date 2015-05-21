package bidimap

import (
	"fmt"
	"sync"
)

// BiDiMapper interface defines an API available to users.
type BiDiMapper interface {
	Insert(key string, val interface{}) error

	Get(key string) (interface{}, error)
	Key(val interface{}) (string, error) // get key for value

	RemoveByKey(key string) error
	RemoveByVal(val interface{}) error

	IsKeyUsed(key string) bool
	IsMember(val interface{}) bool
}

// BiDiMap implements threadsafe bidirectional map.
type BiDiMap struct {
	direct  map[string]interface{}
	inverse map[interface{}]string

	mu sync.RWMutex
}

// New initializes a new BiDiMap and returns it.
func New() *BiDiMap {
	return &BiDiMap{
		make(map[string]interface{}),
		make(map[interface{}]string),
		sync.RWMutex{},
	}
}

// Helpers
func (b *BiDiMap) unsafeIsKeyUsed(key string) bool {
	_, present := b.direct[key]
	return present
}
func (b *BiDiMap) unsafeIsValUsed(val interface{}) bool {
	_, present := b.inverse[val]
	return present
}

// Insert adds a key-value pair into a bidimap.
func (b *BiDiMap) Insert(key string, val interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.unsafeIsKeyUsed(key) {
		return fmt.Errorf("bidimap: key %q is already used", key)
	}
	if b.unsafeIsValUsed(val) {
		return fmt.Errorf("bidimap: val %q is already used", val)
	}
	b.direct[key] = val
	b.inverse[val] = key
	return nil
}

// Get return value by key.
func (b *BiDiMap) Get(key string) (interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	val, presence := b.direct[key]
	if !presence {
		return nil, fmt.Errorf("bidimap: key %q was not set", key)
	}
	return val, nil
}

// Key return key for a given value.
func (b *BiDiMap) Key(val interface{}) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	key, presence := b.inverse[val]
	if !presence {
		return "", fmt.Errorf("bidimap: val %q was not inserted", val)
	}
	return key, nil
}

// RemoveByKey removes record by a given key.
func (b *BiDiMap) RemoveByKey(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	val, presence := b.direct[key]
	if !presence {
		return fmt.Errorf("bidimap: unable to remove unused key %q", key)
	}
	delete(b.direct, key)
	delete(b.inverse, val)
	return nil
}

// RemoveByVal removes record by a given value.
func (b *BiDiMap) RemoveByVal(val interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	key, presence := b.inverse[val]
	if !presence {
		return fmt.Errorf("bidimap: unable to remove unused val %q", val)
	}
	delete(b.direct, key)
	delete(b.inverse, val)
	return nil
}

// IsKeyUsed returns true if the record for given key exists.
func (b *BiDiMap) IsKeyUsed(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.unsafeIsKeyUsed(key)
}

// IsMember returns true if the given value is already stored.
func (b *BiDiMap) IsMember(val interface{}) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.unsafeIsValUsed(val)
}
