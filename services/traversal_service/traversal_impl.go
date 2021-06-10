package traversal_service

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/ignore_service"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"strings"
)

type traversalServiceImpl struct {
	driveSvc         *drive.Service
	cryptoService    intf.CryptoService
	gdrive           google_drive.API
	ignoreService    ignore_service.IgnoreService
	fileNameGuard    ignore_service.FileNameGuard
	handler          *TraversalHandler
	terminate        bool
	totalDirectories int
	scanCount        int
	scanSize         int64
	uploadCount      int
	downloadCount    int
	uploadSize       int64
	downloadSize     int64
	speedUpUpload    bool
	trashFolderId    string
}

func (traversal *traversalServiceImpl) RecursiveLocal(driveSvc *drive.Service, baseDir string, relativeDir string, folderId string, speedUpload bool) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logger.Error("panic: %v\n", r)
		}
	}()

	traversal.fileNameGuard = traversal.ignoreService.LoadRules(baseDir)

	traversal.driveSvc = driveSvc
	traversal.speedUpUpload = speedUpload

	if trashFolderId, err := traversal.gdrive.FindOrCreateTrashFolder(driveSvc, folderId); err != nil {
		fmt.Println(err)
	} else {
		traversal.trashFolderId = trashFolderId

		traversal.recursiveLocal(baseDir, relativeDir, folderId, 0)
	}
}

func (traversal *traversalServiceImpl) recursiveLocal(baseDir string, relativeDir string, folderId string, depth int) {
	files, err := ioutil.ReadDir(baseDir)

	if err != nil {
		traversal.terminate = true
		panic(err.Error())
	}

	subDirMap := make(map[string]string)
	fileInfoMap := make(map[string]*FileNode)

	allUploaded := true

	for _, file := range files { // iterate all files in local directory
		if traversal.fileNameGuard != nil {
			if rule, ignore := traversal.fileNameGuard.ShouldIgnore(relativeDir, file.Name()); ignore {
				fmt.Printf("file: %s ... ignored by rule %d (%s)\n", path.Join(relativeDir, file.Name()), rule.Index, rule.Filename)
				continue
			}
		}

		if file.IsDir() {
			if depth != 0 || file.Name() != constants.GoogleDriveFolder {
				subDirMap[file.Name()] = ``
				traversal.totalDirectories++
			}
		} else if file.Mode().IsRegular() {
			if !strings.HasSuffix(file.Name(), constants.BackUpFileExtension) {
				var data models.SyncFileInfo

				traversal.scanCount++
				traversal.scanSize += file.Size()

				if checksum, err := traversal.handler.Md5Fetcher(baseDir, relativeDir, file, &data); err != nil {
					log.Printf("error: retrieveMd5sum(%s) %v", file.Name(), err)
				} else {
					node := &FileNode{fileInfo: file, md5sum: checksum, found: false, data: &data}

					fileInfoMap[file.Name()] = node

					if !node.data.Uploaded {
						allUploaded = false
					}
				}
			}
		}
	}

	if err := traversal.gdrive.FindDirectoriesInFolder(traversal.driveSvc, folderId, relativeDir, func(file *drive.File, path string) bool {
		if realFileName, _, err := traversal.cryptoService.DecryptFileNameWithMd5(file); err != nil {
			log.Printf("error: (find folder %s) %v", folderId, err)
		} else {
			if v, found := subDirMap[realFileName]; found {
				if v != `` {
					fmt.Printf(">> duplicate folder: %s  %s\n", relativeDir, realFileName)
				} else {
					subDirMap[realFileName] = file.Id
					// fmt.Printf("Add Dir: %s (%s)\n", file.Name, file.Id)
				}
			} else {
				traversal.handler.DeleteExtraFiles(realFileName, file, traversal.trashFolderId)
				// handler.
				//fmt.Printf(">> %s %s: extra folder needs to be deleted\n", relativeDir, realFileName)
			}
		}
		return true
	}); err != nil {
		logger.Error("error: (find folder %s) %v", folderId, err)
	}

	if traversal.speedUpUpload && allUploaded {
		fmt.Printf("directory: %s (done)\n", baseDir)
	} else {
		if err := traversal.gdrive.FindFilesInFolder(traversal.driveSvc, folderId, func(file *drive.File) bool {
			if traversal.terminate {
				return false
			}

			if realFileName, md5sum, err := traversal.cryptoService.DecryptFileNameWithMd5(file); err != nil {
				log.Printf("error: %v", err)
			} else {
				if node, found := fileInfoMap[realFileName]; found {
					node.found = true

					if node.md5sum != md5sum {
						if fileModifiedTime, err := google_drive.ToLocalTime(file.ModifiedTime); err != nil {
							logger.Error(err)
							panic(err.Error())
						} else {
							if fileModifiedTime.Before(node.fileInfo.ModTime()) {
								fmt.Printf("File modified: %s/%s ... (move to trash folder) ... ", baseDir, realFileName)
								if err := traversal.handler.DropToTrash(file, traversal.trashFolderId); err != nil {
									logger.Error(err)
									panic(err.Error())
								} else {
									fmt.Println("ok")
									traversal.uploadFile(baseDir, realFileName, node, folderId)
								}
							} else {
								fmt.Printf("File modified: %s/%s ... skip (older)\n", baseDir, realFileName)
							}
						}

					} else {
						fmt.Printf("file: %s (ok)\n", path.Join(relativeDir, realFileName))
						traversal.handler.MarkUploaded(node.data)
					}
					//} else {
					//	if traversal.handler.Downloader(file, path.Join(baseDir, realFileName)) {
					//		traversal.downloadCount++
					//		traversal.downloadSize += file.Size
					//	}
				} else {
					traversal.handler.DeleteExtraFiles(realFileName, file, traversal.trashFolderId)
				}
			}
			return true
		}); err != nil {
			logger.Error("error: %v", err)
		}

		for fileName, node := range fileInfoMap {
			if traversal.terminate {
				return
			} else if !node.found {
				traversal.uploadFile(baseDir, fileName, node, folderId)
			}
		}
	}

	for subDir, subFolderId := range subDirMap {
		fullPath := path.Join(baseDir, subDir)
		relativePath := path.Join(relativeDir, subDir)
		subDirName := subDir
		// fmt.Println("processing directory: " + fullPath + "  (" + relativePath + ")")

		if encrypted := traversal.cryptoService.EncryptFileNameWithMd5(subDir, "00000000000000000000000000000000"); encrypted != `` {
			subDirName = encrypted
		}

		if subFolderId == `` {
			// fmt.Printf(">> folder [ %s %s ] needs to be created on gsuite drive\n", relativeDir, subDirName)

			if newFolderId, err := traversal.gdrive.CreateFolder(traversal.driveSvc, folderId, subDirName); err != nil {
				logger.Error("error: %v", err)
			} else {
				subDirMap[subDir] = newFolderId
				fmt.Printf(">> create folder: %s\n", subDir)
				subFolderId = newFolderId
			}
		}

		if traversal.terminate {
			return
		}

		if subFolderId != `` {
			traversal.recursiveLocal(fullPath, relativePath, subFolderId, depth+1)
		}
	}
}

func (traversal *traversalServiceImpl) setInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		fmt.Printf("\nReceived an interrupt, stopping services...\n\n")
		traversal.terminate = true
	}()
}

func (traversal *traversalServiceImpl) RecursiveRemote(driveSvc *drive.Service, baseDir string, relativeDir string, folderId string, speedUpload bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic: %v\n", r)
		}
	}()

	traversal.speedUpUpload = speedUpload
	traversal.driveSvc = driveSvc

	traversal.recursiveRemoteFileSystem(baseDir, relativeDir, folderId, true)
}

func (traversal *traversalServiceImpl) recursiveRemoteFileSystem(baseDir string,
	relativeDir string, folderId string, speedUpUpload bool) {
	// fmt.Println("basedir: " + baseDir)

	folders := traversal.retrieveRemoteFolders(baseDir, relativeDir, folderId)
	remoteFiles := traversal.retrieveRemoteFiles(folderId)
	dirs, files := traversal.retrieveLocal(baseDir)

	for _, remoteFile := range remoteFiles {
		encrypted := true

		realFileName, remoteMd5Sum, err := traversal.cryptoService.DecryptFileNameWithMd5(remoteFile)

		if err != nil {
			realFileName = remoteFile.Name
			remoteMd5Sum = remoteFile.Md5Checksum
			encrypted = false
		}

		localFilePath := path.Join(baseDir, realFileName)

		downloadFileFlag := false

		if file, found := files[realFileName]; found {
			var data models.SyncFileInfo

			md5sum, err := traversal.handler.Md5Fetcher(baseDir, relativeDir, file, &data)

			if err != nil || md5sum != remoteMd5Sum {
				if fileModifiedTime, err := google_drive.ToLocalTime(remoteFile.ModifiedTime); err != nil {
					logger.Error(err)
					panic(err.Error())
				} else {
					if fileModifiedTime.After(file.ModTime()) {
						// overwrite id
						// traversal.xdb.Gorm.Delete(&data)
						traversal.handler.DeleteBeforeDownload(localFilePath, &data)
						//traversal.conn.Delete(&data)
						//os.Remove(localFilePath)
						downloadFileFlag = true
						//fmt.Print("Delete local: " + localFilePath)
					} else {
						fmt.Println("Ignore file: " + file.Name() + " (newer)")
					}
				}
			} else {
				fmt.Println("Ignore file: " + file.Name() + " (already exists)")
			}
		} else {
			downloadFileFlag = true
		}

		if downloadFileFlag {
			fmt.Print("Download: " + localFilePath + " ... ")

			if err := traversal.handler.DownloadRemoteFile(remoteFile, localFilePath, remoteMd5Sum, encrypted); err != nil {
				logger.Error("error: %v", err)
			} else {
				fmt.Println("ok")
			}

			traversal.downloadCount++
			traversal.downloadSize += remoteFile.Size
		}

		if traversal.terminate {
			return
		}
	}

	for _, folder := range folders {
		realFolderName, _, _ := traversal.cryptoService.DecryptFileNameWithMd5(folder)

		nextDirectoryBaseDir := path.Join(baseDir, realFolderName)
		nextRelativeDir := path.Join(relativeDir, realFolderName)

		if _, found := dirs[realFolderName]; !found {
			if err := traversal.handler.Mkdir(nextDirectoryBaseDir); err != nil {
				traversal.terminate = true
				return
			}
		}

		traversal.recursiveRemoteFileSystem(nextDirectoryBaseDir, nextRelativeDir, folder.Id, speedUpUpload)

		if traversal.terminate {
			return
		}
	}
}
