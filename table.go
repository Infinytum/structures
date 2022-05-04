package structures

import "errors"

var (
	TableDuplicateKeys = errors.New("the given keys already exist")
	TableKeysNotFound  = errors.New("no value was found for the given keys")
)

// Concurrent table implementation, provides a 3d map
type Table[K1 comparable, K2 comparable, V any] interface {
	// Add stores a new value by the given keys or will error if the keys already exists
	Add(k1 K1, k2 K2, newVal V) error
	// Contains returns whether a value exists for the given keys
	Contains(k1 K1, k2 K2) bool
	// Delete deletes the value by its keys, if the keys does not exist an error will be returned.
	Delete(k1 K1, k2 K2) error
	// Get returns the value by its keys, if the keys does not exist an error will be returned.
	Get(k1 K1, k2 K2) (value V, err error)
	// GetOrDefault returns the value by its keys, if the keys does not exist a given default will be returned
	GetOrDefault(k1 K1, k2 K2, def V) (value V)
	// GetOrSet returns the value by its keys, if the keys does not exist, the given value will be set for them
	GetOrSet(k1 K1, k2 K2, newVal V) (value V)
	// Set stores a new value by the given keys
	Set(k1 K1, k2 K2, newVal V) error
	// ToMap converts the table instance to a native map
	ToMap() map[K1]map[K2]V
}
