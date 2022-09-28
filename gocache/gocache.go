package gocache

import (
	"errors"
	"github.com/ThreeAndTwo/go-idb/iface"
	"github.com/patrickmn/go-cache"
	"time"
)

type Database struct {
	db              *cache.Cache
	expire, cleanup time.Duration
}

type batch struct {
	db   *Database
	size int
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

func New(timeout, cleanup int) *Database {
	_expire := time.Duration(timeout) * time.Second
	_cleanup := time.Duration(cleanup) * time.Second
	db := cache.New(_expire, _cleanup)
	return &Database{
		db:      db,
		expire:  _expire,
		cleanup: _cleanup,
	}
}

func (db *Database) Has(key string) (bool, error) {
	if db.db == nil {
		return false, errMemorydbClosed
	}
	_, ok := db.db.Get(key)
	return ok, nil
}

func (db *Database) Get(key string) (interface{}, error) {
	if db.db == nil {
		return false, errMemorydbClosed
	}
	if entry, ok := db.db.Get(key); ok {
		return entry, nil
	}
	return nil, errMemorydbNotFound
}

func (db *Database) Set(key string, val interface{}) error {
	if db.db == nil {
		return errMemorydbClosed
	}
	db.db.Set(key, val, db.expire)
	return nil
}

func (db *Database) Delete(key string) error {
	if db.db == nil {
		return errMemorydbClosed
	}
	db.db.Delete(key)
	return nil
}

func (db *Database) Count() (int, error) {
	if db.db == nil {
		return 0, errMemorydbClosed
	}
	return db.db.ItemCount(), nil
}

func (db *Database) Batch() iface.IBatch {
	return &batch{
		db: db,
	}
}

func (b *batch) Get(keys []string) ([]interface{}, error) {
	if b.db.db == nil {
		return nil, errMemorydbClosed
	}

	var entries []interface{}
	for _, key := range keys {
		entry, ok := b.db.db.Get(key)
		if !ok {
			return entries, errMemorydbNotFound
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (b *batch) Set(keys []string, values []interface{}) error {

	if len(keys) != len(values) {
		return errKeyValMismatched
	}

	if b.db.db == nil {
		return errMemorydbClosed
	}

	for k, val := range values {
		b.db.db.Set(keys[k], val, b.db.expire)
	}
	return nil
}

func (b *batch) Delete(keys []string) error {
	if b.db.db == nil {
		return errMemorydbClosed
	}

	for _, key := range keys {
		err := b.db.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
