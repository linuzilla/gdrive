package drive_api

import "google.golang.org/api/drive/v3"

func (api *driveApiImpl) FindOrCreateTrashFolder(fileId string) (string, error) {
	trashFolder := ``
	err := api.DriveApi.Connect(func(svc *drive.Service) error {
		var err error
		trashFolder, err = api.DriveApi.FindOrCreateTrashFolder(svc, fileId)
		return err
	})
	return trashFolder, err
}
