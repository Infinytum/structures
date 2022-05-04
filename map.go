package structures

import "errors"

var (
	MapKeyNotFound  = errors.New("key was not found")
	MapDuplicateKey = errors.New("key already exists")
)

// Concurrent map implementation
type Map[K comparable, V any] interface {
	// Add stores a new value by the given key or will error if the key already exists
	Add(key K, val V) error
	// Contains returns a whether the given key is contained in the hashmap or not.
	Contains(key K) bool
	// Delete deletes a value by key, if key does not exist an error will be returned.
	Delete(key K) error
	// Get returns the value by key, if key does not exist an error will be returned.
	Get(key K) (V, error)
	// GetOrDefault returns the value by key, if key does not exist def will be returned.
	GetOrDefault(key K, def V) V
	// GetOrDefault returns the value by key, if key does not exist def will be returned and stored.
	GetOrSet(key K, def V) V
	// Set stores val as a new value by key
	Set(key K, val V) error
	// ToMap converts the map instance to a native map
	ToMap() map[K]V
}
