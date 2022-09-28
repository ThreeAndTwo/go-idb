package monitordb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/go-idb/types"
	"github.com/deng00/go-base/cache/redis"
	"github.com/deng00/go-base/db/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"testing"
	"time"
)

var memoryDB map[string]interface{}
var redisDB *redis.Redis
var mysqlDB *mysql.MySQL
var influxDB influxdb2.Client
var taosDB *sql.DB

type influxConfig struct {
	host   string
	token  string
	org    string
	bucket string
}

var taosUri = "test:123456@tcp(localhost:6030)/"

var influxCnf = influxConfig{
	host:   "http://127.0.0.1:8086",
	token:  "nABt4rcDA1YEVFEdP3J9fxEGxNUlAcL4YoyQZHIXRJAbKBmUJGRwpZU_3agmr0ThhHyXwK0CawD_gKcaa4tLoA==",
	org:    "threeandtwo_org",
	bucket: "dev",
}

func init() {
	//memoryDB = make(map[string]interface{})
	//
	//redisConfig := &redis.Config{
	//	Addr:     "127.0.0.1:6379",
	//	Pass:     "",
	//	DB:       0,
	//	PoolSize: 100,
	//}
	//
	//_redisClient, err := redis.New(redisConfig)
	//if err != nil {
	//	panic("new redis client error:" + err.Error())
	//}
	//redisDB = _redisClient

	//influxDB = influxdb2.NewClient(
	//	influxCnf.host,
	//	influxCnf.token,
	//)
	//
	taos, err := sql.Open("taosSql", taosUri)
	if err != nil {
		panic("new TDEngine client error:" + err.Error())
	}
	taosDB = taos
}

func TestGetTS(t *testing.T) {
	tests := []struct {
		name   string
		dbTy   types.DBTy
		client interface{}
		org    string
		bucket string
	}{
		//{
		//	name:   "normal influxdb",
		//	dbTy:   types.InfluxDBTy,
		//	client: influxDB,
		//	org:    influxCnf.org,
		//	bucket: influxCnf.bucket,
		//},
		{
			name:   "normal TDEngine",
			dbTy:   types.TDEngineTy,
			client: taosDB,
			org:    "",
			bucket: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := GetTS(tt.dbTy, tt.client, tt.org, tt.bucket)
			if err != nil {
				t.Errorf("new no sql db error: %s", err.Error())
				return
			}

			var query string
			var val interface{}
			if tt.dbTy == types.InfluxDBTy {
				query = `from(bucket:"` + tt.bucket + `")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`
				val = influxdb2.NewPoint("stat",
					map[string]string{"unit": "temperature"},
					map[string]interface{}{"avg": 24.5, "max": 45.0},
					time.Now())
			} else {
				query = `select * from d0 limit 10`
				val = "insert into d0 values(NOW, 9.96000, 116, 0.32778)"
			}

			fmt.Println(val)
			err = ts.Insert(val)
			if err != nil {
				t.Errorf("insert data to time-series error: %s", err.Error())
				return
			}

			scanData, err := ts.Query(query)
			if err != nil {
				t.Errorf("query data to time-series error: %s", err.Error())
				return
			}

			marshal, _ := json.Marshal(scanData)
			t.Logf(string(marshal))
		})
	}
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
