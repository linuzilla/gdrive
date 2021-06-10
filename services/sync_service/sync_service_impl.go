package sync_service

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/traversal_service"
	"google.golang.org/api/drive/v3"
	"os"
	"path/filepath"
)

type syncServiceImp struct {
	DriveAPI      intf.DriveAPI                             `inject:"*"`
	EnvSvc        environment_service.EnvironmentService    `inject:"*"`
	Factory       traversal_service.TraversalServiceFactory `inject:"*"`
	GoogleDrive   google_drive.API                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory              `inject:"*"`
	DotService    dot_gdrive_service.DotGDriveService       `inject:"*"`
	//InterruptSvc  interrupt_service.InterruptService        `inject:"*"`
}

func New() SyncService {
	return &syncServiceImp{}
}

func (syncSvc *syncServiceImp) PushUpload(folderName string, uploadFolderId string, speedUpload bool, deleteExtraFiles bool) {
	syncSvc.DotService.PrepareBeforeSync(folderName, uploadFolderId, func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId string) {
		//cryptoSvc := cmd.CryptoFactory.GetInstance()

		syncSvc.Factory.NewInstance(&traversal_service.TraversalHandler{
			//UploadFile: func(filePath string, folderId string, remoteFileName string) error {
			//	_, err := syncSvc.GoogleDrive.UploadFile(service, filePath, folderId, remoteFileName)
			//	return err
			//},
			DropToTrash: func(file *drive.File, trashFolderId string) error {
				// fmt.Printf("Drop to trash %s (trash: %s)\n", file.Name, trashFolderId)
				_, err := syncSvc.GoogleDrive.MoveFile(service, file, trashFolderId)
				return err
			},
			MarkUploaded: func(data *models.SyncFileInfo) {
				data.Uploaded = true
				connection.SaveOrUpdate(data)
			},
			Md5Fetcher: func(baseDir string, relativeDir string, fileInfo os.FileInfo, data *models.SyncFileInfo) (s string, err error) {
				return commons.Md5WithCache(connection, baseDir, relativeDir, fileInfo, data)
			},
			DownloadRemoteFile: func(remoteFile *drive.File, localFilePath string, md5sum string, encrypted bool) error {
				return nil
			},
			DeleteExtraFiles: func(fileName string, remoteFile *drive.File, trashFolderId string) {
				if deleteExtraFiles {
					fmt.Printf("Extra file on remote %s (%s) ... ", fileName, remoteFile.Id)

					_, err := syncSvc.GoogleDrive.MoveFile(service, remoteFile, trashFolderId)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(`ok`)
					}
				} else {
					fmt.Printf("Extra file on remote %s (%s)  (use -d to delete it)\n", fileName, remoteFile.Id)
				}
			},
			DeleteBeforeDownload: func(localFilePath string, data *models.SyncFileInfo) {
			},
			Mkdir: func(dirPath string) error {
				return nil
			},
		}).RecursiveLocal(service, folderName, `/`, folderId, speedUpload)
	})
}

func (syncSvc *syncServiceImp) PullDownload(folderName string, downloadFolderId string) {
	syncSvc.DotService.PrepareBeforeSync(folderName, downloadFolderId, func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId string) {
		//cryptoSvc := syncSvc.CryptoFactory.GetInstance()

		syncSvc.Factory.NewInstance(&traversal_service.TraversalHandler{
			//UploadFile: func(filePath string, folderId string, remoteFileName string) error {
			//	return nil
			//},
			DropToTrash: func(file *drive.File, trashFolderId string) error {
				return nil
			},
			MarkUploaded: func(data *models.SyncFileInfo) {
			},
			Md5Fetcher: func(baseDir string, relativeDir string, fileInfo os.FileInfo, data *models.SyncFileInfo) (s string, err error) {
				return commons.Md5WithCache(connection, baseDir, relativeDir, fileInfo, data)
			},
			DownloadRemoteFile: func(remoteFile *drive.File, localFilePath string, md5sum string, encrypted bool) error {
				//fmt.Printf("download remote file %s\n", remoteFile.Name)

				dir := filepath.Dir(localFilePath)

				return syncSvc.DriveAPI.NativeDownloadFile(service, remoteFile.Id, dir, localFilePath)
				//commons.DownloadAnDecodeRemoteFile(
				//	syncSvc.GoogleDrive,
				//	service,
				//	cryptoSvc,
				//	remoteFile,
				//	localFilePath,
				//	md5sum,
				//	encrypted)
				//return nil
			},
			DeleteExtraFiles: func(fileName string, remoteFile *drive.File, trashFolderId string) {
			},
			DeleteBeforeDownload: func(localFilePath string, data *models.SyncFileInfo) {
				backupFile := localFilePath + constants.BackUpFileExtension
				for i := 1; ; i += 1 {
					if _, err := os.Stat(backupFile); os.IsNotExist(err) {
						break
					}
					backupFile = fmt.Sprintf(`%s.%d%s`, localFilePath, i, constants.BackUpFileExtension)
				}
				fmt.Printf("Backup [ %s ] as [ %s ]\n", localFilePath, filepath.Base(backupFile))
				os.Rename(localFilePath, backupFile)
				connection.Delete(data)
				//os.Remove(localFilePath)
			},
			Mkdir: func(dirPath string) error {
				fmt.Printf("mkdir(%s)\n", dirPath)
				return os.Mkdir(dirPath, 0700)
			},
		}).RecursiveRemote(service, folderName, `/`, folderId, true)
	})
}
