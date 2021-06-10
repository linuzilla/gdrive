package drive_api

import (
	"fmt"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
)

func (api *driveApiImpl) List() {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		if err := api.DriveApi.Read(service, api.folderId, nil); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		logger.Error("%v", err)
	}
}

func (api *driveApiImpl) ListFolder(folderId string) {
	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		if err := api.DriveApi.Read(service, api.cachedName(folderId), nil); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		logger.Error("%v", err)
	}
}

func (api *driveApiImpl) ReadFolderDetail(folderId string, allDrive bool, callback func(file *drive.File) bool) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.DriveApi.ReadFolderDetail(service, folderId, allDrive, callback)
	})
}
func (api *driveApiImpl) ListWithDecoder(fileNameDecoder func(*drive.File) string) {
	fileMap := make(map[string]string)

	if err := api.DriveApi.Connect(func(service *drive.Service) error {
		if err := api.DriveApi.ReadFolder(service, api.folderId, func(file *drive.File) {
			fileName := fileNameDecoder(file)
			fileMap[fileName] = file.Id

			if file.MimeType == google_drive.FolderMimeType {
				fmt.Printf("  [ %s ] <dir> %s\n", file.Id, fileName)
			} else {
				fmt.Printf("  [ %s ] ----- %s  [ %d ]\n", file.Id, fileName, file.Size)
			}
		}); err != nil {
			fmt.Println(err)
		} else {
			api.filesCache = fileMap
		}

		return nil
	}); err != nil {
		logger.Error("%v", err)
	}
}
