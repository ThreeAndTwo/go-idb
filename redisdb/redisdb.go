package redisdb

import (
	"context"
	"errors"
	"github.com/ThreeAndTwo/go-idb/iface"
	"github.com/deng00/go-base/cache/redis"
	"time"
)

type Database struct {
	client  *redis.Redis
	ctx     context.Context
	timeout time.Duration
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
	// errRedisClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errRedisClosed      = errors.New("redis database closed")
	errUnsupported      = errors.New("unsupported function")
	errKeyValMismatched = errors.New("len(keys) != lens(values)")
)

func New(client *redis.Redis, timeout int) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	return &Database{client: client, ctx: ctx, timeout: time.Duration(timeout) * time.Second}
}

func (db *Database) Has(key string) (bool, error) {
	if db.client == nil {
		return false, errRedisClosed
	}
	return db.client.Exist(key)
}

func (db *Database) Get(key string) (interface{}, error) {
	if db.client == nil {
		return nil, errRedisClosed
	}
	return db.client.Get(key)
}

func (db *Database) Set(key string, val interface{}) error {
	if db.client == nil {
		return errRedisClosed
	}
	return db.client.Set(key, val, db.timeout)
}

func (db *Database) Delete(key string) error {
	if db.client == nil {
		return errRedisClosed
	}
	_, err := db.client.Del(key)
	return err
}

func (db *Database) Count() (int, error) {
	if db.client == nil {
		return 0, errRedisClosed
	}
	return 0, errUnsupported
}

func (db *Database) Batch() iface.IBatch {
	return &batch{
		db: db,
	}
}

func (b *batch) Get(keys []string) ([]interface{}, error) {
	if b.db.client == nil {
		return nil, errRedisClosed
	}

	var entries []interface{}
	for _, key := range keys {
		entry, err := b.db.client.Get(key)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (b *batch) Set(keys []string, values []interface{}) error {
	if b.db.client == nil {
		return errRedisClosed
	}

	if len(keys) != len(values) {
		return errKeyValMismatched
	}

	for k, val := range values {
		if err := b.db.client.Set(keys[k], val, b.db.timeout); err != nil {
			return err
		}
	}
	return nil
}

func (b *batch) Delete(keys []string) error {
	if b.db.client == nil {
		return errRedisClosed
	}

	for _, key := range keys {
		_, err := b.db.client.Del(key)
		if err != nil {
			return err
		}
	}
	return nil
}
