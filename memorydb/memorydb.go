package memorydb

import (
	"errors"
	"go-idb/iface"
	"sync"
)

type Database struct {
	db   map[string]interface{}
	lock sync.RWMutex
}

type batch struct {
	db     *Database
	writes []keyVal
	size   int
}

type keyVal struct {
	key    []string
	value  []byte
	delete bool
}

var (
	// errMemorydbClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errMemorydbClosed = errors.New("database closed")

	// errMemorydbNotFound is returned if a key is requested that is not found in
	// the provided memory database.
	errMemorydbNotFound = errors.New("not found")
	errKeyValMismatched = errors.New("len(keys) != lens(values)")
)

func New(db map[string]interface{}) *Database {
	return &Database{db: db}
}

func (db *Database) Has(key string) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, errMemorydbClosed
	}
	_, ok := db.db[key]
	return ok, nil
}

func (db *Database) Get(key string) (interface{}, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, errMemorydbClosed
	}
	if entry, ok := db.db[key]; ok {
		return entry, nil
	}
	return nil, errMemorydbNotFound
}

func (db *Database) Set(key string, val interface{}) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	db.db[key] = val
	return nil
}

func (db *Database) Delete(key string) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	delete(db.db, key)
	return nil
}

func (db *Database) Count() (int, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return 0, errMemorydbClosed
	}
	return len(db.db), nil
}

func (db *Database) Batch() iface.IBatch {
	return &batch{
		db: db,
	}
}

func (b *batch) Get(keys []string) ([]interface{}, error) {
	b.db.lock.RLock()
	defer b.db.lock.RUnlock()

	if b.db.db == nil {
		return nil, errMemorydbClosed
	}

	var entries []interface{}
	for _, key := range keys {
		entry, ok := b.db.db[key]
		if !ok {
			return entries, errMemorydbNotFound
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (b *batch) Set(keys []string, values []interface{}) error {
	b.db.lock.RLock()
	defer b.db.lock.RUnlock()

	if len(keys) != len(values) {
		return errKeyValMismatched
	}

	if b.db.db == nil {
		return errMemorydbClosed
	}

	for k, val := range values {
		b.db.db[keys[k]] = val
	}
	return nil
}

func (b *batch) Delete(keys []string) error {
	b.db.lock.RLock()
	defer b.db.lock.RUnlock()

	if b.db.db == nil {
		return errMemorydbClosed
	}

	for _, key := range keys {
		delete(b.db.db, key)
	}
	return nil
}
