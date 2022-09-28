package iface

import "github.com/ThreeAndTwo/go-idb/types"

type (
	INoSQL interface {
		Has(key string) (bool, error)

		Get(key string) (interface{}, error)
		Set(key string, val interface{}) error
		Delete(key string) error

		Count() (int, error)
		Batch() IBatch
	}

	IBatch interface {
		Get(keys []string) ([]interface{}, error)
		Set(keys []string, values []interface{}) error
		Delete(keys []string) error
	}
)

type (
	ISQL interface {
		Find(fields *types.FindSqlField, val interface{}) error
		InsertOne(tableName string, val interface{}) error
		InsertMany(tableName string, values []interface{}) error
		Update(val *types.UpdateField) error
		Delete(val *types.DeleteField) error

		Raw(raw *types.RawField, res interface{}) error
	}
)

type (
	ITSDB interface {
		Query(raw string, res ...interface{}) ([]interface{}, error)
		Insert(val interface{}) error
	}
)
