package db_bolt_impl

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/go-logger"
	"log"
	"time"
)

type BoltDbBackend struct {
	databaseFile   string
	debug          bool
	databaseBucket []byte
}

type boltConnection struct {
	db      *bolt.DB
	backend *BoltDbBackend
}

var _ intf.DatabaseBackendConnection = (*boltConnection)(nil)
var _ intf.DatabaseBackend = (*BoltDbBackend)(nil)

var BoltBackend BoltDbBackend

///////////////////////////////////////////////////

func (backend *BoltDbBackend) Initialize(databaseFile string, debug bool, callback func(interface{})) error {
	fmt.Println("initialize database: " + databaseFile)

	bucketName := "localFile"

	backend.databaseFile = databaseFile + ".db"
	backend.debug = debug
	backend.databaseBucket = []byte(bucketName)

	if err := backend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		if err := connection.CreateDatabase(); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (backend *BoltDbBackend) ConnectionEstablish(callback func(connection intf.DatabaseBackendConnection) error) error {
	if db, err := bolt.Open(backend.databaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		return err
	} else {
		defer db.Close()
		return callback(&boltConnection{db: db, backend: backend})
	}
}

///////////////////////////////////////////////////

func (conn *boltConnection) CreateDatabase() error {
	return conn.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(conn.backend.databaseBucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
}

func (conn *boltConnection) FindFirstById(id string, data *models.SyncFileInfo) error {
	return conn.db.View(func(tx *bolt.Tx) error {
		return json.Unmarshal(tx.Bucket(conn.backend.databaseBucket).Get([]byte(id)), data)
	})
}

func (conn *boltConnection) Persist(data *models.SyncFileInfo) error {
	return conn.SaveOrUpdate(data)
}

func (conn *boltConnection) SaveOrUpdate(data *models.SyncFileInfo) error {
	return conn.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(conn.backend.databaseBucket)
		if jsonBlob, err := json.Marshal(data); err != nil {
			return err
		} else {
			return bucket.Put([]byte(data.Id), jsonBlob)
		}
	})
}

func (conn *boltConnection) Delete(data *models.SyncFileInfo) error {
	return conn.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(conn.backend.databaseBucket).Delete([]byte(data.Id))
	})
}

func (conn *boltConnection) Find(out interface{}, where ...interface{}) error {
	logger.Fatal("Not implemented !")
	return nil
}

func (conn *boltConnection) ReadConfigByFolderId(folderId string, data *models.GoogleDriveConfig) error {
	logger.Fatal("Not implemented !")
	return nil
}

func (conn *boltConnection) ReadConfig(id string, data *models.GoogleDriveConfig) error {
	logger.Fatal("Not implemented !")
	return nil
}

func (conn *boltConnection) FindAllConfig() ([]models.GoogleDriveConfig, error) {
	logger.Fatal("Not implemented !")
	return nil, nil
}

func (conn *boltConnection) SaveConfig(data *models.GoogleDriveConfig) error {
	logger.Fatal("Not implemented !")
	return nil
}
