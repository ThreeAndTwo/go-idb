package mysql

import (
	"errors"
	"github.com/deng00/go-base/db/mysql"
	"go-idb/types"
)

type Database struct {
	sql *mysql.MySQL
}

var (
	// errMemorydbClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errMysqldbClosed = errors.New("database closed")
)

func New(db *mysql.MySQL) *Database {
	return &Database{sql: db}
}

// Find
// TODO: should be support limit, offset, group, having
func (db *Database) Find(fields *types.FindSqlField, val interface{}) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}

	if fields.Select.Args == nil {
		return db.sql.Client.Table(fields.TBName).Select(fields.Select.Query).
			Where(fields.Where.Query, fields.Where.Args).Scan(val).Error
	}

	if fields.Limit != nil {
		return db.sql.Client.Table(fields.TBName).Select(fields.Select.Query, fields.Select.Args).
			Where(fields.Where.Query, fields.Where.Args).Scan(val).Limit(fields.Limit).Error
	}

	return db.sql.Client.Table(fields.TBName).Select(fields.Select.Query, fields.Select.Args).
		Where(fields.Where.Query, fields.Where.Args).Scan(val).Error
}

func (db *Database) InsertOne(tbName string, val interface{}) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}
	return db.sql.Client.Table(tbName).Create(val).Error
}

func (db *Database) InsertMany(tableName string, values []interface{}) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}
	for _, val := range values {
		if err := db.sql.Client.Table(tableName).Create(val).Error; err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Update(val *types.UpdateField) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}
	return db.sql.Client.Table(val.TBName).Where(val.Where.Query, val.Where.Args).Updates(val.Values).Error
}

func (db *Database) Delete(val *types.DeleteField) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}
	return db.sql.Client.Table(val.TBName).Where(val.Where.Query, val.Where.Args).Delete(val.Values).Error
}

func (db *Database) Raw(raw *types.RawField, res interface{}) error {
	if db.sql.Client == nil {
		return errMysqldbClosed
	}
	if raw.Values == nil {
		return db.sql.Client.Raw(raw.Sql).Scan(res).Error
	}
	return db.sql.Client.Raw(raw.Sql, raw.Values).Scan(res).Error
}
