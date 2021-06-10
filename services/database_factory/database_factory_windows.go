// +build windows

package database_factory

import (
	"github.com/linuzilla/gdrive/backends/db_sqlite3"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/database_service"
	"github.com/linuzilla/summer"
)

func (databaseFactoryImpl) NewDatabase(databaseFile string) intf.DatabaseBackend {
	applicationContext := summer.New()

	databaseBackend := db_sqlite3.New()
	applicationContext.Add(database_service.New())
	applicationContext.Add(databaseBackend)

	databaseBackend.Initialize(databaseFile, false, func(expose interface{}) {
		applicationContext.Add(expose)
	})

	done := applicationContext.Autowiring(func(err error) {})

	if err := <-done; err != nil {
		panic(err.Error())
	}

	return databaseBackend
}
