package drive_api

import "google.golang.org/api/drive/v3"

func (api *driveApiImpl) FileInfo(fileNameOrId string, fileNameDecoder func(*drive.File) string) (string, *drive.File, error) {
	var fileInfo *drive.File
	var fileName string

	err := api.DriveApi.Connect(func(service *drive.Service) error {
		if file, err := api.DriveApi.FileInfo(service, api.cachedName(fileNameOrId)); err != nil {
			return err
		} else {
			fileName = file.Name
			fileInfo = file

			if fileNameDecoder != nil {
				fileName = fileNameDecoder(file)
			}
			return nil
		}
	})
	return fileName, fileInfo, err
}
