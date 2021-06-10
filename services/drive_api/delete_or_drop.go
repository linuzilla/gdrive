package drive_api

import (
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
)

func (api *driveApiImpl) Delete(fileNameOrId string) {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		if err := api.DriveApi.Delete(service, api.cachedName(fileNameOrId)); err != nil {
			logger.Error("%v", err)
		}
		return nil
	}); err != nil {
		logger.Error(err)
	}
}

func (api *driveApiImpl) DropFile(fileNameOrId string) {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		fileId := api.cachedName(fileNameOrId)
		if trashFolder, err := api.DriveApi.FindOrCreateTrashFolder(service, fileId); err != nil {
			return err
		} else {
			if file, err := api.DriveApi.FileInfo(service, fileId); err != nil {
				return err
			} else {
				_, err := api.DriveApi.MoveFile(service, file, trashFolder)

				return err
			}
		}
	}); err != nil {
		logger.Error(err)
	}
}
