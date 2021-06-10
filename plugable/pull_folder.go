package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/sync_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"os"
	"strings"
)

type PullFolderCommand struct {
	SyncService sync_service.SyncService               `inject:"*"`
	EnvSvc      environment_service.EnvironmentService `inject:"*"`
	DotService  dot_gdrive_service.DotGDriveService    `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*PullFolderCommand)(nil)

func (PullFolderCommand) Command() string {
	return `pull-folder`
}

//func (cmd *PullFolderCommand) prepare(folderName string,
//	folderId string,
//	callback func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId, workingDir string)) {
//
//	commons.PrepareBeforeSync(
//		cmd.GoogleDrive,
//		cmd.DatabaseFactory,
//		cmd.CryptoService,
//		folderName,
//		folderId,
//		callback)
//}

func (cmd *PullFolderCommand) pullFolder(folderName string, downloadFolderId string) {
	fmt.Printf("Download folder [ %s ] to %s\n", folderName, downloadFolderId)

	cmd.SyncService.PullDownload(folderName, downloadFolderId)

	//cmd.DotService.PrepareBeforeSync(folderName, downloadFolderId, func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId string, workingDir string) {
	//	cryptoSvc := cmd.CryptoFactory.GetInstance()
	//
	//	cmd.Factory.NewInstance(&traversal_service.TraversalHandler{
	//		UploadFile: func(filePath string, folderId string, remoteFileName string) error {
	//			return nil
	//		},
	//		//UpdateRemoteFile: func(filePath string, remoteFileName string, remote *drive.File) error {
	//		//	return nil
	//		//},
	//		DropToTrash: func(file *drive.File, trashFolderId string) error {
	//			return nil
	//		},
	//		MarkUploaded: func(data *models.SyncFileInfo) {
	//		},
	//		//FileEncryptor: func(dirName string, fileName string, md5sum string, callback func(uploadFileName, remoteFileName string) error) error {
	//		//	return commons.EncryptFile(workingDir, cryptoSvc, path.Join(dirName, fileName), md5sum, callback)
	//		//},
	//		Md5Fetcher: func(baseDir string, relativeDir string, fileInfo os.FileInfo, data *models.SyncFileInfo) (s string, err error) {
	//			return commons.Md5WithCache(connection, baseDir, relativeDir, fileInfo, data)
	//		},
	//		DownloadRemoteFile: func(remoteFile *drive.File, localFilePath string, md5sum string, encrypted bool) error {
	//			//fmt.Printf("download remote file %s\n", remoteFile.Name)
	//
	//			commons.DownloadAnDecodeRemoteFile(
	//				cmd.GoogleDrive,
	//				service,
	//				cryptoSvc,
	//				cmd.EnvSvc.GetWorkingDirectory(),
	//				remoteFile,
	//				localFilePath,
	//				md5sum,
	//				encrypted)
	//			return nil
	//		},
	//		DeleteBeforeDownload: func(localFilePath string, data *models.SyncFileInfo) {
	//			fmt.Printf("delete before download(%s)\n", localFilePath)
	//			//connection.Delete(data)
	//			//os.Remove(localFilePath)
	//		},
	//		Mkdir: func(dirPath string) error {
	//			fmt.Printf("mkdir(%s)\n", dirPath)
	//			return os.Mkdir(dirPath, 0700)
	//		},
	//	}).RecursiveRemote(service, folderName, `/`, folderId, true)
	//})

}
func (cmd *PullFolderCommand) Execute(args ...string) int {
	downloadFolderId := ``

	switch len(args) {
	default:
		if dotEnv := cmd.DotService.LoadEnvironment(); dotEnv != nil {
			cmd.pullFolder(dotEnv.Dir, dotEnv.FolderId)
		} else {
			on := constants.OnOffMap[cmd.EnvSvc.Encoding()]

			fmt.Printf("usage: pull-folder <local-folder-name> <folder-id>    [Encryption: %s]\n", on)
		}
	case 2:
		downloadFolderId = args[1]
		fallthrough
	case 1:
		folderName := utils.StripQuotedString(args[0])

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return 0
		}

		if folderName == `.` {
			folderName = currentDir
		} else if !strings.HasPrefix(folderName, `/`) {
			folderName = currentDir + `/` + folderName
		}

		if strings.HasSuffix(folderName, `/`) {
			folderName = folderName[:len(folderName)-1]
		}

		cmd.pullFolder(folderName, downloadFolderId)
	}
	return 0
}
