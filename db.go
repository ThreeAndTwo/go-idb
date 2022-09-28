package monitordb

import (
	"database/sql"
	"errors"
	"github.com/ThreeAndTwo/go-idb/gocache"
	"github.com/ThreeAndTwo/go-idb/iface"
	"github.com/ThreeAndTwo/go-idb/influxdb"
	"github.com/ThreeAndTwo/go-idb/memorydb"
	"github.com/ThreeAndTwo/go-idb/mysql"
	"github.com/ThreeAndTwo/go-idb/redisdb"
	"github.com/ThreeAndTwo/go-idb/tdengine"
	"github.com/ThreeAndTwo/go-idb/types"
	_mysql "github.com/deng00/go-base/db/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/deng00/go-base/cache/redis"
)

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

func GetTS(dbTy types.DBTy, client interface{}, org, bucket string) (iface.ITSDB, error) {
	switch dbTy {
	case types.InfluxDBTy:
		return influxdb.New(client.(influxdb2.Client), org, bucket)
	case types.TDEngineTy:
		return tdengine.New(client.(*sql.DB))
	default:
		return nil, errDBUnSupported
	}
}
