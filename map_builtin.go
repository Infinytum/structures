package structures

import (
	"errors"
	"sync"
)

// Map is a wrapper for a 2 dimensional map that manages its contents
type builtinMap[K comparable, V any] struct {
	table  map[K]V
	shadow Map[K, V]
	sync.RWMutex
}

// Add stores a new value by the given keys or will error if the keys already exist
func (t *builtinMap[K, V]) Add(key K, newVal V) error {
	t.Lock()
	defer t.Unlock()

	if _, exists := t.table[key]; exists {
		return TableDuplicateKeys
	}
	t.table[key] = newVal
	return nil
}

// Contains returns whether a value exists for the given keys
func (t *builtinMap[K, V]) Contains(key K) bool {
	t.RLock()
	defer t.RUnlock()

	if _, exists := t.table[key]; exists {
		return true
	}

	if t.shadow != nil {
		return t.shadow.Contains(key)
	}
	return false
}

// Delete deletes the value by its keys, if the keys does not exist an error will be returned.
func (t *builtinMap[K, V]) Delete(key K) error {
	t.Lock()
	defer t.Unlock()

	if _, exists := t.table[key]; exists {
		delete(t.table, key)
		return nil
	}

	if t.shadow != nil && t.shadow.Contains(key) {
		return errors.New("cannot delete from shadow table")
	}

	return TableKeysNotFound
}

// Get returns the value by its keys, if the keys does not exist an error will be returned.
func (t *builtinMap[K, V]) Get(key K) (value V, err error) {
	t.RLock()
	defer t.RUnlock()

	err = TableKeysNotFound
	if val, exists := t.table[key]; exists {
		value = val
		err = nil
	}

	if err != nil && t.shadow != nil {
		return t.shadow.Get(key)
	}

	return
}

// GetOrDefault returns the value by its keys, if the keys does not exist a given default will be returned
func (t *builtinMap[K, V]) GetOrDefault(key K, def V) (value V) {
	t.RLock()
	defer t.RUnlock()

	if val, exists := t.table[key]; exists {
		return val
	}

	if t.shadow != nil {
		return t.shadow.GetOrDefault(key, def)
	}

	return def
}

// GetOrSet returns the value by its keys, if the keys does not exist, the given value will be set for them
func (t *builtinMap[K, V]) GetOrSet(key K, newVal V) (value V) {
	t.Lock()
	defer t.Unlock()

	if v, exists := t.table[key]; exists {
		value = v
	} else if t.shadow != nil {
		if v, err := t.shadow.Get(key); err == nil {
			value = v
		}
	} else {
		t.table[key] = newVal
		value = newVal
	}
	return
}

// Set stores val as a new value by key
func (t *builtinMap[K, V]) Set(key K, newVal V) error {
	t.Lock()
	defer t.Unlock()

	t.table[key] = newVal
	return nil
}

// ShadowCopy returns a new map with the current map as its shadow
func (t *builtinMap[K, V]) ShadowCopy() Map[K, V] {
	newMap := NewMap[K, V]()
	newMap.(*builtinMap[K, V]).shadow = t
	return newMap
}

// ToMap converts the map instance to a native map
func (t *builtinMap[K, V]) ToMap() map[K]V {
	t.RLock()
	defer t.RUnlock()

	if t.shadow == nil {
		return t.table
	}

	m := make(map[K]V)
	for k, v := range t.table {
		m[k] = v
	}

	for k, v := range t.shadow.ToMap() {
		m[k] = v
	}
	return m
}

// NewMap will create a new, empty instance of Map
func NewMap[K comparable, V any]() Map[K, V] {
	return &builtinMap[K, V]{
		table:   make(map[K]V),
		RWMutex: sync.RWMutex{},
	}
}
