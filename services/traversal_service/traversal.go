package traversal_service

import (
	"github.com/linuzilla/gdrive/models"
	"google.golang.org/api/drive/v3"
	"os"
)

type TraversalService interface {
	RecursiveLocal(driveSvc *drive.Service, baseDir string, relativeDir string, folderId string, speedUpUpload bool)
	RecursiveRemote(driveSvc *drive.Service, baseDir string, relativeDir string, folderId string, speedUpload bool)
}

type FileNode struct {
	fileInfo os.FileInfo
	md5sum   string
	found    bool
	data     *models.SyncFileInfo
}

type TraversalHandler struct {
	//UploadFile           func(filePath string, folderId string, remoteFileName string) error
	//UpdateRemoteFile     func(filePath string, remoteName string, remote *drive.File) error
	// FileEncryptor        func(dirName string, fileName string, md5sum string, callback func(uploadFileName, remoteFileName string) error) error
	DropToTrash          func(file *drive.File, trashFolderId string) error
	MarkUploaded         func(data *models.SyncFileInfo)
	Md5Fetcher           func(baseDir string, relativeDir string, fileInfo os.FileInfo, data *models.SyncFileInfo) (string, error)
	DownloadRemoteFile   func(remoteFile *drive.File, localFilePath string, md5sum string, encrypted bool) error
	DeleteExtraFiles     func(fileName string, remoteFile *drive.File, trashFolderId string)
	DeleteBeforeDownload func(localFilePath string, data *models.SyncFileInfo)
	Mkdir                func(dirPath string) error
}
