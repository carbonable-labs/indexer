package storage

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cockroachdb/pebble"
)

type PebbleStorageOptsFunc func(*PebbleStorageOptions)

type PebbleStorageOptions struct {
	path string
}

func defaultPebbleOptions() *PebbleStorageOptions {
	return &PebbleStorageOptions{
		path: "sheshat/pebble_storage",
	}
}

type PebbleStorage struct {
	handle *pebble.DB
}

func (p *PebbleStorage) Get(id []byte) []byte {
	value, closer, err := p.handle.Get(id)
	if err != nil {
		log.Debug(err)
		return []byte("")
	}
	_ = closer.Close()

	return value
}

func (p *PebbleStorage) Has(id []byte) bool {
	_, closer, err := p.handle.Get(id)
	if err != nil {
		log.Debug(err)
		return false
	}

	_ = closer.Close()

	return true
}

func (p *PebbleStorage) Set(id []byte, value []byte) error {
	if err := p.handle.Set(id, value, pebble.Sync); err != nil {
		log.Error(err)
		return fmt.Errorf("failed to set value at key : %s (%s)", string(id), err)
	}

	return nil
}

func (p *PebbleStorage) Scan(prefix []byte) [][]byte {
	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil // no upper-bound
	}

	var results [][]byte

	iter, err := p.handle.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: keyUpperBound(prefix),
	})
	if err != nil {
		log.Error("failed to create iterator", "error", err)
	}
	for iter.First(); iter.Valid(); iter.Next() {
		results = append(results, iter.Value())
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return results
}

func NewPebbleStorage(opts ...PebbleStorageOptsFunc) *PebbleStorage {
	o := defaultPebbleOptions()
	for _, optFn := range opts {
		optFn(o)
	}

	handle, err := pebble.Open(o.path, &pebble.Options{})
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return &PebbleStorage{
		handle: handle,
	}
}
