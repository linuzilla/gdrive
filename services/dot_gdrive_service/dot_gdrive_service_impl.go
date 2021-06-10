package dot_gdrive_service

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/database_factory"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"os"
)

type dotGDriveServiceImpl struct {
	DatabaseFactory database_factory.DatabaseFactory       `inject:"*"`
	CryptoFactory   crypto_service.CryptoFactory           `inject:"*"`
	Api             google_drive.API                       `inject:"*"`
	EnvSvc          environment_service.EnvironmentService `inject:"*"`
	databaseBackend intf.DatabaseBackend
	currentEnv      *DotEnvironment
}

func New() DotGDriveService {
	return &dotGDriveServiceImpl{}
}

func (service *dotGDriveServiceImpl) PostSummerConstruct() {
	service.LoadEnvironment()
}

func (service *dotGDriveServiceImpl) updateEnvironment(conf *models.GoogleDriveConfig, dir string) {
	service.EnvSvc.UpdateTopLevelDatabase(conf, dir)
}

func (service *dotGDriveServiceImpl) loadData(conf *models.GoogleDriveConfig, env *DotEnvironment) {
	if conf.Password.Valid {
		if decodeString, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(conf.Password.String); err != nil {
			fmt.Println(err)
		} else {
			env.Password = string(decodeString)
		}
	}

	env.FolderId = conf.FolderId

	if conf.Encoder.Valid && conf.Decoder.Valid {
		confCodec := &models.Codec{
			Encoder: conf.Encoder.String,
			Decoder: conf.Decoder.String,
		}

		service.CryptoFactory.SetCoder(confCodec)
	}
}

func (service *dotGDriveServiceImpl) reload(googleDriveFolder string, dir string, env *DotEnvironment) {
	service.databaseBackend = service.DatabaseFactory.NewDatabase(googleDriveFolder + `/` + constants.DatabaseFile)

	_ = service.databaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		var conf models.GoogleDriveConfig

		if err := connection.ReadConfig(constants.ConfigDatabaseKey, &conf); err == nil {

			service.loadData(&conf, env)
			//fmt.Println(utils.Detail(&conf))

			//if conf.Password.Valid {
			//	if decodeString, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(conf.Password.String); err != nil {
			//		fmt.Println(err)
			//	} else {
			//		env.Password = string(decodeString)
			//	}
			//}
			//
			//env.FolderId = conf.FolderId
			//
			//if conf.Encoder.Valid && conf.Decoder.Valid {
			//	confCodec := &models.Codec{
			//		Encoder: conf.Encoder.String,
			//		Decoder: conf.Decoder.String,
			//	}
			//
			//	service.CryptoFactory.SetCoder(confCodec)
			//}
			//
			if !conf.Encoder.Valid || !conf.Decoder.Valid {
				cryptoSvc := service.CryptoFactory.GetInstance()

				conf.Encoder = sql.NullString{String: cryptoSvc.GetEncoder(), Valid: true}
				conf.Decoder = sql.NullString{String: cryptoSvc.GetDecoder(), Valid: true}

				logger.Notice("Add encoder/decoder to %s", dir)
				connection.SaveConfig(&conf)
			}

			fmt.Printf("Load config [ %s ], folderId: [ %s ], Encoding: %s\n",
				dir, env.FolderId, constants.OnOffMap[env.Password != ``])

			service.updateEnvironment(&conf, dir)
		}
		return nil
	})
}

func (service *dotGDriveServiceImpl) LoadTopLevelEnvByFolder(folderId string) *DotEnvironment {
	if conf, err := service.EnvSvc.LoadFromTopLevelDatabaseByFolderId(folderId); err != nil {
		fmt.Println(err)
		return nil
	} else {
		env := &DotEnvironment{}
		if currentDir, err := os.Getwd(); err == nil {
			env.Dir = currentDir
		}
		service.loadData(conf, env)
		return env
	}
}

func (service *dotGDriveServiceImpl) LoadTopLevelEnvironment(targetDir string) *DotEnvironment {
	if conf, err := service.EnvSvc.LoadFromTopLevelDatabase(targetDir); err != nil {
		fmt.Println(err)
		return nil
	} else {
		env := &DotEnvironment{Dir: targetDir}
		service.loadData(conf, env)
		return env
	}
}

func (service *dotGDriveServiceImpl) LoadEnvironment() *DotEnvironment {
	if dir, err := os.Getwd(); err == nil {
		if service.currentEnv != nil && service.currentEnv.Dir == dir {
			return service.currentEnv
		} else {
			googleDriveFolder := dir + `/` + constants.GoogleDriveFolder

			if fileInfo, err := os.Stat(googleDriveFolder); err == nil {
				if fileInfo.IsDir() {
					env := &DotEnvironment{Dir: dir}
					service.reload(googleDriveFolder, dir, env)
					return env
				}
			}
		}
	}
	return nil
}

func (service *dotGDriveServiceImpl) creationGDriveFolder(folderName string) (string, error) {
	googleDriveFolder := folderName + `/` + constants.GoogleDriveFolder

	if fileInfo, err := os.Stat(googleDriveFolder); os.IsNotExist(err) {
		if err := os.Mkdir(googleDriveFolder, 0700); err != nil {
			return "", err
		} else {
			fmt.Printf("Create directory: %s\n", googleDriveFolder)
		}
	} else if !fileInfo.IsDir() {
		return ``, fmt.Errorf("%s : should be a directory", googleDriveFolder)
	}

	return googleDriveFolder, nil
}

func (service *dotGDriveServiceImpl) PrepareBeforeSync(folderName string, folderId string,
	callback func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId string)) {

	googleDriveFolder, createErr := service.creationGDriveFolder(folderName)

	if createErr != nil {
		fmt.Println(createErr)
		return
	}

	databaseFile := googleDriveFolder + `/` + constants.DatabaseFile
	databaseBackend := service.DatabaseFactory.NewDatabase(databaseFile)

	cryptoSvc := service.CryptoFactory.GetInstance()

	if err := databaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
		var conf models.GoogleDriveConfig

		if err := connection.ReadConfig(constants.ConfigDatabaseKey, &conf); err != nil {
			if folderId == `` {
				return fmt.Errorf("ERROR!!! Specify a folder-id to upload to/download from")
			}

			conf.Id = constants.ConfigDatabaseKey
			conf.FolderId = folderId

			if cryptoSvc.IsEnabled() {
				conf.Password = sql.NullString{
					String: base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(cryptoSvc.GetPassword())),
					Valid:  true,
				}
			} else {
				conf.Password = sql.NullString{Valid: false}
			}

			conf.TrashFolderId = sql.NullString{Valid: false}

			conf.Encoder = sql.NullString{String: cryptoSvc.GetEncoder(), Valid: true}
			conf.Decoder = sql.NullString{String: cryptoSvc.GetDecoder(), Valid: true}

			connection.SaveConfig(&conf)
		} else {
			//fmt.Println(utils.Detail(&conf))
			password := ``

			if conf.Password.Valid {
				if decodeString, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(conf.Password.String); err != nil {
					fmt.Println(err)
				} else {
					password = string(decodeString)
				}
			}

			if folderId == `` {
				folderId = conf.FolderId

				fmt.Printf("Upload to / Download from [ %s ]\n", folderId)
			} else if folderId != conf.FolderId {
				return fmt.Errorf("ERROR! upload/download folder was %s (vs %s)", conf.FolderId, folderId)
			}

			if cryptoSvc.IsEnabled() {
				if password == `` {
					fmt.Println("WARNING!! Folder was NOT encrypted with passwords, encrypt new files")
					conf.Password = sql.NullString{
						String: base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(cryptoSvc.GetPassword())),
						Valid:  true,
					}
					connection.SaveConfig(&conf)
				} else if password != cryptoSvc.GetPassword() {
					return fmt.Errorf("ERROR!! Folder [ %s ] encrypted with different password", folderName)
				}
			} else {
				if password != `` {
					fmt.Println("WARNING!! Folder was encrypted with a password, use stored password")
					cryptoSvc.SetPassword(password)
				}
			}
		}
		service.updateEnvironment(&conf, folderName)

		return service.Api.Connect(func(service *drive.Service) error {
			callback(connection, service, folderId)
			return nil
		})
	}); err != nil {
		logger.Error(err)
	}
}
