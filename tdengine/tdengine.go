package tdengine

import (
	"database/sql"
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

func (db *Database) Query(raw string, res ...interface{}) ([]interface{}, error) {
	if err, ok := db.check(); ok {
		return nil, err
	}

	var result *sql.Rows
	var err error

	result, err = db.db.Query(raw)
	//if res == nil || len(res) == 0 {
	//	fmt.Println("111")
	//	result, err = db.db.Query(raw)
	//} else {
	//	fmt.Println("222")
	//	result, err = db.db.Query(raw, res)
	//}

	if err != nil {
		return nil, err
	}

	//marshal, _ := json.Marshal(result)
	//
	//fmt.Println("result:,", string(marshal))
	fmt.Println("aaaa", result.Next())

	var data []interface{}
	for result.Next() {
		var _res interface{}
		if result.Scan(_res) != nil {
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
