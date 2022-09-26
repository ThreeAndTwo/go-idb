package monitordb

import (
	"errors"
	_mysql "github.com/deng00/go-base/db/mysql"
	"go-idb/gocache"
	"go-idb/iface"
	"go-idb/memorydb"
	"go-idb/mysql"
	"go-idb/redisdb"
	"go-idb/types"

	"github.com/deng00/go-base/cache/redis"
)

type MonitorDB struct {
}

var (
	errDBUnSupported = errors.New("unsupported db type")
	errClientIsNull  = errors.New("client for db is nil")
)

func GetNoSqlDB(dbTy types.DBTy, client interface{}, timeout int) (iface.INoSQL, error) {
	//check
	if dbTy != types.GoCacheTy && client == nil {
		return nil, errClientIsNull
	}

	switch dbTy {
	case types.MemDBTy:
		_client := memorydb.New(client.(map[string]interface{}))
		return _client, nil
	case types.GoCacheTy:
		_client := gocache.New(timeout, timeout)
		return _client, nil
	case types.RedisDBTy:
		_client := redisdb.New(client.(*redis.Redis), timeout)
		return _client, nil
	default:
		return nil, errDBUnSupported
	}
}

func GetSQL(dbTy types.DBTy, client interface{}) (iface.ISQL, error) {
	switch dbTy {
	case types.MysqlDBTy:
		_client := mysql.New(client.(*_mysql.MySQL))
		return _client, nil
	default:
		_client := mysql.New(client.(*_mysql.MySQL))
		return _client, nil
	}
}
