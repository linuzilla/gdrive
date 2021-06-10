package drive_api

import (
	"fmt"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
)

func (api *driveApiImpl) CreateFolder(folderName string) {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		if folderId, err := api.DriveApi.CreateFolder(service, api.folderId, folderName); err == nil {
			fmt.Println("Push current folder to stack: " + api.folderId)
			fmt.Println("Create and change folder to: " + folderId + " (" + folderName + ")")
			api.folderStack.Push(api.folderId)
			api.folderId = folderId
		} else {
			logger.Error("%v", err)
		}
		return nil
	}); err != nil {
		logger.Error("%v", err)
	}
}
