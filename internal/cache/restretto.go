package cache

import (
	"errors"

	"github.com/dgraph-io/ristretto/v2"
)

type RistrettoCache[K ristretto.Key, V any] struct {
	cache *ristretto.Cache[K, V]
}

func NewRistrettoCache[K ristretto.Key, V any](maxCost int64, numCounters int64, bufferItems int64) (*RistrettoCache[K, V], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[K, V]{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		return nil, err
	}
	return &RistrettoCache[K, V]{cache: cache}, nil
}

func (r *RistrettoCache[K, V]) Set(key K, value V, cost int64) bool {
	return r.cache.Set(key, value, cost)
}

func (r *RistrettoCache[K, V]) Get(key K) (V, error) {
	value, found := r.cache.Get(key)
	if !found {
		var zeroValue V
		return zeroValue, errors.New("key not found")
	}

	return value, nil
}

func (r *RistrettoCache[K, V]) Delete(key K) {
	r.cache.Del(key)
}

func (r *RistrettoCache[K, V]) GetOrSet(key K, cost int64, fn func() (V, error)) (V, error) {
	value, found := r.cache.Get(key)
	if found {
		return value, nil
	}

	value, err := fn()
	if err != nil {
		return value, err
	}

	r.cache.Set(key, value, cost)
	return value, nil
}
