package drive_api

import (
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
)

func (api *driveApiImpl) Duplicate(fileNameOrId string, fileName string) {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		fileId := api.cachedName(fileNameOrId)
		_, err := api.DriveApi.Duplicate(service, fileId, fileName)
		return err
	}); err != nil {
		logger.Error(err)
	}
}
