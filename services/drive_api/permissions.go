package drive_api

import (
	"fmt"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
)

func (api *driveApiImpl) PermissionInfo(driveId string, callback func(permission *drive.Permission)) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.DriveApi.PermissionInfo(service, driveId, callback)
	})
}

func (api *driveApiImpl) ShareInfo(fileNameOrId string) error {
	err := api.DriveApi.Connect(func(service *drive.Service) error {
		if _, err := api.DriveApi.ShareInfo(service, api.cachedName(fileNameOrId)); err != nil {
			return err
		} else {

			return nil
		}
	})
	return err
}

func (api *driveApiImpl) ShareFileWith(fileId string, email string, role models.Role) error {
	if err := role.IsValid(); err != nil {
		return err
	} else {
		return api.DriveApi.Connect(func(service *drive.Service) error {
			perm, err := api.DriveApi.ShareFileWith(service, fileId, email, string(role))
			if err != nil {
				fmt.Printf("Permission: %s (%s) [ %s ]", perm.EmailAddress, perm.Role, perm.Id)
			}
			return err
		})
	}
}

func (api *driveApiImpl) UnshareFileWith(fileId string, email string) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		if permissionId, err := api.DriveApi.GetPermissionIdByEmail(service, fileId, email); err != nil {
			return err
		} else {
			logger.Notice("email: %s, permission-id: %s", email, permissionId)
			return api.DriveApi.RemovePermission(service, fileId, permissionId)
		}
	})
}

func (api *driveApiImpl) RemovePermission(fileId string, permissionId string) error {
	return api.DriveApi.Connect(func(service *drive.Service) error {
		return api.DriveApi.RemovePermission(service, fileId, permissionId)
	})
}
