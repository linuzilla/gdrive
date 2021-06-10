package intf

import "github.com/linuzilla/gdrive/models"

type DatabaseBackend interface {
	Initialize(databaseFile string, debug bool, callback func(expose interface{})) error
	ConnectionEstablish(callback func(connection DatabaseBackendConnection) error) error
}

type DatabaseBackendConnection interface {
	CreateDatabase() error
	FindFirstById(id string, data *models.SyncFileInfo) error
	Persist(data *models.SyncFileInfo) error
	SaveOrUpdate(data *models.SyncFileInfo) error
	Delete(data *models.SyncFileInfo) error

	ReadConfig(id string, data *models.GoogleDriveConfig) error
	ReadConfigByFolderId(folderId string, data *models.GoogleDriveConfig) error
	SaveConfig(data *models.GoogleDriveConfig) error
	FindAllConfig() ([]models.GoogleDriveConfig, error)

	Find(out interface{}, where ...interface{}) error
}
