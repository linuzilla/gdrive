package intf

import "github.com/linuzilla/gdrive/utils/dbconn"

type SqlDatabaseInterface interface {
	Connection(callback func(xdb *dbconn.DB) error) error
}
