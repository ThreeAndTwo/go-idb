package monitordb

import (
	"github.com/deng00/go-base/cache/redis"
	"github.com/deng00/go-base/db/mysql"
	"go-idb/types"
	"testing"
)

var memoryDB map[string]interface{}
var redisDB *redis.Redis
var mysqlDB *mysql.MySQL

func init() {
	memoryDB = make(map[string]interface{})

	redisConfig := &redis.Config{
		Addr:     "127.0.0.1:6379",
		Pass:     "",
		DB:       0,
		PoolSize: 100,
	}

	_redisClient, err := redis.New(redisConfig)
	if err != nil {
		panic("new redis client error:" + err.Error())
	}
	redisDB = _redisClient

	//mysqlConfig := new(mysql.Config)
	//
	//mysqlDB, err = mysql.New(mysqlConfig)
	//if err != nil {
	//	panic("new mysql client error:" + err.Error())
	//}
}

func TestGetNoSqlDB(t *testing.T) {
	key0 := "1:99"
	val0 := "aaa"

	key1 := "56:1000"
	val1 := "bsc"

	keys, values := []string{key0, key1}, []interface{}{val0, val1}

	var tests = []struct {
		name   string
		dbTy   types.DBTy
		client interface{}
	}{
		{
			name:   "test memoryDB",
			dbTy:   types.MemDBTy,
			client: memoryDB,
		},
		{
			name:   "test go cache",
			dbTy:   types.GoCacheTy,
			client: nil,
		},
		{
			name:   "test redis db",
			dbTy:   types.RedisDBTy,
			client: redisDB,
		},
		{
			name:   "test memoryDB for nil",
			dbTy:   types.MemDBTy,
			client: nil,
		},
		{
			name:   "test redisDB for nil",
			dbTy:   types.RedisDBTy,
			client: nil,
		},
		{
			name:   "null db type",
			dbTy:   "",
			client: types.MemDBTy,
		},
		{
			name:   "both of db type, client",
			dbTy:   "",
			client: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := GetNoSqlDB(tt.dbTy, tt.client, 0)
			if err != nil {
				t.Errorf("new no sql db error: %s", err.Error())
				return
			}

			err = db.Set(key0, val0)
			if err != nil {
				t.Errorf("set val for %s key error: %s", key0, err.Error())
				return
			}

			_val, err := db.Get(key0)
			if err != nil {
				t.Errorf("get val for %s key error: %s", key0, err.Error())
				return
			}
			t.Logf("val: %s", _val)

			has, err := db.Has(key0)
			if err != nil {
				t.Errorf(" %s key error: %s", key0, err.Error())
				return
			}
			t.Logf("is has key? %t", has)

			err = db.Batch().Set(keys, values)
			if err != nil {
				t.Errorf("batch set error: %s", err.Error())
				return
			}

			_values, err := db.Batch().Get(keys)
			if err != nil {
				t.Errorf("batch get error: %s", err.Error())
				return
			}
			t.Logf("batch get value: %v", _values)

			err = db.Batch().Delete(keys)
			if err != nil {
				t.Errorf("delete keys error: %s", err.Error())
			}
			t.Logf("db execution success!")
		})
	}
}

func TestGetSQL(t *testing.T) {
	var tests = []struct {
		name   string
		dbTy   types.DBTy
		client interface{}
	}{
		{
			name:   "test mysql db",
			dbTy:   types.MysqlDBTy,
			client: mysqlDB,
		},
		{
			name:   "db type is null",
			dbTy:   "",
			client: mysqlDB,
		},
		{
			name:   "client is null",
			dbTy:   types.MysqlDBTy,
			client: nil,
		},
		{
			name:   "both of null for db type, client",
			dbTy:   "",
			client: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//sql, err := GetSQL(tt.dbTy, tt.client)
			//if err != nil {
			//	t.Errorf("new no sql db error: %s", err.Error())
			//	return
			//}
			//
			//sql.Find()
		})
	}
}
