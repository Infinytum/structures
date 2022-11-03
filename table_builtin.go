package structures

import (
	"errors"
	"sync"
)

// Table is a wrapper for a 3 dimensional map that manages its contents
type builtinTable[K1 comparable, K2 comparable, V2 any] struct {
	table  map[K1]map[K2]V2
	shadow Table[K1, K2, V2]
	sync.RWMutex
}

// Add stores a new value by the given keys or will error if the keys already exist
func (t *builtinTable[K1, K2, V2]) Add(k1 K1, k2 K2, newVal V2) error {
	t.Lock()
	defer t.Unlock()

	if _, exists := t.table[k1]; !exists {
		t.table[k1] = make(map[K2]V2)
	} else if _, exists := t.table[k1][k2]; exists {
		return TableDuplicateKeys
	}
	t.table[k1][k2] = newVal

	return nil
}

// Contains returns whether a value exists for the given keys
func (t *builtinTable[K1, K2, V2]) Contains(k1 K1, k2 K2) bool {
	t.RLock()
	defer t.RUnlock()

	if table2, exists := t.table[k1]; exists {
		if _, exists := table2[k2]; exists {
			return true
		}
	}

	if t.shadow != nil {
		return t.shadow.Contains(k1, k2)
	}
	return false
}

// Delete deletes the value by its keys, if the keys does not exist an error will be returned.
func (t *builtinTable[K1, K2, V2]) Delete(k1 K1, k2 K2) error {
	t.Lock()
	defer t.Unlock()

	if table2, exists := t.table[k1]; exists {
		if _, exists := table2[k2]; exists {
			delete(t.table[k1], k2)
			return nil
		}
	}

	if t.shadow != nil && t.shadow.Contains(k1, k2) {
		return errors.New("cannot delete from shadow table")
	}

	return TableKeysNotFound
}

// Get returns the value by its keys, if the keys does not exist an error will be returned.
func (t *builtinTable[K1, K2, V2]) Get(k1 K1, k2 K2) (value V2, err error) {
	t.RLock()
	defer t.RUnlock()

	err = TableKeysNotFound
	if table2, exists := t.table[k1]; exists {
		if val, exists := table2[k2]; exists {
			value = val
			err = nil
		}
	}

	if err != nil && t.shadow != nil {
		return t.shadow.Get(k1, k2)
	}

	return
}

// GetOrDefault returns the value by its keys, if the keys does not exist a given default will be returned
func (t *builtinTable[K1, K2, V2]) GetOrDefault(k1 K1, k2 K2, def V2) (value V2) {
	t.RLock()
	defer t.RUnlock()

	if table2, exists := t.table[k1]; exists {
		if val, exists := table2[k2]; exists {
			return val
		}
	}

	if t.shadow != nil {
		return t.shadow.GetOrDefault(k1, k2, def)
	}

	return def
}

// GetOrSet returns the value by its keys, if the keys does not exist, the given value will be set for them
func (t *builtinTable[K1, K2, V2]) GetOrSet(k1 K1, k2 K2, newVal V2) (value V2) {
	t.Lock()
	defer t.Unlock()

	if table2, exists := t.table[k1]; exists {
		if val, exists := table2[k2]; exists {
			return val
		} else if t.shadow != nil && t.shadow.Contains(k1, k2) {
			if v, err := t.shadow.Get(k1, k2); err == nil {
				return v
			}
		} else {
			table2[k2] = newVal
			return newVal
		}
	} else if t.shadow != nil && t.shadow.Contains(k1, k2) {
		if v, err := t.shadow.Get(k1, k2); err == nil {
			return v
		}
	} else {
		t.table[k1] = make(map[K2]V2)
		t.table[k1][k2] = newVal
		return newVal
	}

	return
}

// Set stores a new value by the given keys
func (t *builtinTable[K1, K2, V2]) Set(k1 K1, k2 K2, newVal V2) error {
	t.Lock()
	defer t.Unlock()

	if _, exists := t.table[k1]; !exists {
		t.table[k1] = make(map[K2]V2)
	}
	t.table[k1][k2] = newVal

	return nil
}

// ShadowCopy returns a new table with the current table as its shadow
func (t *builtinTable[K1, K2, V]) ShadowCopy() Table[K1, K2, V] {
	newTable := NewTable[K1, K2, V]()
	newTable.(*builtinTable[K1, K2, V]).shadow = t
	return newTable
}

// ToMap converts the map instance to a native map
func (t *builtinTable[K1, K2, V]) ToMap() map[K1]map[K2]V {
	t.RLock()
	defer t.RUnlock()

	if t.shadow == nil {
		return t.table
	}

	m := NewTable[K1, K2, V]()
	for k1, table2 := range t.shadow.ToMap() {
		for k2, v := range table2 {
			m.Set(k1, k2, v)
		}
	}
	for k1, table2 := range t.table {
		for k2, v := range table2 {
			m.Set(k1, k2, v)
		}
	}

	return m.ToMap()
}

// NewTable will create a new, empty instance of Table
func NewTable[K1 comparable, K2 comparable, V any]() Table[K1, K2, V] {
	return &builtinTable[K1, K2, V]{
		table:   make(map[K1]map[K2]V),
		RWMutex: sync.RWMutex{},
	}
}
