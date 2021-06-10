package traversal_service

import (
	"fmt"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func (traversal *traversalServiceImpl) retrieveRemoteFolders(baseDir string, relativeDir string,
	folderId string) (folders map[string]*drive.File) {
	folders = make(map[string]*drive.File)

	if err := traversal.gdrive.FindDirectoriesInFolder(traversal.driveSvc, folderId, relativeDir, func(file *drive.File, path string) bool {
		//fmt.Println(baseDir + " " + file.Name)
		folders[file.Name] = file
		return true
	}); err != nil {
		log.Printf("error: %v", err)
		traversal.terminate = true
	}

	return folders
}

func (traversal *traversalServiceImpl) retrieveRemoteFiles(folderId string) (remoteFiles map[string]*drive.File) {
	remoteFiles = make(map[string]*drive.File)

	if err := traversal.gdrive.FindFilesInFolder(traversal.driveSvc, folderId, func(file *drive.File) bool {
		remoteFiles[file.Name] = file
		return true
	}); err != nil {
		log.Printf("error: (folder= %s) %v", folderId, err)
		traversal.terminate = true
	}

	return remoteFiles
}

func (traversal *traversalServiceImpl) retrieveLocal(baseDir string) (dirs, files map[string]os.FileInfo) {
	dirs = make(map[string]os.FileInfo)
	files = make(map[string]os.FileInfo)

	if allFiles, err := ioutil.ReadDir(baseDir); err != nil {
		log.Printf("error: (%s) %v", baseDir, err)
		traversal.terminate = true
	} else {
		for _, file := range allFiles {
			if file.IsDir() {
				dirs[file.Name()] = file
			} else {
				files[file.Name()] = file
			}
		}
	}

	return dirs, files
}

func (traversal *traversalServiceImpl) uploadFile(baseDir string, fileName string, node *FileNode, folderId string) {
	fullName := path.Join(baseDir, fileName)

	fmt.Printf("Upload file [ %s ] ... ", fullName)

	if err := commons.UploadFile(traversal.gdrive, traversal.driveSvc, traversal.cryptoService,
		fullName, folderId, fileName, node.md5sum); err != nil {
		logger.Error(err)
	} else {
		traversal.uploadCount++
		traversal.uploadSize += node.fileInfo.Size()
		traversal.handler.MarkUploaded(node.data)
		fmt.Printf("%d bytes\n", node.fileInfo.Size())
	}

	//if traversal.cryptoService.IsEnabled() {
	//	encryptFileName := traversal.cryptoService.EncryptFileNameWithMd5(fileName, node.md5sum)
	//	logger.Debug("Encrypted Name: %s", encryptFileName)
	//
	//	if err := traversal.cryptoService.EncodeAndUploadReader(fullName, func(reader io.Reader) error {
	//		file, err := traversal.gdrive.PushFile(traversal.driveSvc, folderId, encryptFileName, reader)
	//
	//		if err == nil {
	//			logger.Debug("file pushed: %s", file.Id)
	//		}
	//		return err
	//	}); err != nil {
	//		logger.Error(err)
	//	} else {
	//		fmt.Printf("%d bytes (encrypted)\n", node.fileInfo.Size())
	//		traversal.uploadCount++
	//		traversal.uploadSize += node.fileInfo.Size()
	//		traversal.handler.MarkUploaded(node.data)
	//	}
	//} else {
	//	if err := traversal.handler.UploadFile(fullName, folderId, fileName); err != nil {
	//		logger.Error(err)
	//	} else {
	//		fmt.Printf("%d bytes\n", node.fileInfo.Size())
	//		traversal.uploadCount++
	//		traversal.uploadSize += node.fileInfo.Size()
	//		traversal.handler.MarkUploaded(node.data)
	//	}
	//}

	//fmt.Printf("Upload file [ %s/%s ] ... ", baseDir, fileName)

	// commons.DownloadAnDecodeRemoteFile(traversal.gdrive, traversal.driveSvc, traversal.cryptoService, `/tmp`, )
	//if err := traversal.handler.FileEncryptor(baseDir, fileName, node.md5sum, func(uploadFileName, remoteFileName string) error {
	//	if err := traversal.handler.UploadFile(uploadFileName, folderId, remoteFileName); err != nil {
	//		return err
	//	} else {
	//		fmt.Printf("%d bytes\n", node.fileInfo.Size())
	//		traversal.uploadCount++
	//		traversal.uploadSize += node.fileInfo.Size()
	//		traversal.handler.MarkUploaded(node.data)
	//		return nil
	//	}
	//}); err != nil {
	//	fmt.Println(err)
	//}
}
