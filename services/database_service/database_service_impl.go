package database_service

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/utils/dbconn"
	"log"
)

type databaseServiceImpl struct {
	Conn dbconn.IDBConnection `inject:"*"`
}

var _ intf.SqlDatabaseInterface = (*databaseServiceImpl)(nil)

func (db *databaseServiceImpl) PostSummerConstruct() {
	if _, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		xdb.Gorm.AutoMigrate(&models.SyncFileInfo{})
		xdb.Gorm.AutoMigrate(&models.GoogleDriveConfig{})
		return nil, nil
	}); err != nil {
		log.Fatalf("%v", err)
	}
}

func New() intf.SqlDatabaseInterface {
	return &databaseServiceImpl{}
}

func (db *databaseServiceImpl) Connection(callback func(xdb *dbconn.DB) error) error {
	_, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		return nil, callback(xdb)
	})
	return err
}

func (db *databaseServiceImpl) ClearLocalCache() {
	if _, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		xdb.Gorm.Delete(&models.SyncFileInfo{})
		return nil, nil
	}); err != nil {
		log.Fatalf("%v", err)
	}
}

func (db *databaseServiceImpl) SaveLocalFile(file *models.SyncFileInfo) {
	if _, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		xdb.Gorm.Save(file)
		return nil, nil
	}); err != nil {
		log.Fatalf("%v", err)
	}
}

func (db *databaseServiceImpl) Find(data interface{}) error {
	_, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		return nil, xdb.Gorm.Find(data).Error
	})
	return err
}

func (db *databaseServiceImpl) WhereFirst(where interface{}, out interface{}) error {
	_, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		return nil, xdb.Gorm.Where(where).First(out).Error
	})
	return err
}

func (db *databaseServiceImpl) Where(where interface{}, out interface{}) error {
	_, err := db.Conn.GetConnection(func(xdb *dbconn.DB) (interface{}, error) {
		return nil, xdb.Gorm.Where(where).Find(out).Error
	})
	return err
}
