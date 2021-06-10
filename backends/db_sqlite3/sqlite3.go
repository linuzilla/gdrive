package db_sqlite3

import (
	"github.com/linuzilla/gdrive/intf"
)

import (
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/utils/dbconn"
	"github.com/linuzilla/gdrive/utils/dbsqlite3"
)

type sqlite3Backend struct {
	SqlDb intf.SqlDatabaseInterface `inject:"*"`
}

type sqliteConnection struct {
	xdb *dbconn.DB
}

var _ intf.DatabaseBackendConnection = (*sqliteConnection)(nil)
var _ intf.DatabaseBackend = (*sqlite3Backend)(nil)

func (conn *sqliteConnection) CreateDatabase() error {
	return nil
}

func (conn *sqliteConnection) FindFirstById(id string, data *models.SyncFileInfo) error {
	return conn.xdb.Gorm.Where(&models.SyncFileInfo{Id: id}).First(data).Error
}

func (conn *sqliteConnection) Persist(data *models.SyncFileInfo) error {
	return conn.xdb.Gorm.Create(data).Error
}

func (conn *sqliteConnection) SaveOrUpdate(data *models.SyncFileInfo) error {
	return conn.xdb.Gorm.Save(data).Error
}

func (conn *sqliteConnection) Delete(data *models.SyncFileInfo) error {
	return conn.xdb.Gorm.Delete(data).Error
}

func (conn *sqliteConnection) ReadConfig(id string, data *models.GoogleDriveConfig) error {
	return conn.xdb.Gorm.Where(&models.GoogleDriveConfig{Id: id}).First(data).Error
}

func (conn *sqliteConnection) ReadConfigByFolderId(folderId string, data *models.GoogleDriveConfig) error {
	return conn.xdb.Gorm.Where(&models.GoogleDriveConfig{FolderId: folderId}).First(data).Error
}

func (conn *sqliteConnection) SaveConfig(data *models.GoogleDriveConfig) error {
	return conn.xdb.Gorm.Save(data).Error
}

func (conn *sqliteConnection) Find(out interface{}, where ...interface{}) error {
	return conn.xdb.Gorm.Find(out, where).Error
}

func (conn *sqliteConnection) FindAllConfig() ([]models.GoogleDriveConfig, error) {
	var data []models.GoogleDriveConfig

	rows, err := conn.xdb.Gorm.Model(&data).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		entry := &models.GoogleDriveConfig{}
		conn.xdb.Gorm.ScanRows(rows, entry)

		data = append(data, *entry)
	}
	return data, nil

	//return conn.xdb.Gorm.Where(where).Find(out).Error
	//return conn.xdb.Gorm.Find(out, where).Error
}

func (backend *sqlite3Backend) ConnectionEstablish(callback func(connection intf.DatabaseBackendConnection) error) error {
	return backend.SqlDb.Connection(func(xdb *dbconn.DB) error {
		return callback(&sqliteConnection{xdb: xdb})
	})
}

func (backend *sqlite3Backend) Initialize(databaseFile string, debug bool, callback func(interface{})) error {
	callback(dbsqlite3.New(&dbconn.DBConnectionInfo{
		File:    databaseFile + ".sqlite3",
		LogMode: debug,
	}))
	return nil
}

func New() intf.DatabaseBackend {
	return &sqlite3Backend{}
}
