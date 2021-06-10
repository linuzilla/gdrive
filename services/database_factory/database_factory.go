package database_factory

import "github.com/linuzilla/gdrive/intf"

type DatabaseFactory interface {
	NewDatabase(databaseFile string) intf.DatabaseBackend
}

type databaseFactoryImpl struct {
}

func New() DatabaseFactory {
	return &databaseFactoryImpl{}
}
