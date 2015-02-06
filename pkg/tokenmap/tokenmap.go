package tokenmap

import (
	"fmt"
	"sync"

	"gatewayd/utils"
)

// TokenMap implements threadsafe bidirectional map.
// It is a lot like bidimap, but extended for working with
// random generated tokens.
// Code is mostly copy-pasted. No inheritance in go. :(
//                     ^--- this may be the wrong solution though.
type TokenMap struct {
	direct  map[string]interface{}
	inverse map[interface{}]string

	mu sync.RWMutex
}

// New initializes a new BiDiMap and returns it.
func New() *TokenMap {
	return &TokenMap{
		make(map[string]interface{}),
		make(map[interface{}]string),
		sync.RWMutex{},
	}
}

// Helpers
func (b *TokenMap) unsafeIsKeyUsed(key string) bool {
	_, present := b.direct[key]
	return present
}
func (b *TokenMap) unsafeIsValUsed(val interface{}) bool {
	_, present := b.inverse[val]
	return present
}

// Insert adds a key-value pair into a bidimap.
func (b *TokenMap) Insert(key string, val interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.unsafeIsKeyUsed(key) {
		return fmt.Errorf("tokenmap: key %q is already used", key)
	}
	if b.unsafeIsValUsed(val) {
		return fmt.Errorf("tokenmap: val %q is already used", val)
	}
	b.direct[key] = val
	b.inverse[val] = key
	return nil
}

// Get return value by key.
func (b *TokenMap) Get(key string) (interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	val, presence := b.direct[key]
	if !presence {
		return nil, fmt.Errorf("tokenmap: key %q was not set", key)
	}
	return val, nil
}

// Key return key for a given value.
func (b *TokenMap) Key(val interface{}) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	key, presence := b.inverse[val]
	if !presence {
		return "", fmt.Errorf("tokenmap: val %q was not inserted", val)
	}
	return key, nil
}

// RemoveByKey removes record by a given key.
func (b *TokenMap) RemoveByKey(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	val, presence := b.direct[key]
	if !presence {
		return fmt.Errorf("tokenmap: unable to remove unused key %q", key)
	}
	delete(b.direct, key)
	delete(b.inverse, val)
	return nil
}

// RemoveByVal removes record by a given value.
func (b *TokenMap) RemoveByVal(val interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	key, presence := b.inverse[val]
	if !presence {
		return fmt.Errorf("tokenmap: unable to remove unused val %q", val)
	}
	delete(b.direct, key)
	delete(b.inverse, val)
	return nil
}

// IsKeyUsed returns true if the record for given key exists.
func (b *TokenMap) IsKeyUsed(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.unsafeIsKeyUsed(key)
}

// IsMember returns true if the given value is already stored.
func (b *TokenMap) IsMember(val interface{}) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.unsafeIsValUsed(val)
}

// InsertWithRandomKey inserts specified value assigning
// it unique random key of 32 symbols length (or error).
func (b *TokenMap) InsertWithRandomKey(val interface{}) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.unsafeIsValUsed(val) {
		return "", fmt.Errorf("tokenmap: val %q is already used", val)
	}
	for attempts := 0; attempts < 10000; attempts++ {
		key := utils.RandStr(32)
		if !b.unsafeIsKeyUsed(key) {
			b.direct[key] = val
			b.inverse[val] = key
			return key, nil
		}
	}
	return "", fmt.Errorf("tokenmap: too many attempts to generate key for %q", val)
}
