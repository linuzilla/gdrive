package drive_api

import (
	"bufio"
	"fmt"
	"github.com/linuzilla/gdrive/services/progress_service"
	"github.com/linuzilla/gdrive/utils"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
	"path"
)

func (api *driveApiImpl) Download(fileNameOrId string, localFile string) (string, error) {
	fileName := ``

	err := api.DriveApi.Connect(func(service *drive.Service) error {
		if name, err := api.DriveApi.DownloadFile(service, api.cachedName(fileNameOrId), localFile, false); err != nil {
			return err
		} else {
			fileName = name
			return nil
		}
	})

	return fileName, err
}

func (api *driveApiImpl) RetrieveFile(fileNameOrId string, fileRetriever func(*drive.File, io.Reader) error) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.DriveApi.RetrieveFile(service, api.cachedName(fileNameOrId), fileRetriever)
	})
}

func (api *driveApiImpl) DownloadFile(fileNameOrId string, dir string, localFile string) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.NativeDownloadFile(service, fileNameOrId, dir, localFile)
	})
}

func (api *driveApiImpl) NativeDownloadFile(driveSvc *drive.Service, fileNameOrId string, dir string, localFile string) error {
	progressService := api.ProgressFactory.NewInstance(`download`)
	cryptoSvc := api.CryptoFactory.GetInstance()

	defer progressService.Close()

	var md5sum string
	var fileName string

	fileId := api.cachedName(fileNameOrId)

	err := progressService.ExecAndWait(func(progressSvc progress_service.ProgressService, errChan chan error) {
		errChan <- api.DriveApi.RetrieveFile(driveSvc, fileId, func(file *drive.File, reader io.Reader) error {
			var err error
			fileName, md5sum, err = cryptoSvc.DecryptFileNameWithMd5(file)
			progressSvc.SetName(fileName)
			progressSvc.SetSize(file.Size)
			r := progressSvc.WrapperReader(reader)

			if err != nil {
				return err
			}

			if localFile == `.` {
				localFile = path.Join(dir, fileName)
			}

			if cryptoSvc.IsEnabled() {
				err = cryptoSvc.DecodeAndDownloadFromReader(localFile, r)
			} else {
				f, errOpenFile := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)

				if errOpenFile != nil {
					return errOpenFile
				}

				defer f.Close()

				writer := bufio.NewWriter(f)
				_, err = io.Copy(writer, r)
				writer.Flush()
			}

			return err
		})
	})

	if err == nil {
		if md5, err := utils.Md5sum(localFile); err != nil {
			return err
		} else if md5 != md5sum {
			return fmt.Errorf("checksum mismatch (%s vs %s)", md5, md5sum)
		}
		fmt.Printf("%s: checksum verified\n", fileName)
		return nil
	} else {
		return err
	}
}
