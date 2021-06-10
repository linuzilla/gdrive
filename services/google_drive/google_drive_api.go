package google_drive

import (
	"google.golang.org/api/drive/v3"
	"io"
	"time"
)

type API interface {
	RootFolder() string
	GetCredentialEmail() (string, error)
	SetImpersonate(subject string)
	GetImpersonate() string
	SetDomainAdminAccess(useDomainAdminAccess bool)
	GetDomainAdminAccess() bool
	FindWithQuery(service *drive.Service, query string, callback func(file *drive.File) bool) error
	FindFilesByNameInFolder(service *drive.Service, folderId string, fileName string, callback func(file *drive.File) bool) error
	FindDirectoriesInFolder(service *drive.Service, folderId string, baseDir string, callback func(file *drive.File, path string) bool) error
	RecursiveFindDirectoriesInFolder(service *drive.Service, folderId string, baseDir string, callback func(file *drive.File, path string)) error
	FindFilesInFolder(service *drive.Service, folderId string, callback func(file *drive.File) bool) error
	Connect(callback func(service *drive.Service) error) error
	FileInfo(service *drive.Service, fileId string) (*drive.File, error)
	Read(service *drive.Service, folderId string, fileNameDecoder func(file *drive.File) string) error
	ReadFolder(service *drive.Service, folderId string, callback func(file *drive.File)) error
	ReadFolderDetail(service *drive.Service, folderId string, allDrive bool, callback func(file *drive.File) bool) error
	Delete(service *drive.Service, fileId string) error
	UploadFile(service *drive.Service, localFile string, folderId string, fileName string, readerWrapper func(reader io.Reader) io.Reader) (*drive.File, error)

	ReplaceFile(service *drive.Service, localFile string, source *drive.File, fileName string) error
	ReadTeamDrives(service *drive.Service, callback func(teamDrive *drive.TeamDrive) bool) error
	DownloadFile(service *drive.Service, fileId string, localFile string, overwrite bool) (string, error)
	RetrieveFile(service *drive.Service, fileId string, fileRetriever func(*drive.File, io.Reader) error) error
	PushFile(service *drive.Service, folderId string, fileName string, reader io.Reader) (*drive.File, error)
	CreateFolder(service *drive.Service, parentFolder string, folderName string) (string, error)
	MoveFile(service *drive.Service, file *drive.File, newParentId string) (*drive.File, error)
	FindOrCreateTrashFolder(service *drive.Service, fileId string) (string, error)
	Rename(service *drive.Service, fileId string, newFileName string) (*drive.File, error)
	Duplicate(service *drive.Service, sourceFileId string, targetFileName string) (*drive.File, error)
	ShareInfo(service *drive.Service, fileId string) (*drive.File, error)
	ShareFileWith(service *drive.Service, fileId string, email string, role string) (*drive.Permission, error)
	RemovePermission(service *drive.Service, fileId string, permissionId string) error
	PermissionInfo(service *drive.Service, driveId string, callback func(permission *drive.Permission)) error
	GetPermissionIdByEmail(service *drive.Service, fileId string, email string) (string, error)
}

func ToLocalTime(timeString string) (*time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, timeString); err != nil {
		return nil, err
	} else {
		in := parsed.In(time.Now().Location())
		return &in, nil
	}
}
