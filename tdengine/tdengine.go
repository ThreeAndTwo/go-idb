package tdengine

import (
	"database/sql"
	"errors"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"reflect"
	"time"
)

type Database struct {
	db *sql.DB
}

var (
	// errClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errClosed     = errors.New("database closed")
	errInsertType = errors.New("value type invalidate")
)

func New(db *sql.DB) (*Database, error) {
	_db := &Database{db: db}

	err, _ := _db.check()
	return _db, err
}

func (db *Database) check() (error, bool) {
	if db.db == nil {
		return errClosed, true
	}
	return nil, false
}

func (db *Database) Query(query string, args ...interface{}) ([]interface{}, error) {
	if err, ok := db.check(); ok {
		return nil, err
	}

	var result *sql.Rows
	var err error

	result, err = db.db.Query(query)
	if err != nil {
		return nil, err
	}

	var data []interface{}

	for result.Next() {
		var _res struct {
			Ts      time.Time
			Current float64
			Voltage int
			Phase   float64
		}

		if result.Scan(&_res.Ts, &_res.Current, &_res.Voltage, &_res.Phase) != nil {
			return nil, err
		}
		data = append(data, _res)
	}
	return data, nil
}

func (db *Database) Insert(val interface{}) error {
	if err, ok := db.check(); ok {
		return err
	}

	if reflect.TypeOf(val).Kind() != reflect.String {
		return errInsertType
	}

	_, err := db.db.Exec(val.(string))
	return err
}
