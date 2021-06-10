package drive_api

import "google.golang.org/api/drive/v3"

func (api *driveApiImpl) ReadTeamDrives(callback func(teamDrive *drive.TeamDrive) bool) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.DriveApi.ReadTeamDrives(service, callback)
	})
}
