package tdengine

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"reflect"
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

func (db *Database) Query(raw string, res interface{}) ([]interface{}, error) {
	if err, ok := db.check(); ok {
		return nil, err
	}

	fmt.Println("raw:", raw)
	result, err := db.db.Query(raw)
	if err != nil {
		return nil, err
	}

	marshal, _ := json.Marshal(result)

	fmt.Println("result:,", string(marshal))

	var data []interface{}
	for result.Next() {
		if result.Scan(res) != nil {
			return nil, err
		}
		data = append(data, res)
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
