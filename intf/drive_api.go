package intf

import (
	"github.com/linuzilla/gdrive/models"
	"google.golang.org/api/drive/v3"
	"io"
)

type DriveAPI interface {
	ReadTeamDrives(callback func(teamDrive *drive.TeamDrive) bool) error
	ChangeFolder(folderNameOrId string)
	FileInfo(fileNameOrId string, fileNameDecoder func(*drive.File) string) (string, *drive.File, error)
	List()
	ListFolder(folderId string)
	ReadFolderDetail(folderId string, allDrive bool, callback func(file *drive.File) bool) error
	ListWithDecoder(fileNameDecoder func(*drive.File) string)
	Upload(localFile string) error
	//Download(fileNameOrId string, localFile string) (string, error)
	RetrieveFile(fileNameOrId string, fileRetriever func(*drive.File, io.Reader) error) error
	CreateFolder(folderName string)
	CurrentFolder() string
	Delete(fileNameOrId string)
	DropFile(fileNameOrId string)
	Pushd()
	Popd()
	Dirs()
	Version()
	Duplicate(fileNameOrId string, fileName string)
	Rename(fileNameOrId string, fileName string)
	FindOrCreateTrashFolder(fileId string) (string, error)
	EncryptedUpload(cryptoService CryptoService, fileFullPath string,
		folderId string, fileName string, md5sum string) error
	ShareInfo(fileNameOrId string) error
	ShareFileWith(fileId string, email string, role models.Role) error
	UnshareFileWith(fileId string, email string) error
	RemovePermission(fileId string, permissionId string) error
	PermissionInfo(driveId string, callback func(permission *drive.Permission)) error
	SetDomainAdminAccess(useDomainAdminAccess bool)
	GetDomainAdminAccess() bool
	SetImpersonate(impersonate string)
	GetImpersonate() string
	GetCredentialEmail() (string, error)

	DownloadFile(fileNameOrId string, dir string, localFile string) error
	NativeDownloadFile(driveSvc *drive.Service, fileNameOrId string, dir string, localFile string) error
}
