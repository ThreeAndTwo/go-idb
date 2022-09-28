package influxdb

import (
	"context"
	"errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"reflect"
)

var (
	InfluxWApi api.WriteAPIBlocking
)

type Database struct {
	db   influxdb2.Client
	conf config
}

type config struct {
	org    string
	bucket string
}

var (
	// errClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errClosed      = errors.New("database closed")
	errDbConfig    = errors.New("config invalidate for influxdb")
	errPointIsNull = errors.New("point is nil")
	errInsertType  = errors.New("value type invalidate")
)

func New(
	db influxdb2.Client,
	org string,
	bucket string,
) (*Database, error) {
	InfluxWApi = db.WriteAPIBlocking(org, bucket)
	_db := &Database{
		db: db,
		conf: config{
			org:    org,
			bucket: bucket,
		},
	}

	err, _ := _db.check()
	return _db, err
}

func (db *Database) check() (error, bool) {
	if db.db == nil {
		return errClosed, true
	} else if db.conf.bucket == "" || db.conf.org == "" {
		return errDbConfig, true
	}
	return nil, false
}

func (db *Database) Query(raw string, res ...interface{}) ([]interface{}, error) {
	if err, ok := db.check(); ok {
		return nil, err
	}

	result, err := db.db.QueryAPI(db.conf.org).Query(context.Background(), raw)
	if err != nil {
		return nil, err
	}

	var data []interface{}
	for result.Next() {
		data = append(data, result.Record().Value())
	}
	return data, nil
}

func (db *Database) Insert(val interface{}) error {
	if err, ok := db.check(); ok {
		return err
	}

	if reflect.TypeOf(val).String() != "*write.Point" {
		return errInsertType
	}

	p := val.(*write.Point)
	if p == nil {
		return errPointIsNull
	}
	return InfluxWApi.WritePoint(context.Background(), p)
}
