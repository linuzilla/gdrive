package environment_service

import (
	"fmt"
	"github.com/linuzilla/gdrive/config"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/go-logger"
	"os"
)

type EnvironmentService interface {
	RootFolder() string
	GetWorkingDirectory() string
	GetEditor() string
	SetEncoding(enabled bool)
	Encoding() bool
	UpdateTopLevelDatabase(conf *models.GoogleDriveConfig, dir string) error
	FetchTopLevel() ([]models.GoogleDriveConfig, error)
	SetCoder(codec *models.Codec)
	GetCodec() *models.Codec
	ListSyncInfos()
	LoadCache(index int) *models.GoogleDriveConfig
	LoadFromTopLevelDatabase(id string) (*models.GoogleDriveConfig, error)
	LoadFromTopLevelDatabaseByFolderId(folderId string) (*models.GoogleDriveConfig, error)
}

type environmentServiceImpl struct {
	Conf             *config.Config       `inject:"*"`
	DatabaseBackend  intf.DatabaseBackend `inject:"*"`
	workingDirectory string
	encoding         bool
	codec            models.Codec
	cache            map[int]models.GoogleDriveConfig
}

func New() EnvironmentService {
	return &environmentServiceImpl{}
}

func (svc *environmentServiceImpl) PostSummerConstruct() {
	svc.codec = svc.Conf.Codec

	if svc.Conf.Application.WorkingDir == `` {
		if dir, err := os.Getwd(); err != nil {
			svc.workingDirectory = "/tmp"
		} else {
			svc.workingDirectory = dir
		}
	} else {
		svc.workingDirectory = svc.Conf.Application.WorkingDir
	}
}

func (svc *environmentServiceImpl) LoadCache(index int) *models.GoogleDriveConfig {
	if svc.cache != nil {
		if entry, found := svc.cache[index]; found {
			return &entry
		}
	}
	return nil
}

func (svc *environmentServiceImpl) ListSyncInfos() {
	driveConfig := map[int]models.GoogleDriveConfig{}

	if topLevel, err := svc.FetchTopLevel(); err != nil {
		logger.Error(err)
	} else {
		for i, entry := range topLevel {
			driveConfig[i+1] = entry
			fmt.Printf("%4d. [ %s ] %s\n", i+1, entry.FolderId, entry.Id)
		}
	}
	svc.cache = driveConfig
}

func (svc *environmentServiceImpl) GetEditor() string {
	if svc.Conf.Application.Editor != `` {
		return svc.Conf.Application.Editor
	} else {
		return constants.DefaultEditor
	}
}
func (svc *environmentServiceImpl) RootFolder() string {
	return svc.Conf.GoogleDrive.FolderId
}

func (svc *environmentServiceImpl) GetCodec() *models.Codec {
	return &svc.codec
}

func (svc *environmentServiceImpl) GetWorkingDirectory() string {
	return svc.workingDirectory
}

func (svc *environmentServiceImpl) SetEncoding(enable bool) {
	svc.encoding = enable
}

func (svc *environmentServiceImpl) Encoding() bool {
	return svc.encoding
}

func (svc *environmentServiceImpl) FetchTopLevel() ([]models.GoogleDriveConfig, error) {
	var driveConfig []models.GoogleDriveConfig

	err := svc.DatabaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		if data, err := connection.FindAllConfig(); err != nil {
			return err
		} else {
			driveConfig = data
			return nil
		}
	})

	return driveConfig, err
}

func (svc *environmentServiceImpl) SetCoder(codec *models.Codec) {
	svc.codec = *codec
}

func (svc *environmentServiceImpl) LoadFromTopLevelDatabaseByFolderId(folderId string) (*models.GoogleDriveConfig, error) {
	var stored models.GoogleDriveConfig

	if err := svc.DatabaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		return connection.ReadConfigByFolderId(folderId, &stored)
	}); err != nil {
		return nil, err
	} else {
		return &stored, nil
	}
}

func (svc *environmentServiceImpl) LoadFromTopLevelDatabase(id string) (*models.GoogleDriveConfig, error) {
	var stored models.GoogleDriveConfig

	if err := svc.DatabaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		return connection.ReadConfig(id, &stored)
	}); err != nil {
		return nil, err
	} else {
		return &stored, nil
	}
}

func (svc *environmentServiceImpl) UpdateTopLevelDatabase(conf *models.GoogleDriveConfig, dir string) error {
	return svc.DatabaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		var stored models.GoogleDriveConfig

		if err := connection.ReadConfig(dir, &stored); err != nil {
			stored = *conf
			stored.Id = dir
			if err := connection.SaveConfig(&stored); err != nil {
				logger.Error(err)
				return err
			} else {
				logger.Notice("config [ %s ] saved", dir)
			}
		} else {
			current := *conf
			current.Id = dir
			current.CreatedAt = stored.CreatedAt
			current.UpdatedAt = stored.UpdatedAt

			if stored != current {
				if err := connection.SaveConfig(&current); err != nil {
					logger.Error(err)
					return err
				} else {
					logger.Notice("config [ %s ] updated", dir)
				}
			}
		}
		return nil
	})
}
