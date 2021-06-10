package dbsqlite3

import (
	"github.com/linuzilla/gdrive/utils/dbconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqlite3Connection struct {
	dbFile  string
	logmode bool
}

func (db *sqlite3Connection) testConnection() {
}

func (svc *sqlite3Connection) GetConnection(callback func(xdb *dbconn.DB) (interface{}, error)) (interface{}, error) {
	if conn, err := gorm.Open(sqlite.Open(svc.dbFile)); err != nil {
		return nil, err
	} else {
		//defer conn.Close()
		//conn.LogMode(svc.logmode)
		db, err := conn.DB()
		if err != nil {
			return nil, err
		}
		return callback(&dbconn.DB{Gorm: conn, SqlDB: db})
	}
}

func (svc *sqlite3Connection) GetExtraProperty(property string) interface{} {
	return nil
}

func (svc *sqlite3Connection) GetDialect() dbconn.IDbDialect {
	return dialect
}

func New(dbinfo *dbconn.DBConnectionInfo) dbconn.IDBConnection {
	return &sqlite3Connection{
		dbFile:  dbinfo.File,
		logmode: dbinfo.LogMode,
	}
}
