package drive_api

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
)

func (api *driveApiImpl) Upload(localFile string) error {
	if fileState, err := os.Stat(localFile); err != nil {
		return err
	} else {
		return api.DriveApi.Connect(func(service *drive.Service) error {
			progressReader := utils.NewProgressReader(nil)
			done := make(chan error)

			go func() {
				_, err := api.DriveApi.UploadFile(service, localFile, api.folderId, localFile, func(reader io.Reader) io.Reader {
					return progressReader.SetReader(reader)
				})
				done <- err
			}()

			return commons.UploadProgress(localFile, fileState.Size(), done, progressReader)
		})
	}
}

func (api *driveApiImpl) EncryptedUpload(cryptoService intf.CryptoService, fileFullPath string,
	folderId string, fileName string, md5sum string) error {

	return api.DriveApi.Connect(func(service *drive.Service) error {
		return commons.UploadFile(api.DriveApi, service, cryptoService, fileFullPath, folderId, fileName, md5sum)
	})
}

func (api *driveApiImpl) UploadWithProgress(service *drive.Service, fullPath string, folderId string, fileName string, md5sum string) error {
	fileState, err := os.Stat(fullPath)

	if err != nil {
		return err
	}
	fmt.Printf("Upload [ %s ] (%d)\n", fullPath, fileState.Size())

	cryptoService := api.CryptoFactory.GetInstance()

	done := make(chan error)
	progressReader := utils.NewProgressReader(nil)

	if cryptoService.IsEnabled() {
		encryptedFileName := cryptoService.EncryptFileNameWithMd5(fileName, md5sum)

		go func() {
			done <- cryptoService.EncodeAndUploadReader(fullPath, func(writer io.Writer) io.Writer {
				return writer
			}, func(reader io.Reader) error {
				progressReader.SetReader(reader)

				file, err := api.DriveApi.PushFile(service, folderId, encryptedFileName, progressReader)

				if err == nil {
					logger.Debug("file pushed: %s", file.Id)
				}
				return err
			})
		}()
	} else {
		go func() {
			_, err := api.DriveApi.UploadFile(service, fullPath, folderId, fileName, func(reader io.Reader) io.Reader {
				return progressReader.SetReader(reader)
			})
			done <- err
		}()
	}
	return commons.UploadProgress(fileName, fileState.Size(), done, progressReader)
}
